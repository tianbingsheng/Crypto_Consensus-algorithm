package main

import (
	"io"
	"crypto/cipher"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// 加密函数
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	 // NewOFB返回一个在输出反馈模式下使用分组密码b进行加密或解密的Stream。初始化矢量iv的长度必须等于b的块大小。
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// 解密函数
func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.

	iv := ciphertext[:aes.BlockSize]

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	plaintext2 := make([]byte, len(ciphertext))
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])

	return plaintext2, nil
}

func main() {
	key := []byte("6368616e676520746869732070617374")
	// 加密
	result, err := AesEncrypt([]byte("hello world"), key)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result))

	// 解密
	origData, err := AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
