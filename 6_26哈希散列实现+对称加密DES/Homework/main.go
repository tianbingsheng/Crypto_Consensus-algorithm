package main

import (

	"fmt"
	"time"
)

//保存头节点信息
var hNode *Node
//保存当前节点信息
var cNode *Node

//创建节点结构体
type Node struct {
	Data int
	NextNode *Node
}


//创建头节点
func CreateHead(data int)  *Node{
	var node =&Node{data,nil}
	hNode=node
	cNode = node
	return node
}

//添加新节点
func AddNode(data int) *Node {
	var newNode= &Node{data,nil}
	cNode.NextNode = newNode
	//改变当前节点为新添加的节点
	cNode = newNode
	return newNode
}



//遍历链表
func ShowNode() {
	node:=hNode

	cnt:=1

	for {

		fmt.Println(node.Data)

		//当循环为5时剔除某个节点
		if cnt==4 {
			fmt.Println("节点被剔除",node.NextNode.Data)
			//剔除4后边的那个节点
			node.NextNode = node.NextNode.NextNode

			cnt = 0
		}


		cnt++


		//每隔一秒打印一次
		time.Sleep(1*time.Second)
		if node.NextNode == nil {
			break
		} else {
			node = node.NextNode
		}
	}
}



func main() {

	head:=CreateHead(1)
	AddNode(2)
	AddNode(3)
	AddNode(4)
	AddNode(5)
	tail:=AddNode(6)

	//确保实现循环链表
	tail.NextNode = head


	ShowNode()

}
