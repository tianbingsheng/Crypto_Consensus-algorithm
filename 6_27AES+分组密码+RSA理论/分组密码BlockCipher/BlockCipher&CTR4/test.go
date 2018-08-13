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
	// 申明初始化获取一个新的密钥块。关键参数应该是AES密钥，16,24或32个字节来选择AES-128，AES-192或AES-256。
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 切片处理申明初始化一个较大长度的新字符串变量
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// 申明初始化，同时调用加密函数得到流接口
	stream := cipher.NewCTR(block, iv)
	// 流处理
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

	// 申明初始化，同时调用加密函数得到流接口
	stream := cipher.NewCTR(block, iv)
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