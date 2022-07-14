package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"spirit-chat/goclient/src/chat"
	"spirit-chat/lib/tcp"
	"spirit-chat/server/src/etc"
	"syscall"
)

func main() {
	app := &cli.App{
		Name: "spirit-chat-client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "tcp_addr",
				Aliases: []string{ "addr", "a"},
				Usage: "tcp listen addr",
				Value: "127.0.0.1:9999",
				Destination: &chat.Etc.Tcp_addr,
			},

			&cli.StringFlag{
				Name: "log_file",
				Aliases: []string{ "log", "l"},
				Usage: "log file name",
				Value: "../log/spirit-chat-cli.log",
				Destination: &chat.Etc.Log_file,
			},

			&cli.IntFlag{
				Name: "head_len",
				Aliases: []string{ "hl", },
				Usage: "msg head length",
				Value: 4,
				Destination: &chat.Etc.Head_len,
			},

			&cli.BoolFlag{
				Name: "include_len",
				Aliases: []string{ "il", },
				Usage: "length include head_len",
				Value: true,
				Destination: &chat.Etc.Include_len,
			},

			&cli.BoolFlag{
				Name: "big_endian",
				Aliases: []string{ "be", },
				Usage: "big endian or little",
				Value: true,
				Destination: &chat.Etc.Big_endian,
			},
		},
		Action: action,
	}

	app.Run(os.Args)
}

func action(c *cli.Context) (err error) {
	chat.Init()

	tcp_cli := &tcp.TcpClient{}
	go tcp_cli.Dial(chat.Etc.Tcp_addr, true, &chat.ChatCli{}, chat.Codec, etc.InitTcpOpts()...)

	//等待结束
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGKILL)
	csig := <- sig

	tcp_cli.Stop()

	logrus.Infof("close by %v", csig)
	return
}
