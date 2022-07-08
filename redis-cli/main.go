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
	defer conn.Close()

	fmt.Println("======Welcome to spirit-moon test-tools redis-cli(golang)======")

	cmd := &RedisCmd{
		Num_to_str: num_to_str,
	}
	buffer := bufio.NewReader(os.Stdin)


	buf := bytes.NewBuffer(nil)
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

		req_cmd := arr[0]
		if req_cmd == "exit" {
			break
		}

		cmd.Clear()
		args := []interface{}{}
		for _, arg_str := range arr[1:] {
			args = append(args, arg_str)
		}

		cmd_str, err := cmd.Encode(req_cmd, args...)
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

		data_len, rcount, err := cmd.Read(conn, buf)
		if err != nil {
			panic(err)
		}

		resp_data := make([]byte, data_len)
		buf.Read(resp_data)
		if show_cmd {
			fmt.Printf("RECV(%d): %s\n--------------\n", rcount, strings.ReplaceAll(string(resp_data), "\r\n", "\\r\\n"))
		}

		cmd.Clear()
		_, err = cmd.Decode(string(resp_data))
		if err != nil {
			fmt.Println("cli_decode_err:", err)
			continue
		}

		cmd.Cmd.ShowCommand("")

		req_cmd = strings.ToLower(req_cmd)
		if req_cmd == "quit" {
			break
		}else if req_cmd == "monitor" { //原始转发数据
			rdata := make([]byte, 1024)
			for {
				n, err := conn.Read(rdata)
				if err != nil {
					panic(err)
				}

				resp_data := string(rdata[:n])
				fmt.Printf(string(resp_data))
			}
		} else if req_cmd == "subscribe" { //订阅
			for {
				data_len, rcount, err = cmd.Read(conn, buf)
				if err != nil {
					panic(err)
				}
				resp_data := make([]byte, data_len)
				buf.Read(resp_data)

				cmd.Clear()
				_, err = cmd.Decode(string(resp_data))
				if err != nil {
					panic(err)
				}

				cmd.Cmd.ShowCommand("")
			}
		}
	}

}
