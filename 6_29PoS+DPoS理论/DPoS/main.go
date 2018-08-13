package main

import (
	"math/rand"
	"time"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)





//实现DPoS原理
var Blockchain [] *Block



//选举
type Node struct {
	Name string //节点名字
	Votes int // 被选举的票数
}


//创建数组，保存所有的节点
var n = make([]*Node ,5)


//创建节点
func createNodes() {

	//创建随机种子
	rand.Seed(time.Now().Unix())
	node1:=Node{"node1",rand.Intn(10)}
	node2:=Node{"node2",rand.Intn(10)}
	node3:=Node{"node3",rand.Intn(10)}
	node4:=Node{"node4",rand.Intn(10)}
	node5:=Node{"node5",rand.Intn(10)}
	n[0] = &node1
	n[1] = &node2
	n[2] = &node3
	n[3] = &node4
	n[4] = &node5



}


//DPoS中选出票数最高的前n位
func sortNodes() []*Node {
	//对所有选民的票数进行排序
	for i:=0;i<4;i++{
		for j:=0;j<4-i;j++ {
			if n[j].Votes<n[j+1].Votes {
				//二者的位置交换
				t:=n[j]
				n[j]=n[j+1]
				n[j+1]=t
			}

		}
	}

	return n[:3]

}


type Block struct {
	Index int
	Timestamp string
	Prehash string
	Hash string
	Data int

	//增加代理
	delegate *Node

}

func generateNextBlock (oldBlock Block,data int) *Block {
	var newBlock = Block{oldBlock.Index +1,time.Now().String(),
	oldBlock.Hash,"",data ,&Node{"",0}}
	calculateHash(&newBlock)
	//将新的区块添加到数组
	Blockchain=append(Blockchain ,&newBlock)
	return &newBlock
}



//计算新区块的hash
func calculateHash(block * Block) {
	record:=strconv.Itoa(block.Index)+strconv.Itoa(block.Data)+block.Timestamp+
		block.Prehash
		h:=sha256.New()
		h.Write([]byte(record))
		hashed:=h.Sum(nil)
		block.Hash = hex.EncodeToString(hashed)

}

//设置代理人的方法
func (block *Block)setDelete(node *Node) {
	block.delegate = node
}


//创世区块
func genesisBlock () Block {
	genesis:=Block{0,time.Now().String(),"","",0,&Node{"",0}}
	calculateHash(&genesis)
	Blockchain=append(Blockchain,&genesis)
	return genesis
}


func main() {

	//创建所有选民
	createNodes()

	//通过投票选出指定的人员做旷工
	c:=sortNodes()


	//for _,v :=range  c {
	//	fmt.Println(v.Votes)
	//}

	//c数组中就是被选民选出的3个主节点

	//下边所有的挖矿均有这三个主节点完成，三个人轮流挖矿
	g:=genesisBlock()

	//创建新区块
	newBlock :=generateNextBlock(g,1)
	//设置设置选取的第一位为旷工
	newBlock.setDelete(c[0])

	newBlock = generateNextBlock(*newBlock,2)
	newBlock.setDelete(c[1])


	newBlock = generateNextBlock(*newBlock,3)
	newBlock.setDelete(c[2])


	newBlock = generateNextBlock(*newBlock,4)
	newBlock.setDelete(c[0])

	//循环遍历当前区块链中每个区块的信息


	for _,blockInfo:=range Blockchain {
		var block = *blockInfo
		fmt.Println(block.delegate.Name)
		fmt.Println(block.Hash)
		fmt.Println(block.Data)
	}

}


