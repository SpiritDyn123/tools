package codec

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type StringMsgMarshaler struct {

}

func (this *StringMsgMarshaler) Marshal(msg interface{}) (data []byte, err error) {
	switch msg.(type) {
	case []byte:
		data = msg.([]byte)
	case string:
		data = []byte(msg.(string))
	default:
		err = fmt.Errorf("StringMsgMarshaler::Marshal invalid msg type:%v", reflect.TypeOf(msg))
	}
	return
}

func (this *StringMsgMarshaler) Unmarshal(data []byte) (msg interface{}, err error) {
	return string(data), nil
}

type JsonMsgMarshaler struct {
	m_msgs map[string]reflect.Type
}

func (this *JsonMsgMarshaler) getKey(tname string) string {
	index := strings.LastIndex(tname, ".")
	if index < 0 {
		return tname
	}

	return tname[index+1:]
}

func (this *JsonMsgMarshaler) Register(msg interface{}) (err error) {
	if msg == nil {
		err = fmt.Errorf("JsonMsgMarshaler::Register msg nil")
		return
	}

	rt := reflect.TypeOf(msg)
	key := this.getKey(rt.String())
	if this.m_msgs == nil {
		this.m_msgs = map[string]reflect.Type{}
	}

	if this.m_msgs[key] != nil {
		err = fmt.Errorf("JsonMsgMarshaler::Register msg:%s repeated", key)
		return
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	this.m_msgs[key] = rt
	return
}

func (this *JsonMsgMarshaler) Marshal(msg interface{}) (data []byte, err error) {
	mdata, err := json.Marshal(msg)
	if err != nil {
		return
	}

	rt := reflect.TypeOf(msg)
	key := this.getKey(rt.String())
	data = make([]byte, 4 + len(key) + len(mdata))
	binary.LittleEndian.PutUint32(data, uint32(len(key)))
	copy(data[4:], []byte(key))
	copy(data[4 + len(key):], mdata)
	return
}

func (this *JsonMsgMarshaler) Unmarshal(data []byte) (msg interface{}, err error) {
	if len(data) <= 0 {
		return
	}

	if len(data) < 4 {
		err = fmt.Errorf("JsonMsgMarshaler::Unmarshal data len < 4")
		return
	}

	klen := int(binary.LittleEndian.Uint32(data))
	if len(data[4:]) < klen {
		err = fmt.Errorf("JsonMsgMarshaler::Unmarshal key len err")
		return
	}

	key := string(data[4:4 + klen])
	mt := this.m_msgs[key]
	if mt == nil {
		err = fmt.Errorf("JsonMsgMarshaler::Unmarshal invalid key:%s", key)
		return
	}

	msg = reflect.New(mt).Interface()
	err = json.Unmarshal(data[4 + klen:], msg)
	return
}



