package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)


//go test -v redis_test.go redis.go -run=Encode
func TestRedisCmd_Encode(t *testing.T) {
	rcmd := &RedisCmd{}

	req_data, err := rcmd.Encode("HMGET", "thash", "k1", "k2")
	if err != nil {
		t.Errorf("ecnode cmd err:%v", err)
		return
	}

	//输出 *4\r\n$5\r\nHMGET\r\n$5\r\nthash\r\n$2\r\nk1\r\n$2\r\nk2\r\n
	t.Logf("send cmd: %s\n", strings.ReplaceAll(req_data, "\r\n", "\\r\\n"))
}

//go test -v redis_test.go redis.go -run=Decode
func TestRedisCmd_Decode(t *testing.T) {
	rcmd := &RedisCmd{}

	//hgetall k1 k2的值
	resp_data := "*3\r\n$2\r\nv1\r\n$2\r\nv2\r\n$-1\r\n"
	cmd, err := rcmd.Decode(resp_data)
	if err != nil {
		t.Errorf("decode cmd err:%v", err)
		return
	}

	//testing模式 “可能” 不显示
	cmd.ShowCommand("")

	jsdata, _ := json.Marshal(cmd)
	t.Logf("decode cmd to JSON:%s", string(jsdata))
}

//go test -v redis_test.go redis.go -run=Read
func TestRedisCmd_Read(t *testing.T) {
	test_cmd_str := "*4\r\n$5\r\nHMGET\r\n$5\r\nthash\r\n$2\r\nk1\r\n$2\r\nk2\r\n"
	src_buf := bytes.NewBuffer(nil)
	write_cnt := 3
	for i := 0;i < write_cnt;i++ { //多次放在一起粘包
		src_buf.Write([]byte(test_cmd_str))
	}

	dst_buf := bytes.NewBuffer(nil)
	cmd := &RedisCmd{}
	for write_cnt > 0 {
		n, rc, err := cmd.Read(src_buf, dst_buf)
		if err != nil {
			t.Errorf("read err:%v\n", err)
			return
		}
		rdata := make([]byte, n)
		dst_buf.Read(rdata)
		write_cnt--
		t.Logf("%d、=====data_len:%d====readcount:%d====%s\n", 3-write_cnt, n, rc,
			strings.ReplaceAll(string(rdata), "\r\n", "\\r\\n"))
	}
}
