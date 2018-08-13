package main

import (
	"fmt"
	"net"
	"bufio"
	"io"
	"strconv"
)

//拜占庭理论的代码实现部分

type Node struct {
	Name string
	Status int//1 代表去，０代表不去
	Votes []*Node//记录账本

}


//保存４个node
var nodes = make([]*Node,0)


func createNodes() {

	//监听终端
	netListen,_:=net.Listen("tcp","127.0.0.1:1234")
	defer netListen.Close()

	//监听终端连接
	for {
		conn,_:=netListen.Accept()
		//扫描终端
		myscan:=bufio.NewScanner(conn)

		io.WriteString(conn,"input name")
		//对终端的内容监听
		go func() {
			//接收中的内容
			myscan.Scan()
			name:=myscan.Text()
			io.WriteString(conn,"input status")
			myscan.Scan()
			status:=myscan.Text()

			//fmt.Println(name,status)
			//创建新的node对象
			sts,_:=strconv.Atoi(status)
			node:=Node{name,sts,make([]*Node,0)}
			//将对象添加到数组
			nodes=append(nodes,&node)

			//投票
			votes()

		}()


	}



	/*A:=Node{"A",1,make([]*Node,0)}//去
	B:=Node{"B",1,make([]*Node,0)}//去
	C:=Node{"C",1,make([]*Node,0)}//去
	D:=Node{"D",0,make([]*Node,0)}//不去
	//按照拜占庭的1/3理论，这次打仗是可行的
	nodes=append(nodes,&A)
	nodes=append(nodes,&B)
	nodes=append(nodes,&C)
	nodes=append(nodes,&D)*/

}


//互相转达
func votes() {
	for i:=0;i<len(nodes);i++{
		node:=nodes[i]
		//可以将每个人的进攻状态获取出来
		fmt.Println(node.Status)

		//将此人的状态分发给其他人
		for j:=0;j<len(nodes);j++{
			inode := nodes[j]
			node.Votes=append(node.Votes,inode)
		}

	}


	//打印是否可行
	fmt.Println(isValid())

}


//判断本次进攻是否可行，判断叛徒是否小于三分之一
func isValid() bool{
	//在数组中取出最后一个对象
	node:=nodes[len(nodes)-1]
	votes:=node.Votes

	cnt:=0
	for _,n :=range votes {
		fmt.Println(n.Status)
		if n.Status == 0 {
			cnt++
		}
	}

	//判断cnt只有小于n/3的情况下，才能成功,拜占庭
	if float32(cnt)<float32(len(nodes))/float32(3.0) {

		return true
	}

	return false

}


func main() {

	createNodes()

}
