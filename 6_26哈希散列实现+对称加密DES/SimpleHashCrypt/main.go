package main

import "fmt"

//优秀的Hash散列应该具备以下特征
//1,正向快速
//2,逆向困难
//3,输入敏感
//4,避免冲突

func HashCode(data int ) int {
	//秘密学加密，一般都通过左移、右移，取模，异或来完成
	return (data+2+(data<<1))%8^5
}

func main(){

	//通过一个简单的算法，学习如何尽量实现Hash以上特征
	fmt.Println(HashCode(233))
	fmt.Println(HashCode(234))
	fmt.Println(HashCode(235))


}
