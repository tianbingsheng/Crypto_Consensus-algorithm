package main

import (
	"net"
	"fmt"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"crypto"
	"crypto/md5"
)

//公钥验签
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyOnQ8Dbm1/UIkmSfeMKd
K1LRJGX7T18vjZ7P4w3f/Jft/LKCwkxyC2H7x03An+EdHP7dreRhNytzbQaseIgH
EYjFapaCEz+JpMNm+qY4ZpApzvPvqm/tut4T1J0HG33iiBqnyMJRZg8LjUXV2tEw
fnHm5yCX36kOkN/YCW7ZbeO6aqw7gMyvJDiLGIYgCy2Daqe1MH1RP91djrt6tWcf
qVmUR+HxvJFkvUZZHqFUUZyJefNcY7JQDLSz5F22VB7ZLd9sSX38My353pNy4D19
yeo5/54Z5AbSeRUMYJSFFbxwJzfewyVq2nV7EUJEj7lk0NmksB+S6w1a+a8cWydJ
/QIDAQAB
-----END PUBLIC KEY-----
`)

//接收数据
func Recive() []byte {
	//创建tcp
	netListen, _ := net.Listen("tcp", "127.0.0.1:1234")
	defer netListen.Close()

	for {
		//等待链接
		conn, _ := netListen.Accept()

		//创建数组，保存接收到的数据
		data := make([]byte, 2048)
		for {
			n, _ := conn.Read(data)
			return data[:n]
		}

	}
}




func main() {

	data:=Recive()

	//首先切割数组
	plaintxt:=data[:10]
	fmt.Println("收到的数据为：",plaintxt)

	//用公钥做验签工作
	sig := data[10:]

	//将字节数组转换成publickey类型
	//公钥加密
	block ,_:= pem.Decode(publicKey)
	//解析公钥
	pubInterface ,_:=x509.ParsePKIXPublicKey(block.Bytes)
	//设置刚才公钥为public key 类型断言
	pub:= pubInterface.(*rsa.PublicKey)


	//在此应该和”IAMA“做校验
	h:=md5.New()
	h.Write([]byte("IAMA"))
	hashed:=h.Sum(nil)


	e:=rsa.VerifyPSS(pub,crypto.MD5,hashed,sig,nil)

	if e==nil {
		fmt.Println("验签成功，您可放心收接收到的",string(plaintxt))
	}



}




