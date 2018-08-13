package main

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/md5"
	"crypto"

	"fmt"
)

//研究RSA实现签名和验签

//签名和验签的作用为，能够确定数据是否与发送者一致

func main() {

	//生成私钥
	priv,_:=rsa.GenerateKey(rand.Reader,1024)

	//通过私钥生成公钥
	pub:=&priv.PublicKey

	//通过hash散列对准备签名的名为做hash散列
	plaitxt:=[]byte("hello world")

	//实现散列过程
	h:=md5.New()
	h.Write(plaitxt)
	hashed:=h.Sum(nil)

	//通过pss函数，实现对明文hello world的签名
	//pss函数可以添加杂质，能够使得签名过程更安全
	opts:=rsa.PSSOptions{rsa.PSSSaltLengthAuto,crypto.MD5}

	//实现签名
	sig,_:=rsa.SignPSS(rand.Reader,priv,crypto.MD5,hashed,&opts)

	//sig就是RSA对“hello world”签名结果


	fmt.Println(sig)


	//通过公钥实现验签

	err:=rsa.VerifyPSS(pub,crypto.MD5,hashed,sig,&opts)
	if err==nil {
		fmt.Println("验签成功")
	}

}
