package main

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"crypto/sha256"
	"os"
	"io"
	"golang.org/x/crypto/ripemd160"
)

//Md5
func MyMd5() {
	//测试Md5的编写方法

	//Md5第一写法
	//准备加密的明文
	data:=[]byte("hello world")
	//用Md5加密
	s:=fmt.Sprintf("%x",md5.Sum(data))
	//打印密文

	//密文为十六进制的数字
	//Md5加密的密文占16字节 = （16*8 = 128位）
	fmt.Println(s)



	//Md5第二种写法
	data2:=[]byte("hello world")
	m:=md5.New()
	m.Write(data2)
	//字节数组转换成字符串
	s2:=hex.EncodeToString(m.Sum(nil))
	fmt.Println(s2)



	//openssl是密码学最普遍使用的三方库，在此三方库中集成Md5,Sha256,RSA等典型加密算法



}


//测试sha256加密算法的使用
func MySha256() {

	//三种用法
	//第一种用法
	data:=[]byte("hello world")
	s:=fmt.Sprintf("%x",sha256.Sum256(data))
	fmt.Println(s)


	//第二种写法
	data2:=[]byte("hello world")
	m:=sha256.New()
	m.Write(data2)
	fmt.Println(hex.EncodeToString(m.Sum(nil)))

	//总结，使用sha256方式加密，占用32字节 = （32字节×8=256位），Sha256加密方法，通用在公链中


	//第三种用法
	//可将文件中的数据进行sha256加密处理
	//通过文件流首先找到文件,将文件读入内存
	f,_:=os.Open("test")
	//创建sha256对象
	h:=sha256.New()
	//将f内存中的数据，拷贝到sha256中
	io.Copy(h,f)
	//实现sha256计算过程
	s2:=h.Sum(nil)
	fmt.Println(hex.EncodeToString(s2))



}


//如果利用ripemd160加密，需要引入三方库
//引入三方法库的步骤
//1,进入gopath下，创建golang.org目录
//2,进入golang.org，创建x目录
//3,进入x目录，并在翻墙情况下,在github上下载三方库
//git clone https://github.com/golang/crypto.git


//以上的三个步骤可以通过一行命令在终端直接实现
//cd $GOPATH/src $ mkdir golang.org $ cd golang.org $ mkdir x $ cd x $ git clone https://github.com/golang/crypto.git


func MyRipemd160() {
	//只有一种写法
	hasher:=ripemd160.New()
	hasher.Write([]byte("hello world"))
	fmt.Println(hex.EncodeToString(hasher.Sum(nil)))


}


func main() {

	//MyMd5()
	//MySha256()
	MyRipemd160()
}
