package tcp

type ITcpOptions interface {
	GetZipper() IZipper
	GetEncryptor() IEncryptor
	GetMsgMarshaler() IMsgMarshaler
}

type TcpOption struct {
	f func (opt *tcpOptions)
}

type tcpOptions struct {
	zipper        IZipper
	encrptor      IEncryptor
	msg_marshaler IMsgMarshaler
}

func (this *tcpOptions) GetZipper() IZipper {
	return this.zipper
}

func (this *tcpOptions) GetEncryptor() IEncryptor {
	return this.encrptor
}

func (this *tcpOptions) GetMsgMarshaler() IMsgMarshaler {
	return this.msg_marshaler
}

func ZipOption(zipper IZipper) TcpOption {
	return TcpOption{
		f: func(opt *tcpOptions) {
			opt.zipper = zipper
		},
	}
}

func EncryptOption(encrptor IEncryptor) TcpOption {
	return TcpOption{
		f: func(opt *tcpOptions) {
			opt.encrptor = encrptor
		},
	}
}

func MsgMarshalerOption(msg_marshaler IMsgMarshaler) TcpOption {
	return TcpOption{
		f: func(opt *tcpOptions) {
			opt.msg_marshaler = msg_marshaler
		},
	}
}