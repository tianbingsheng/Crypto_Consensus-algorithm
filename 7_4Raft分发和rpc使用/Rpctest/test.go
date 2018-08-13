package main

import (
	"net/rpc"
	"fmt"
)

//如何调用rpc

type Param struct {
	Width,Height int
}



func main() {

	//调用rpc 服务器
	rp,err:=rpc.DialHTTP("tcp","127.0.0.1:9000")
	if err!=nil {
		fmt.Println(err)
	}

	ret:=0
	e:=rp.Call("Rect.Aera",Param{100,100},&ret)
	if e!=nil {
		fmt.Println(e)
	}
	fmt.Println(ret)


}