package chat

import (
	"encoding/binary"
	"spirit-chat/lib/log"
	"spirit-chat/lib/tcp"
	"spirit-chat/lib/tcp/codec"
)

var (
	Etc settings
)

type settings struct {
	Tcp_addr string
	Log_file string

	Head_len  int
	Big_endian bool
	Include_len bool
}

func Init() {
	initLog()
	initCodec()
}

func initLog() {
	log.InitLog(Etc.Log_file, "debug")
}

var (
	Codec tcp.ICodec
)
func initCodec() {
	var endian binary.ByteOrder = binary.BigEndian
	if !Etc.Big_endian {
		endian = binary.LittleEndian
	}
	Codec = codec.NewBinaryCodec(Etc.Include_len, byte(Etc.Head_len), endian)
}