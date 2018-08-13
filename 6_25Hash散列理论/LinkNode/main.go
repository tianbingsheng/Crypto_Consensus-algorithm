package main

import "LinkNode/LinkNode"

func main() {


	LinkNode.CreateHead(1)
	LinkNode.AddNode(2)
	LinkNode.AddNode(3)
	LinkNode.AddNode(4)

	//测试链表的遍历
	LinkNode.ShowNode()

}
