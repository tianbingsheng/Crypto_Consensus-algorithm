package main

import (
	"fmt"
	"Nodes/LinkNodes"
)

func main() {

	fmt.Println("hello world")

	//调用CreateHeadNode,创建头节点
	head:=LinkNodes.CreateHeadNode(1)
	//添加新节点
	node:=LinkNodes.AddNode(2,head)
	node = LinkNodes.AddNode(3,node)
	node = LinkNodes.AddNode(4,node)

	//head = LinkNodes.InsertNodeWithIndex(100,0,head)

	//测试插入尾节点
	//LinkNodes.InsertNodeWithIndex(100,3,head)


	//测试修改节点信息
	//LinkNodes.UpdateNodeByIndex(100,3,head)


	//测试删除节点
	head=LinkNodes.DeleteNodeByIndex(0,head)

	//测试便利节点
	LinkNodes.ShowNodes(head)

	//测试链表总长度
	//fmt.Println("链表总长度为：",LinkNodes.NLen(head))

}

