package main

import (
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
