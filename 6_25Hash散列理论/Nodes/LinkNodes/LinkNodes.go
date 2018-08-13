package LinkNodes

import (

	"fmt"

)

//通过代码完成链表的常用功能


//通过全局变量记录当前链表的头节点
var hNode *Node


//创建节点类型
type Node struct {
	//数据域
	Data int
	//地址域
	NextNode *Node
}


//创建头节点
func CreateHeadNode(data int) *Node{
	//创建Node对象，返回该对象地址
	var node = &Node{data,nil}
	//保存头节点
	hNode = node
	return node
}

//通过尾插法添加新节点
func AddNode(data int,node *Node) *Node {
	var newNode = &Node{data,nil}
	node.NextNode = newNode
	return  newNode
}


//链表的便利
func ShowNodes(head *Node) {
	//接收参数中的头节点
	node :=head
	for {
		//打印节点的数据域信息
		fmt.Println(node.Data)
		//如果位移到尾节点，则跳出循环，停止打印
		if node.NextNode == nil {
			break
		} else  {
			//实现节点的位移
			node = node.NextNode
		}
	}

}


//计算链表总长度
func NLen(head *Node) int{
	node:=head
	//通过cnt实现节点个数累加
	var cnt = 1
	for {
		if node.NextNode == nil {
			break
		} else {
			//控制节点向后位移
			node = node.NextNode
			cnt++
		}
	}
	return cnt
}



//按照下标插入新节点
func InsertNodeWithIndex(data int,index int,node *Node) *Node {

	//如何在下标为0的位置插入新节点，相当于插入了新的头结点
	if index == 0 {
		var insertedNode = &Node{data,nil}
		insertedNode.NextNode = node
		//保存头节点信息
		hNode = insertedNode
		//返回新的头节点
		return insertedNode

	} else if index >= NLen(hNode)-1 {
		//当插入的节点下标大于等于节点长度-1时，就相当于在链表上插入了一个新节点

		//1,获得到尾节点
		node := hNode
		for {
			if node.NextNode == nil {
				//node就是尾节点
				//2,在尾节点出插入新节点
				var insertedNode = &Node{data,nil}
				node.NextNode = insertedNode

				break
			} else {
				node =node.NextNode
			}

		}


	} else {
		//在链表中间的位置插入新节点
		head :=node
		var cnt = 1
		for {
			if cnt == index {
				//创建插入的新节点
				var insertedNode = &Node{data,nil}
				//修改链的方向
				insertedNode.NextNode = head.NextNode
				head.NextNode = insertedNode
				break
			} else {
				head = head.NextNode
				cnt++
			}
		}
	}

	return nil
}

//修改指定下标上的节点数据
func UpdateNodeByIndex(data int ,index int ,head *Node)  {


	if index == 0 {
		//修改头节点的信息
		hNode.Data = data
	} else {

		//通过遍历找到指定下标的节点
		node := head
		var cnt = 1
		for {
			if node.NextNode == nil {
				break
			} else {
				if cnt == index {
					//修改指定下标下的节点信息
					node.NextNode.Data = data

					break
				} else {
					node = node.NextNode
					cnt++
				}
			}
		}

	}
}



//删除指定下标的节点
func DeleteNodeByIndex(index int,head *Node) *Node {

	if index == 0 {
		//相当于准备删除头节点
		//让下一个节点当做头节点即可
		node:= head
		hNode = node.NextNode
		return hNode


	} else {

		node := head
		//通过遍历查询到准备删除的节点
		var cnt = 1

		for {
			if node.NextNode == nil {
				break
			} else {
				if cnt == index {
					//删除这个节点
					node.NextNode = node.NextNode.NextNode
					break

				} else {
					//节点位移，实现寻找index下标下的节点
					node = node.NextNode
					cnt++
				}
			}
		}
		return nil
	}


}

