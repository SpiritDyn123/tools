package etc

import (
	"spirit-chat/lib/log"
)

var (
	Tcp_addr string
	Log_file string

	Head_len  int
	Big_endian bool
	Include_len bool
)

func Init() {
	initLog()

	initCodec()
}

func initLog() {
	log.InitLog(Log_file, "debug")
}

