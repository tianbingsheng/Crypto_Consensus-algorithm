package main

import (
	"os"
	"fmt"
	"net/http"
	"io"
)

//创建存储节点的对象
//存放四个国家的地址
var nodeTable = make(map[string]string)

//声明节点类型
type Node struct {
	NodeId string
	Path string
	Writer http.ResponseWriter
}

func main() {



	//接收终端传的参数
	nodeId :=os.Args[1]
	fmt.Println("打印终端参数",nodeId)

	//创建四个国家的地址
	nodeTable = map[string]string {
		"A":"10.0.151.215:1111",
		"B":"10.0.151.235:1112",
		"C":"10.0.151.238:1113",
	}



	//创建节点对象
	node:=Node{nodeId, nodeTable[nodeId],nil }

	//对node对象的http的请求做监听
	http.HandleFunc("/req",node.request)
	//处理接收到的数据
	http.HandleFunc("/prepare",node.prepare)

	//处理ＢＣ返回来的确认消息
	http.HandleFunc("/commit",node.commit)



	//http://localhost:1111/req?data=1234

	//启动http服务
	if err:=http.ListenAndServe(node.Path,nil );err!=nil {
		fmt.Println(err)
	}


}

//从浏览器向Ａ对象传参数
func (node *Node)request(writer http.ResponseWriter,reqeust *http.Request){
	//设置能够解析参数
	reqeust.ParseForm()
	//解析参数
	if len(reqeust.Form["data"])>0 {
		fmt.Println(reqeust.Form["data"][0])
		//writer,就是conn
		node.Writer = writer
		//将收到的消息转发给另外两个人
		node.broadcast(reqeust.Form["data"][0],"/prepare")
	}
}

//转发函数,A转发数据
func(node *Node) broadcast(data string ,path string ){
	for nodeId,url :=range nodeTable {
		//设置不能自己给自己转发
		if nodeId==node.NodeId {
			continue
		}

		//给另外除了自己之外的人做转发
		//http://127.0.0.1:1112/prepare?data=data&nodeid=A
		http.Get("http://"+url+path+"?data="+data+"&nodeId="+node.NodeId)
	}
}

//Ｂ,C　在接收数据
func(node *Node)prepare(writer http.ResponseWriter,request *http.Request){

	request.ParseForm()
	if len(request.Form["data"])>0 && len(request.Form["nodeId"])>0{
		fmt.Println("data",request.Form["data"][0])
		fmt.Println("nodeId",request.Form["nodeId"][0])

		//当接收到A转发过来的数据，然后在通知给A,告诉Ａ我已经收到消息了
		var path = nodeTable[request.Form["nodeId"][0]]

		var url ="http://"+path+"/commit?data="+"IRecieved"+"&nodeId="+node.NodeId

		http.Get(url)
	}

}

var cnt = 0
//编写Ａ接收到确认返回的消息后，处理一下
func(node *Node)commit(writer http.ResponseWriter,request *http.Request){


	request.ParseForm()
	fmt.Println("commit function")
	if len(request.Form["data"])>0 {
		fmt.Println(request.Form["data"][0])
		//如何接收到了两次IRecevied则，反馈给浏览器本次分发成功
		cnt++
	}
	if cnt ==2 {
		io.WriteString(node.Writer,"本次分发成功了")
	}

}