package tcp

import (
	"io"
	"net"
)

type ISession interface {
	ID() uint64
	Addr() net.Addr
	Send(...interface{})
	Close()
}

type ICodec interface {
	Read(io.Reader, ITcpOptions) ([]interface{}, error)
	Write(io.Writer, []interface{}, ITcpOptions) error
}

type IEventListener interface {
	OnConnect(ISession)
	OnMessage(ISession, []interface{})
	OnDisconnect(ISession, error)
}

type IZipper interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}

type IEncryptor interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type IMsgMarshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte) (interface{}, error)
}