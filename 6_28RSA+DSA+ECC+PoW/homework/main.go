package main

import (
	"bytes"
	"fmt"
)

//PKCS8补码
func PKCS8Padding(org []byte)[]byte {
	//长度就是８
	pad:=8-len(org)%8
	//补码的字节数组
	pading:=bytes.Repeat([]byte{byte(pad)},pad)

	//补码
	return append(org,pading...)



}

func PKCS8Unpadding(org []byte) []byte  {
	l:=len(org)
	pad:=int(org[l-1])
	return org[:l-pad]
}


//实现密文分组的过程
func gcrypt(org []byte) [][]byte {
	//因为利用pkcs8实现的补码,所以就以８字节为单位
	var g [][]byte
	var t []byte
	//遍历补码后的数据
	for i:=0;i<len(org);i++{
		t=append(t,org[i]+1)
		if i%7==0&&i!=0 {
			//做分组
			g=append(g,t)
			t=make([]byte,0)
		}
	}


	return g
}



func main() {

	cipher:=PKCS8Padding([]byte("helloworldaaa222b"))
	fmt.Println(cipher)

	//fmt.Println(PKCS8Unpadding(cipher))

	g:=gcrypt(cipher)

	//通过遍历二位数组，测试分组是否成功
	for index,value :=range g {
		for _,val :=range value {
			fmt.Printf("当前组%d,值为%d\n",index,val-1)
		}
	}

}
