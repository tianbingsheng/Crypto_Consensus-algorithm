package main

import (
	"crypto/aes"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"fmt"
	"encoding/hex"
	"encoding/base64"
)

//密码分组中CFB分组模式的编程

//通过CFB分组模式加密
func AesCFBEncrypt(plainTxt []byte,key []byte) []byte{
	//key是否合法
	block,_:=aes.NewCipher(key)
	cipherTxt:=make([]byte,aes.BlockSize+len(plainTxt))
	iv:=cipherTxt[:aes.BlockSize]

	//向iv切片数组初始化rand.Reader（随机内存流）
	io.ReadFull(rand.Reader,iv)

	//设置加密模式为CFB
	stream:=cipher.NewCFBEncrypter(block,iv)

	//加密
	stream.XORKeyStream(cipherTxt[aes.BlockSize:],plainTxt)

	//cipherTxt 包含了key和明问的两部分加密的内容
	return cipherTxt

}


//通过AES算法，利用CFB分组模式解密
func AesCFBDecrypt(cipherTxt []byte,key []byte) []byte {
	block,_:=aes.NewCipher(key)

	//拆分iv和密文
	iv:=cipherTxt[:aes.BlockSize]
	cipherTxt = cipherTxt[aes.BlockSize:]

	//设置解密模式
	stream:=cipher.NewCFBDecrypter(block,iv)

	var des =make([]byte,len(cipherTxt))

	//解密
	stream.XORKeyStream(des,cipherTxt)

	return des
}

func main() {


	//对称加密DES，key为8
	//对称加密3DES，key为24
	//对称加密AES，可以16,24,32
	var cipher = AesCFBEncrypt([]byte("hello 123"),[]byte("1234567890123456"))


	//通过编码，编译用户可以看到的密文
	fmt.Println(hex.EncodeToString(cipher))
	fmt.Println(base64.StdEncoding.EncodeToString(cipher))


	//解密
	var des=AesCFBDecrypt(cipher,[]byte("1234567890123456"))
	fmt.Println(string(des))




}
