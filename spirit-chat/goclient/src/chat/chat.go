package chat

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"spirit-chat/lib/tcp"
	"spirit-chat/lib/tcp/codec"
	"spirit-chat/server/src/etc"
	"time"
)

type ChatCli struct {
	
}

func (this *ChatCli) OnConnect(session tcp.ISession) {
	log.Printf("implement me OnConnect")
	session.Send(&codec.Head{
		Seq: 100,
		Zip: true,
		Encrypt: true,
	}, &etc.MsgHello{
		Id: 1,
		Msg: "hello",
	})
}

func (this *ChatCli) OnMessage(session tcp.ISession, msgs []interface{}) {
	data, _ := json.Marshal(msgs)
	log.Printf("implement me OnMessage:%+v", string(data))
	time.Sleep(time.Second)

	head := msgs[0].(*codec.Head)
	head.Encrypt = true
	head.Zip = true
	session.Send(msgs...)}

func (this *ChatCli) OnDisconnect(session tcp.ISession, err error) {
	log.Printf("implement me OnDisconnect")
}

