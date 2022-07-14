package chat

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"spirit-chat/lib/tcp"
	"spirit-chat/lib/tcp/codec"
	"time"
)

type ChatSvr struct {
	
}

func (this *ChatSvr) OnConnect(session tcp.ISession) {
	log.Printf("implement me OnAccept")
}

func (this *ChatSvr) OnMessage(session tcp.ISession, msgs []interface{}) {
	data, _ := json.Marshal(msgs)
	log.Printf("implement me OnMessage:%+v", string(data))
	time.Sleep(time.Second)

	head := msgs[0].(*codec.Head)
	head.Zip = false
	head.Encrypt = false
	session.Send(msgs...)
}

func (this *ChatSvr) OnDisconnect(session tcp.ISession, err error) {
	log.Printf("implement me OnDisconnect")
}

