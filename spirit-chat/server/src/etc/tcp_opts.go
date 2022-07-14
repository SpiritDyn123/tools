package etc

import (
	"encoding/binary"
	"spirit-chat/lib/tcp"
	"spirit-chat/lib/tcp/codec"
)

var (
	Codec tcp.ICodec
)

func initCodec() {
	var endian binary.ByteOrder = binary.BigEndian
	if !Big_endian {
		endian = binary.LittleEndian
	}

	Codec = codec.NewBinaryCodec(Include_len, byte(Head_len), endian)
}

var (
	Msg_marshaler 		tcp.IMsgMarshaler
	Zipper 				tcp.IZipper
	Encryptor			tcp.IEncryptor
)

type MsgHello struct {
	Id int
	Msg string
}

type MsgAdd struct {
	A int
	B int
	Sum int
}

func InitTcpOpts() (opts []tcp.TcpOption) {
	msg_m := &codec.JsonMsgMarshaler{}
	msg_m.Register(MsgHello{})
	msg_m.Register(MsgAdd{})
	Msg_marshaler = msg_m

	Zipper = &codec.GZipper{}

	encryptor := &codec.AESEncryptor{}
	err := encryptor.Init([]byte("spiritmoon123456"))
	if err != nil {
		panic(err)
	}

	Encryptor = encryptor

	opts = append(opts, tcp.MsgMarshalerOption(Msg_marshaler))
	opts = append(opts, tcp.ZipOption(Zipper))
	opts = append(opts, tcp.EncryptOption(encryptor))
	return
}

