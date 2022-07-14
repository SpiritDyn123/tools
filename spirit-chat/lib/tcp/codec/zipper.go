package codec

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type GZipper struct {
	
}

func (this *GZipper) Zip(data []byte) (zdata[]byte, err error) {
	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)
	_, err = w.Write(data)
	if err != nil {
		return
	}

	err = w.Flush()
	if err != nil {
		w.Close()
		return
	}

	err = w.Close()
	if err != nil {
		return
	}

	zdata = buf.Bytes()
	return
}

func (this *GZipper) Unzip(zdata []byte) (data []byte, err error) {
	buf := bytes.NewBuffer(zdata)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	defer r.Close()
	data, err = ioutil.ReadAll(r)
	return
}
