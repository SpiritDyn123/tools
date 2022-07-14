package codec

import "testing"

func Test_gzip(t *testing.T) {
	zipper := &GZipper{}

	data := make([]byte, 1024)
	data, err := zipper.Zip(data)
	if err != nil {
		t.Errorf("Gzip zip err:%v", err)
		return
	}
	t.Logf("zip data:%v len:%d", data, len(data))

	data, err = zipper.Unzip(data)
	if err != nil {
		t.Errorf("Gzip unzip err:%v", err)
		return
	}
	t.Logf("unzip data len:%d", len(data))
}

//go test -bench=GZip
func BenchmarkGZip(b *testing.B) {
	zipper := &GZipper{}
	data := make([]byte, 1024)
	var err error
	for i := 0;i < b.N;i++ {
		data, err = zipper.Zip(data)
		if err != nil {
			b.Errorf("Gzip zip err:%v", err)
			return
		}

		//b.Logf("zip len:%d", len(data))

		data, err = zipper.Unzip(data)
		if err != nil {
			b.Errorf("Gzip unzip err:%v", err)
			return
		}
		//b.Logf("unzip data len:%d", len(data))
	}
}

