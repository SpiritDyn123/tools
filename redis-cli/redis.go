package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

/*
规范格式(redis-cli) RESP
	1）间隔符号，在Linux下是\r\n，在Windows下是\n
	2）简单字符串 Simple Strings, 以 "+"加号 开头
	3）错误 Errors, 以"-"减号 开头
	4）整数型 Integer， 以 ":" 冒号开头
	5）大字符串类型 Bulk Strings, 以 "$"美元符号开头，长度限制512M
	6）数组类型 Arrays，以 "*"星号开头
*/
type RedisCmdItem struct {
	Empty    bool
	Status  string "+,-,*,:,$"
	Value	string "comands(get,set...), length or args"
	Cmds    []*RedisCmdItem
}

//原始数据
func (this *RedisCmdItem) Raw() (raw string) {
	if !this.Empty {
		raw += fmt.Sprintf("%s%s\r\n", this.Status, this.Value)
	}

	for _, citem := range this.Cmds {
		raw += citem.Raw()
	}

	return
}

//命令行显示
func (this *RedisCmdItem) ShowCommand(prefix string) {
	if this.Empty {
		for _, citem := range this.Cmds {
			citem.ShowCommand(prefix)
		}
		return
	}

	switch this.Status {
	case "+":
		fmt.Println(this.Value)
	case "-":
		fmt.Printf("(error) %s\n", this.Value)
	case "$":
		if this.Value == "-1" {
			fmt.Println("(nil)")
		}
	case ":":
		fmt.Printf("(integer) %s\n", this.Value)
	case "*":
		index := 1
		if len(this.Cmds) == 0 {
			fmt.Println("(empty list or set)")
			return
		}

		for _, ccitem := range this.Cmds {
			fmt.Printf(prefix)
			if ccitem.Status != "$" || ccitem.Value == "-1" {
				fmt.Printf("%d)", index)
				index++
			}

			ccitem.ShowCommand(prefix + "\t")
		}
	default: //常规字符串 case "":
		fmt.Printf("\"%s\"\n", this.Value)
	}
}

type RedisCmd struct {
	Raw		string
	Cmd		*RedisCmdItem
	Num_to_str bool
}

func (this *RedisCmd) Clear() {
	this.Raw = ""
	this.Cmd = nil
}

func (this *RedisCmd) Encode(cmd string, args ...interface{}) (string, error) {
	if this.Raw != "" {
		return this.Raw, nil
	}

	args_len := len(args)
	this.Cmd = &RedisCmdItem{
		Empty: true,
	}
	//if args_len > 0 { //多个args，以数组方式拼接（貌似必须以数组方式发送）
		this.Cmd  = &RedisCmdItem{
			Empty: false,
			Status: "*",
			Value: strconv.Itoa(args_len + 1),
		}
	//}

	this.Cmd.Cmds = append(this.Cmd.Cmds,
		&RedisCmdItem{
			Status: "$",
			Value: strconv.Itoa(len(cmd)),
		},
		&RedisCmdItem{
			Status: "",
			Value: cmd,
		})

	for i, arg := range args  {
		str_arg := fmt.Sprintf("%v", arg)
		citem := &RedisCmdItem{
			Status: "$",
			Value: strconv.Itoa(len(str_arg)),
		}

		rkind := reflect.TypeOf(arg).Kind()
		if rkind == reflect.Bool { //bool 值转为int
			if arg.(bool) {
				arg = 1
			} else {
				arg = 0
			}
			rkind = reflect.Int
		}

		add_arg := true
		switch rkind {
		case reflect.String:
		case reflect.Float32, reflect.Float64:
		case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint,
			reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int: //其实都可以转成字符串
			if !this.Num_to_str {
				citem.Status = ":"
				citem.Value = str_arg
				add_arg = false
			}
		default:
			return "", fmt.Errorf("invalid args[%d]:%v type:%v", i, arg, rkind)
		}

		this.Cmd.Cmds = append(this.Cmd.Cmds, citem)
		if add_arg {
			this.Cmd.Cmds = append(this.Cmd.Cmds, &RedisCmdItem{
				Status: "",
				Value: str_arg,
			})
		}
	}

	//生成字符串
	this.Raw = this.Cmd.Raw()
	return this.Raw, nil
}

func (this *RedisCmd) genCmd(cmd_str string) (cmds []*RedisCmdItem, offset int, err error) {
	sp_index := strings.Index(cmd_str, "\r\n")
	if sp_index <= 0 {
		err = fmt.Errorf("invalid format split index")
		return
	}

	cmd := &RedisCmdItem{
		Status: cmd_str[:1],
		Value: cmd_str[1:sp_index],
	}
	cmds = append(cmds, cmd)

	offset = sp_index + 2
	cstr_len := len(cmd_str)
	switch(cmd_str[0]) {
	case '+', '-': //操作成功或失败的描述
	case ':': //整数
	case '*': //数组
		args_len, gerr := strconv.Atoi(cmd.Value)
		if gerr != nil {
			return nil, 0, gerr
		}

		for args_len > 0 {
			if offset >= cstr_len {
				err = fmt.Errorf("invalid format *(array) len")
				return
			}

			c_cmds, coffset, gerr := this.genCmd(cmd_str[offset:])
			if gerr != nil {
				return nil, 0, gerr
			}

			cmd.Cmds = append(cmd.Cmds, c_cmds...)
			offset += coffset
			args_len--
		}
	case '$': //字符串长度
		if cmd.Value == "-1" {
			return
		}

		str_len, gerr := strconv.Atoi(cmd.Value)
		if gerr != nil {
			return nil, 0, gerr
		}

		if offset + str_len >= cstr_len {
			err = fmt.Errorf("invalid format $(string) len")
			return
		}

		cmds = append(cmds, &RedisCmdItem{
			Status: "",
			Value: cmd_str[offset:offset + str_len],
		})
		offset += str_len + 2
	default:
		err = fmt.Errorf("invalid format status:%v", cmd_str[0])
		return
	}

	return
}

func (this *RedisCmd) Read(c io.Reader, buf []byte) (resp_data []byte, rcount int, err error) {
	rbuf := bytes.NewBuffer(nil)
	var n int
	for {
		rcount++
		n, err = c.Read(buf)
		if err != nil {
			return
		}

		rbuf.Write(buf[:n])
		resp_data = rbuf.Bytes()
		_, n, err = this.genCmd(string(resp_data))
		if err == nil {
			resp_data = resp_data[:n]
			break
		}
	}

	return
}

func (this *RedisCmd) Decode(raw_data string) (cmd *RedisCmdItem, err error) {
	//eg: raw_data = "*3\r\n$3\r\naaa\r\n$5\r\n12345\r\n:6" is "aaa" "12345" 6
	c, _, err := this.genCmd(raw_data)
	if err != nil {
		return
	}

	this.Raw = raw_data
	this.Cmd = &RedisCmdItem{
		Empty: true,
		Cmds: c,
	}
	cmd = this.Cmd
	//fmt.Println("==============>", strings.ReplaceAll(this.Cmd.Raw(), "\r\n", "\\r\\n"))

	return
}
