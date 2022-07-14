package codec

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

//go test -bench=AES
func BenchmarkAESEncryptor(b *testing.B) {
	encryptor := &AESEncryptor{}
	err := encryptor.Init([]byte("abcdef0123456789"))
	if err != nil {
		b.Errorf("encryptor init error:%v",  err)
		return
	}

	prefix := "hello world================>"
	var wg sync.WaitGroup
	for i := 0;i < b.N;i++ {
		id := i+1
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data := []byte(fmt.Sprintf("%s%d", prefix, id))
			edata, err := encryptor.Encrypt(data)
			if err != nil {
				b.Errorf("id:%d encryptor encrypt error:%v", id,  err)
				return
			}

			//b.Logf("id:%d encryptor encrypt data:%v len:%d", id, edata, len(edata))

			ddata, err := encryptor.Decrypt(edata)
			if err != nil {
				b.Errorf("id:%d encryptor decrypt error:%v", id, err)
				return
			}

			//b.Logf("id:%d encryptor decrypt data:%v, len:%d", id, string(ddata), len(ddata))
			if string(ddata[len(prefix):]) != strconv.Itoa(id) {
				b.Fatalf("id:%d error!!!", id)
			}
		}(id)
	}
	wg.Wait()
	b.Logf("------------end----------------")
}
