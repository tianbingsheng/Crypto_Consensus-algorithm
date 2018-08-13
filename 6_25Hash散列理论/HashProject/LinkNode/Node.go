package LinkNode

import "fmt"

//实现链表


//声明键值对类型
type KV struct {
	Key string
	Value string
}


//声明节点类型
type Node struct {
	Data KV
	NextNode *Node
}


//创建头节点，并通过尾插法，插入新节点
func CreateHead() *Node {
	var node = &Node{KV{"头节点","头节点"},nil}
	return node
}


//添加新节点
func AddNode(data KV,node *Node) *Node {
	var newNode = &Node{data,nil}
	node.NextNode = newNode
	return newNode
}


//遍历链表
func ShowNode(key string ,head *Node) {
	node:=head
	for {
		//当遍历指定下标下的链表时，判断打印只有node.Data.key == key的value即可
		if node.Data.Key == key {
			fmt.Println(node.Data.Value)
		}

		if node.NextNode == nil {
			break
		} else {
			node = node.NextNode
		}
	}
}



//遍历到尾节点
func GetTailNode(head *Node ) *Node {
	node:=head
	for {
		if node.NextNode == nil {
			return node
		} else {
			node =node.NextNode
		}
	}
}











