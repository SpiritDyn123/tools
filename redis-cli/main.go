package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

func init() {

}

func main() {
	//原生tcp不封装，模拟redis-cli
	host_addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", host_addr)
	if err != nil {
		panic(err)
	}

	cmd := &RedisCmd{}
	buffer := bufio.NewReader(os.Stdin)


	rdata := make([]byte, 5) //处理粘包问题，buf刻意设置很小
	rbuf := bytes.NewBuffer(nil)

	for {
		fmt.Printf("%s> ", host_addr)

		cmd_info, err := buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		cmd_info = strings.TrimSpace(cmd_info)
		if cmd_info == "" {
			continue
		}

		arr := strings.Split(cmd_info, " ")
		if len(arr) <= 0 {
			fmt.Println("invalid cmd:", cmd_info)
			continue
		}

		cmd.Clear()
		args := []interface{}{}
		for _, arg_str := range arr[1:] {
			args = append(args, arg_str)
		}

		cmd_str, err := cmd.Encode(arr[0], args...)
		if err != nil {
			fmt.Println("cli_encode_err:", err)
			continue
		}

		if show_cmd {
			fmt.Printf("SEND: %s\n--------------\n", strings.ReplaceAll(cmd_str, "\r\n", "\\r\\n"))
		}

		_, err = conn.Write([]byte(cmd_str))
		if err != nil {
			panic(err)
		}

		var cmd_data []byte
		rcount := 0
		for {
			rcount++
			n, err := conn.Read(rdata)
			if err != nil {
				panic(err)
			}

			rbuf.Write(rdata[:n])
			cmd_data = rbuf.Bytes()
			_, n, err = cmd.genCmd(string(cmd_data))
			if err == nil {
				cmd_data = cmd_data[:n]
				rbuf.Reset()
				rbuf.Write(cmd_data[n:])
				break
			}
		}

		resp_data := string(cmd_data)
		if show_cmd {
			fmt.Printf("RECV: %s\n--------------\n", strings.ReplaceAll(resp_data, "\r\n", "\\r\\n"))
		}

		cmd.Clear()
		_, err = cmd.Decode(resp_data)
		if err != nil {
			fmt.Println("cli_decode_err:", err)
			continue
		}

		cmd.Cmd.ShowCommand("")

		if cmd_info == "monitor" {
			for {
				n, err := conn.Read(rdata)
				if err != nil {
					panic(err)
				}

				resp_data := string(rdata[:n])
				fmt.Println(string(resp_data))
			}
		}
	}

}
