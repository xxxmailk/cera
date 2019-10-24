package view

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// padding data
func padData(src []byte, size int) []byte {
	num := size - len(src)%size
	pad := bytes.Repeat([]byte{byte(num)}, num)
	return append(src, pad...)
}

// unpadding data
func unPadData(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// encrypt by aes
func AESEncrypt(src, key []byte) (dst []byte, err error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	src = padData(src, b.BlockSize())
	blockMode := cipher.NewCBCEncrypter(b, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// decrypt by aes
func AESDecrypt(src, key []byte) (dst []byte, err error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockMode := cipher.NewCBCDecrypter(b, key)
	blockMode.CryptBlocks(src, src)
	return unPadData(src), nil
}
