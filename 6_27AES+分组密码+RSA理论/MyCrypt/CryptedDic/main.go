package CryptedDic

import (
	"bytes"

)


//自行完成对陈加密的算法，实现通过密码对明文加密，同时也可以用秘钥对密文进行解密


//创建一个加密的方法
func EnCrypt(key string ,data []byte) [] byte {
	//加密算法：首先计算key的总和,利用key的总和与明文参与运算

	var sum = 0
	for i:=0;i<len(key);i++ {
		sum = sum +int(key[i])
	}

	//首先对明文进行补码
	var pad = PKCS5Padding(data,len(key))

	//通过加法运算，实现加密过程
	for i:=0;i<len(pad);i++{
		pad[i]= pad[i]+byte(sum)
	}

	return pad

}


//创建一个解密的方法
//解密的过程是加密过程的逆运算
func Decrypt(cipherTxt []byte,key string) []byte {

	//计算key的总和
	var sum =0
	for i:=0;i<len(key);i++ {
		sum += int(key[i])
	}

	//减法运算
	for i:=0;i<len(cipherTxt);i++{
		cipherTxt[i]=cipherTxt[i]-byte(sum)
	}

	//去码
	var p = PKCS5UnPadding(cipherTxt)


	return p
}


//PKCS5Unpadding 去码
func PKCS5UnPadding(cipherTxt []byte) []byte {
	var l = len(cipherTxt)
	var txt = int(cipherTxt[l-1])
	return cipherTxt[:l-txt]
}


//实现PKCS5Padding补码
func PKCS5Padding(cipherTxt[] byte,blockSize int) []byte {
	//计算准备添加的数字
	padding:=blockSize-len(cipherTxt)%blockSize

	//55555
	padTxt:=bytes.Repeat([]byte{byte(padding)},padding)

	//叠加两个数组
	var byteTxt =append(cipherTxt,padTxt...)

	return byteTxt

}

