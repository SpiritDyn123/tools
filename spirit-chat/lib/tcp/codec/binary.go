package codec

import (
	"encoding/binary"
	"fmt"
	"io"
	"spirit-chat/lib/tcp"
)

type EndianToNFunc func(binary.ByteOrder, []byte)uint64
type EndianToBFunc func(binary.ByteOrder, []byte, uint64)

var (
	endian_ton_fs  = []EndianToNFunc{
		func(endian binary.ByteOrder, b []byte) uint64 { return uint64(b[0]) },
		func(endian binary.ByteOrder, b []byte) uint64 { return uint64(endian.Uint16(b)) },
		func(endian binary.ByteOrder, b []byte) uint64 { return uint64(endian.Uint32(b)) },
		func(endian binary.ByteOrder, b []byte) uint64 { return endian.Uint64(b) },
	}

	endian_tob_fs = []EndianToBFunc{
		func(endian binary.ByteOrder, b []byte, n uint64)  { b[0] = byte(n) },
		func(endian binary.ByteOrder, b []byte, n uint64) { endian.PutUint16(b, uint16(n)) },
		func(endian binary.ByteOrder, b []byte, n uint64) { endian.PutUint32(b, uint32(n)) },
		func(endian binary.ByteOrder, b []byte, n uint64) { endian.PutUint64(b, n) },
	}
)

type LenHead struct {
	Length uint64
}

type Head struct {
	LenHead
	Seq      			uint64
	Encrypt  			bool
	Zip  				bool
}

func NewBinaryCodec(include_head bool, head_len byte, endian binary.ByteOrder) tcp.ICodec {
	codec := &BinaryCodec{
		Include_head: include_head,
		Head_len: head_len,
		Endian: endian,
	}

	//修正headlen
	switch head_len {
	case 1:
		codec.Head_len = 1
	case 2:
		codec.Head_len = 2
	case 3,4:
		codec.Head_len = 4
	case 5,6,7,8:
		codec.Head_len = 8
	default:
		codec.Head_len = 4
	}

	codec.endian_index = int(codec.Head_len / 2)
	if codec.endian_index >= 3 {
		codec.endian_index = 3
	}

	codec.head_data = make([]byte, codec.Head_len)

	return codec
}

type BinaryCodec struct {
	Include_head         bool	// 长度是否包含包头长度
	Head_len	         byte   // head len length
	Endian			     binary.ByteOrder  // 大端还是小端

	head_data			[]byte
	endian_index		int
}

func (this *BinaryCodec) Read(conn io.Reader, opts tcp.ITcpOptions) (msgs []interface{}, err error) {
	_, err = io.ReadFull(conn, this.head_data)
	if err != nil {
		return
	}

	head := &Head{}
	head.Length = endian_ton_fs[this.endian_index](this.Endian, this.head_data)
	data_len := int(head.Length)
	if this.Include_head {
		data_len -= int(this.Head_len)
	}

	if data_len < 10 {
		err = fmt.Errorf("BinaryCodec::Read invalid msg len")
		return
	}

	buf := make([]byte, data_len)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return
	}

	head.Seq = this.Endian.Uint64(buf)
	head.Encrypt = buf[8] == 1
	head.Zip     = buf[9] == 1

	buf = buf[10:]
	//压缩
	if head.Zip {
		zipper := opts.GetZipper()
		if zipper == nil {
			err = fmt.Errorf("BinaryCodec::Read has no zipper")
			return
		}

		buf, err = zipper.Unzip(buf[:])
		if err != nil {
			return
		}
	}

	//解密
	if head.Encrypt {
		encryptor := opts.GetEncryptor()
		if encryptor == nil {
			err = fmt.Errorf("BinaryCodec::Read has no encryptor")
			return
		}

		buf, err = encryptor.Decrypt(buf, )
		if err != nil {
			return
		}
	}

	var msg interface{} = buf

	//msg 解析
	msg_marshaler := opts.GetMsgMarshaler()
	if msg_marshaler != nil {
		msg, err = msg_marshaler.Unmarshal(buf)
		if err != nil {
			return
		}
	}


	msgs = []interface{}{head, msg}

	return
}

func (this *BinaryCodec) Write(conn io.Writer, msgs []interface{}, opts tcp.ITcpOptions) (err error) {
	if len(msgs) < 1 {
		err = fmt.Errorf("BinaryCodec::Write need head at least")
		return
	}

	head, ok := msgs[0].(*Head)
	if !ok {
		err = fmt.Errorf("BinaryCodec::Write head type err")
		return
	}

	var msg []byte
	if len(msgs) > 1 {
		switch msgs[1].(type) {
		case []byte:
			msg = msgs[1].([]byte)
		case string:
			msg = []byte(msgs[1].(string))
		default:
			//msg 序列化
			msg_mashaler := opts.GetMsgMarshaler()
			if msg_mashaler == nil {
				err = fmt.Errorf("BinaryCodec::Write msg type hash no msg_mashaler")
				return
			}

			msg, err = msg_mashaler.Marshal(msgs[1])
			if err != nil {
				return
			}
		}
	}

	//加密
	if head.Encrypt {
		encryptor := opts.GetEncryptor()
		if encryptor == nil {
			err = fmt.Errorf("BinaryCodec::Write has no encryptor")
			return
		}

		msg, err = encryptor.Encrypt(msg)
		if err != nil {
			return
		}
	}

	//解密
	if head.Zip {
		zipper := opts.GetZipper()
		if zipper == nil {
			err = fmt.Errorf("BinaryCodec::Write has no zipper")
			return
		}

		msg, err = zipper.Zip(msg)
		if err != nil {
			return
		}
	}

	head.Length = uint64(8 + 1 + 1 + len(msg))
	buf := make([]byte, int(head.Length) + int(this.Head_len))
	if this.Include_head {
		head.Length += uint64(this.Head_len)
	}

	endian_tob_fs[this.endian_index](this.Endian, buf, head.Length)
	offset := int(this.Head_len)

	this.Endian.PutUint64(buf[offset:], head.Seq)
	offset+=8

	if head.Encrypt {
		buf[offset] = 1
	}
	offset+=1

	if head.Zip {
		buf[offset] = 1
	}
	offset+=1

	copy(buf[offset:], msg)

	var t, n int
	for {
		n, err = conn.Write(buf[t:])
		if err != nil {
			return
		}

		t += n
		if t == int(head.Length) {
			break
		}
	}

	return
}
