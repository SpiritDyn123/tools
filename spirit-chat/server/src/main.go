package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"spirit-chat/lib/tcp"
	"spirit-chat/server/src/chat"
	"spirit-chat/server/src/etc"
	"syscall"
)


func main() {
	app := &cli.App{
		Name: "spirit-chat-server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "tcp_addr",
				Aliases: []string{ "addr", "a"},
				Usage: "tcp listen addr",
				Value: ":9999",
				Destination: &etc.Tcp_addr,
			},

			&cli.StringFlag{
				Name: "log_file",
				Aliases: []string{ "log", "l"},
				Usage: "log file name",
				Value: "../log/spirit-chat-svr.log",
				Destination: &etc.Log_file,
			},

			&cli.IntFlag{
				Name: "head_len",
				Aliases: []string{ "hl", },
				Usage: "msg head length",
				Value: 4,
				Destination: &etc.Head_len,
			},

			&cli.BoolFlag{
				Name: "include_len",
				Aliases: []string{ "il", },
				Usage: "length include head_len",
				Value: true,
				Destination: &etc.Include_len,
			},

			&cli.BoolFlag{
				Name: "big_endian",
				Aliases: []string{ "be", },
				Usage: "big endian or little",
				Value: true,
				Destination: &etc.Big_endian,
			},
		},
		Action: action,
	}

	app.Run(os.Args)
}

func action(c *cli.Context) (err error) {
	etc.Init()

	//启动服务器
	svr := tcp.TcpServer{}
	go svr.RunLoop(etc.Tcp_addr, &chat.ChatSvr{}, etc.Codec, etc.InitTcpOpts()...)

	logrus.Infof("=========sever start success========")

	//等待结束
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGKILL)
	csig := <- sig

	svr.Stop()

	logrus.Infof("close by %v", csig)
	return
}

