package main

import (
	"time"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
)

//实现PoS挖矿的原理

type Block struct {
	 Data int //交易记录
	 Prehash string
	 Hash string
	 Timestamp string
	 Index int
	 //记录挖矿的那个节点的地址
	 validator *Node

}


//创建创世区块
func genesisBlock() Block {
	var genesisBlock = Block{0,"0","",time.Now().String(),
	0,nil}
	calculateHash(&genesisBlock)
	return genesisBlock
}


//计算区块的hash值
func calculateHash(block *Block) [] byte {
	record:=strconv.Itoa(block.Data) + strconv.Itoa(block.Index)+block.Prehash+
		block.Timestamp
		h:=sha256.New()
		h.Write([]byte(record))
		hashed:=h.Sum(nil)
		block.Hash = hex.EncodeToString(hashed)
		return hashed
}


//创建全节点类型，可以理解成持有币的村民类型
type Node struct {
	Tokens int //持有币个数
	Days int   //持有币的时间
	Address string //地址
}

//创建５个村民
var n =make([]Node,5)
//存放每个村民的地址
var addr = make([]*Node,15)

func initNode()  {
	n[0] = Node{1,1,"1111"}
	n[1] = Node{2,1,"2222"}
	n[2] = Node{3,1,"3333"}
	n[3] = Node{4,1,"4444"}
	n[4] = Node{5,1,"5555"}

	cnt :=0
	for i:=0;i<5;i++ {
		for j:=0;j<n[i].Tokens*n[i].Days;j++{
			addr[cnt] = &n[i]
			cnt++
		}
	}

}




//采用PoS共识算法实现挖矿
func generateNextBlock(oldBlock Block,data int) Block {
	var newBlock Block
	newBlock.Index = oldBlock.Index+1
	newBlock.Timestamp = time.Now().String()
	newBlock.Prehash= oldBlock.Hash
	newBlock.Data= data
	calculateHash(&newBlock)

	//通过pos计算由那个村民挖矿
	//设置随机种子
	rand.Seed(time.Now().Unix())
	//[0,15)产生0-15的随机值
	var rd =rand.Intn(15)

	//选出挖矿的旷工
	node := addr[rd]
	//设置当前区块挖矿地址为旷工
	newBlock.validator = node
	//当前node 旷工的原有的币增加
	node.Tokens = node.Tokens +1


	return newBlock

}


func main () {

	initNode()

	//创建创世区块
	var genesisBlock = genesisBlock()

	//创建新区快
	var newBlock = generateNextBlock(genesisBlock,1)

	//打印新区快信息
	fmt.Println(newBlock.validator.Address)
	fmt.Println(newBlock.validator.Tokens)

}

