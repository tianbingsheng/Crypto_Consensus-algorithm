package main

import (
	"net"
	"fmt"
	"MyCrypt/CryptedDic"
)

func MyClient(cipherTxt []byte) {

	//创建准备链接的服务器
	netAddr,_:=net.ResolveTCPAddr("tcp4","127.0.0.1:1234")
	//链接服务器
	conn,_:=net.DialTCP("tcp",nil,netAddr)
	//发送数据
	conn.Write(cipherTxt)
	fmt.Println("发送数据成功")

}

func main() {


	data:=[]byte("hello 123")
	//cipherTxt就是加密出的密文
	cipherTxt:=CryptedDic.EnCrypt("12345678",data)
	//发送数据
	MyClient(cipherTxt)
}