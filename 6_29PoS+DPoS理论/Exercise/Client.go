package main

import (
	"net"
	"fmt"
	"crypto/md5"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
	"encoding/pem"
	"crypto/x509"
	"math/big"
)


//私钥做签名
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyOnQ8Dbm1/UIkmSfeMKdK1LRJGX7T18vjZ7P4w3f/Jft/LKC
wkxyC2H7x03An+EdHP7dreRhNytzbQaseIgHEYjFapaCEz+JpMNm+qY4ZpApzvPv
qm/tut4T1J0HG33iiBqnyMJRZg8LjUXV2tEwfnHm5yCX36kOkN/YCW7ZbeO6aqw7
gMyvJDiLGIYgCy2Daqe1MH1RP91djrt6tWcfqVmUR+HxvJFkvUZZHqFUUZyJefNc
Y7JQDLSz5F22VB7ZLd9sSX38My353pNy4D19yeo5/54Z5AbSeRUMYJSFFbxwJzfe
wyVq2nV7EUJEj7lk0NmksB+S6w1a+a8cWydJ/QIDAQABAoIBAAQo+z+OE3eTRksp
tDee5/w2qcf0KKD7GpP3HtzXs7SaPL5Hv/df99iOfdUhogRtd9na2SI5oV2wE6LF
SZrxThwp1dSgKy9U2HfF6AL2oCJXh9YWLPc9fBGreYOkgLosAB3LV4ALrf3L//Q7
5vKx9CwaFarhfOOPr5KGYAXJ+syQqi3CjQrPGTLsoyYPB5oc5CA45eHIctoS90M3
cCRb5pu8vlbmeMUh9G9GMdjD3zuefndOBnwcpErLf2xPuM/Qav9LI7bP25UaZe1u
zuTm93AjAtjS9zTvyqbVx/xq7C+LA4EaEeBzxNuUPHAGEhuf4kQGOPl48XKM3aNk
lc/UoUECgYEA5vTg6lJKkyHvA5JJvOLSRqrGd220TvO0EPmdp3PUGSFWdldC1ev1
M42f8tbzBBeQJlIMBTeGi5Sa8QRVVZXWYmjKkcpDreZJCKz4dVPyeg93MRUhDA7J
8+2GSypKO+MpTty3WY7y0K0Lyk7381to7QTfqXzMc1d/Q/W2rqdrITECgYEA3rL3
4EzaxisRG9yaxu1ndGFcaTX9JUm+RTDPd5wecfE2mrSqfWwcXjsJ/BOqim1pwPQe
1/7e6PwyPqqd9It+E5f3zLwN5fvHISSjawU+sCLgpPY4YQvybf2cLsfyQrIQw1Ig
4Mo+DTBp4EAGYLszn/8yk7A6PIkev/+W22s1oo0CgYEArYriGpC9QrOj4t4glFMn
yDv4p7MCYrUS4BKA2VvaywtW6dOGADKubA+GWbYUo0WFZpYOPDlN20JN7r7f1FC1
6Axv8n62vkwlIuS+SYNL11YZrQMAcwwFykn3uDFN2JRH7N9C0oPshssQ6fLOs8lD
HZ6k5blF84GSuqE+pRxeDnECgYAagUJvN7ZyD34hysfBFVBS0aOmAf73CSfUJFMM
8U/OT98n2ojQFKXrsSsnDVAm2Y7UoDlri7IMGLgCLVxPVqrSuiuL+nXNAYJZt3qb
qiwj2oLSH1vmcP2RibWk+7chqP/Fv2iaWHe6KiDvx676pE0opb7nRPopakh2oXz1
8I+ZoQKBgDR/aXBDXcDzHC4dM3dYzkSf0bn8LXyANkEjHovSH/QSs4m+d2BkGlSy
yB3kgNSnEa9vNoffQcRvRlQbS58zTF8Z4iGjnoCHS6Q2yJBFm9L+EaRJlF6tOERk
ngLn8mAtV/IGigWBpZCVeEIHH1nG1DLatF2VDCQifQXZ5oRcZZr6
-----END RSA PRIVATE KEY-----
`)



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



//用私钥做签名
func SignData() []byte {

	//首先用hash做明文散列
	plaintxt:=[]byte("IAMA")
	h:=md5.New()
	h.Write(plaintxt)
	hashed:=h.Sum(nil)

	//将字节数组，转换成*private类型
	block ,_:=pem.Decode(privateKey)
	priv,_:=x509.ParsePKCS1PrivateKey(block.Bytes)

	//通过pss做签名
	opts:=&rsa.PSSOptions{rsa.PSSSaltLengthAuto,crypto.MD5}
	sig,_:=rsa.SignPSS(rand.Reader,priv,crypto.MD5,hashed,opts)

	return sig

}



//通过tcp发送数据

func Send( data []byte) {
	//创建准备链接的服务器
	conn,_:=net.ResolveTCPAddr("tcp4","127.0.0.1:1234")
	//开始链接
	n,_:=net.DialTCP("tcp",nil,conn)
	//发送数据
	n.Write(data)

	fmt.Println("发送结束")

}



func main() {

	//发送的收，将"hello baby"　和　sig　两个数组拼接到一起

	sg:=SignData()

	//明文
	message:=[]byte("hello kongyixueyuan hello world hello china")

	//将字节数组转换成publickey类型
	//公钥加密
	block ,_:= pem.Decode(publicKey)
	//解析公钥
	pubInterface ,_:=x509.ParsePKIXPublicKey(block.Bytes)
	//设置刚才公钥为public key 类型断言
	pub:= pubInterface.(*rsa.PublicKey)

	//用公钥对明文加密后，然后在网络上传输密文
	ciphertTxt,_:=rsa.EncryptOAEP(md5.New(),rand.Reader,pub,message,nil)


	var data= make([]byte ,len(ciphertTxt)+len(sg)+4)

	var l = len(ciphertTxt)
	var b = big.NewInt(int64(l))

	copy(data[0:2],b.Bytes())

	l = len(sg)
	b = big.NewInt(int64(l))

	copy(data[2:4],b.Bytes())
	copy(data[4:4+len(ciphertTxt)],ciphertTxt)
	copy(data[4+len(ciphertTxt):],sg)



	fmt.Println("发送的总数据为",data)
	fmt.Println("发送的内容为",[]byte("hello baby"))
	fmt.Println("发送的签名数据为",sg)

	Send(data)


	//优化，要求可以发送任意长度的明文


}