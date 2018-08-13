package main

import (
	"crypto/aes"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"fmt"
)

//通过AES测试OFB
func AesOFBEncryt(plaintxt []byte ,key []byte) []byte {
	block,_:=aes.NewCipher(key)
	ciphertxt:=make([]byte,aes.BlockSize+len(plaintxt))
	iv:=ciphertxt[:block.BlockSize()]
	io.ReadFull(rand.Reader,iv)
	//设置加密模式
	//strem:=cipher.NewOFB(block,iv)
	strem:=cipher.NewCTR(block,iv)
	//其实是将加密后的密文存到ciphertxt[aes.BlockSize:]
	strem.XORKeyStream(ciphertxt[aes.BlockSize:],plaintxt)
	return ciphertxt

}

//解密
func AesOFBDecrypt(ciphertxt []byte,key []byte) []byte {
	block,_:=aes.NewCipher(key)
	iv:=ciphertxt[:aes.BlockSize]
	//存储解密后的信息
	plaintxt:=make([]byte,len(ciphertxt)-aes.BlockSize)
	//设置解密方式OFB
	//stream:=cipher.NewOFB(block,iv)
	stream:=cipher.NewCTR(block,iv)
	stream.XORKeyStream(plaintxt,ciphertxt[aes.BlockSize:])

	return plaintxt
}


func main () {
	ciphertxt:=AesOFBEncryt([]byte("hello world 123"),[]byte("1234567890123456"))
	fmt.Println(ciphertxt)

	fmt.Println(string(AesOFBDecrypt(ciphertxt,[]byte("1234567890123456"))))
}
