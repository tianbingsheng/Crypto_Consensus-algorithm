package main

import (
	"net/http"
	"fmt"
)

type RaftServer struct {

}
func main() {
	rs:=RaftServer{}
	rs.setHettServers()
}
//设置服务器
func (rs *RaftServer)setHettServers() {


	http.HandleFunc("/req",rs.request)
	if err:=http.ListenAndServe("127.0.0.1:5566",nil);err!=nil {
		fmt.Println(err)
	}

}

func (rs *RaftServer)request(writer http.ResponseWriter,request *http.Request) {
	request.ParseForm()
	//接收网页的消息
	if len(request.Form["data"])>0 {
		fmt.Println(request.Form["data"][0])
	}

}

