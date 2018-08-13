package main

import (
	"net"
	"fmt"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"crypto"
	"crypto/md5"
	"math/big"
	"crypto/rand"
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


	fmt.Println("接收到的完整数据",data)

	ct:=data[0:2]
	var ctInt = big.NewInt(0)
	ctInt.SetBytes(ct)
	fmt.Println("收到的密文的长度为",ctInt)


	//获得密文
	cipherTxt:=data[4:4+ctInt.Int64()]

	//用私钥做解密
	//将字节数组，转换成*private类型
	block ,_:=pem.Decode(privateKey)
	priv,_:=x509.ParsePKCS1PrivateKey(block.Bytes)
	plaintxt,_:=rsa.DecryptOAEP(md5.New(),rand.Reader,priv,cipherTxt,nil)
	fmt.Println("解密后的明文为",string(plaintxt))


	ct=data[2:4]
	ctInt = big.NewInt(0)
	ctInt.SetBytes(ct)
	fmt.Println("收到的签名的长度为",ctInt)



	//用公钥做验签工作
	sig := data[4+ctInt.Int64():]

	//将字节数组转换成publickey类型
	//公钥加密
	block ,_= pem.Decode(publicKey)
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




