package main

import (
	"net"

	"MyCrypt/CryptedDic"
	"fmt"
)

func MyServer() []byte {
	//监听本机1234端口
	netListen,_:=net.Listen("tcp","127.0.0.1:1234")
	//延时关闭
	defer netListen.Close()
	for  {
		//等待连接
		conn,_:=netListen.Accept()

		//接收数据
		data:=make([]byte,2048)
		for {
			//代表对方发送数据的长度
			n,_:=conn.Read(data)
			//fmt.Println("您接收到了对方的数据为：",string(data[:n]))
			//返回接收到的密文
			return data[:n]
			//break
		}


	}
}


func main() {


	cipherTxt:=MyServer()
	//解密密文
	orig:=CryptedDic.Decrypt(cipherTxt,"12345678")
	fmt.Println("解密后的数据为",string(orig))
}
