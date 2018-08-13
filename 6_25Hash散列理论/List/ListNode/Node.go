package ListNode

import (

	"fmt"
)

type Node struct {
	PreNode *Node
	Data int
	NextNode *Node
}

//保存头节点
var HNode *Node
//保存当前节点
var CNode *Node
//保存尾节点
var TNode *Node

//创建头节点
func CreateHead(data int) *Node {
	var node=&Node{nil,data,nil}
	HNode = node
	CNode = node
	TNode = node
	return node
}

//添加新节点
func AddNode(data int ) {
	var newNode = &Node{nil,data,nil}
	newNode.PreNode = CNode
	CNode.NextNode = newNode

	CNode = newNode
	TNode = newNode

}


func PreShow() {

	node :=HNode
	for {
		fmt.Println(node.Data)
		if node.NextNode == nil {
			break
		}else {
			//向后移动
			node= node.NextNode
		}
	}
}


func BackShow() {
	node :=TNode
	for {
		fmt.Println(node.Data)
		if node.PreNode== nil {
			break
		}else {
			//向前移动
			node = node.PreNode
		}
	}
}