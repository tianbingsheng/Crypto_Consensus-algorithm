package main

import "List/ListNode"

func main()  {

	ListNode.CreateHead(1)
	ListNode.AddNode(2)
	ListNode.AddNode(3)
	ListNode.AddNode(4)

	ListNode.PreShow()
	ListNode.BackShow()
}
