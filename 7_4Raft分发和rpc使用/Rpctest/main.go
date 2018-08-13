package main

import (
	"net/rpc"
	"net/http"
	"fmt"
)

//raft库中默认使用rpc通讯
//rpc在p2p中使用非常广泛，学习成本底

//rpc对tcp的封装

//创建服务器
//声明传参类型
type Params struct {
	Width,Height int
}

//声明矩形对象类型
type Rect struct {

}

//计算Params的周长
func (r *Rect)Permiter(p Params,ret *int ) error{
	*ret = (p.Width+p.Height)*2
	return nil
}



//计算Params的面积
func (r *Rect)Aera(p Params,ret *int) error {
	*ret = p.Width *p.Height
	return  nil

}



func main() {
	//注册服务
	rect:=new(Rect)
	rpc.Register(rect)
	rpc.HandleHTTP()
	err:=http.ListenAndServe(":9000",nil)

	if err!=nil {
		fmt.Println(err)
	}


}
