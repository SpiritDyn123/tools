package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AESEncryptor struct {
	key []byte
	block 	cipher.Block
}

func (this *AESEncryptor) Init(key []byte) (err error) {
	this.block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	this.key = key
	return
}

func (this *AESEncryptor) pcs5Panding(data []byte, block_size int) []byte {
	padding := block_size - len(data)%block_size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func (this *AESEncryptor) pcs5Unpanding(data []byte) []byte {
	data_len := len(data)
	padding_len := int(data[data_len - 1])
	return data[:data_len - padding_len]
}

func (this *AESEncryptor) Encrypt(data []byte) (edata []byte, err error) {
	block := this.block
	blockSize := block.BlockSize()
	padding_data := this.pcs5Panding(data, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, this.key[:blockSize])
	edata = make([]byte, len(padding_data))
	blockMode.CryptBlocks(edata, padding_data)
	return
}

func (this *AESEncryptor) Decrypt(edata []byte) (data []byte, err error) {
	block := this.block
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, this.key[:blockSize])
	data = make([]byte, len(edata))
	blockMode.CryptBlocks(data, edata)
	data = this.pcs5Unpanding(data)
	return
}

type RSAEncryptor struct {

}

func (this *RSAEncryptor) Encrypt(data []byte) (edata []byte, err error) {

	return
}

func (this *RSAEncryptor) Decrypt(edata []byte) (data []byte, err error) {
	return
}


type RSA_AESEcryptor struct {

}

func (this *RSA_AESEcryptor) Encrypt(data []byte) (edata []byte, err error) {

	return
}

func (this *RSA_AESEcryptor) Decrypt(edata []byte) (data []byte, err error) {
	return
}

