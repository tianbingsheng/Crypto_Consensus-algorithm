package main

import (
	"time"
	"encoding/hex"
	"crypto/sha256"
	"fmt"
	"strings"
)

//完成比特币中PoW挖矿的共识算法

//创建区块的结构体


//声明区块，挖矿就是挖区块
//挖区块的方式，常用的有三种PoW（比特币）,PoS,DPoS方式
//编程实现通过PoW实现挖矿的过程
type Block struct {

	//上一个区块的hash值
	PreHash []byte
	//时间戳
	Timestamp int64
	//交易信息
	Data []byte
	//当前区块的hash值
	Hash []byte
	//随机数
	Nonce int

}


//创建创世区块
func CreateGenesisBlock() *Block {
	var genesisBlock = &Block{[]byte{0},time.Now().Unix(),
	[]byte("ab交易1bitcoin"),nil,0}
	//计算当前区块的hash
	genesisBlock.getBlockHash()
	return genesisBlock
}


//计算block的当前hash值
func(block *Block)getBlockHash() []byte {
	//拼接区块信息
	var blockInfo = hex.EncodeToString(block.Data)+string(block.Nonce)+
		string(block.Timestamp)+hex.EncodeToString(block.PreHash)

		h:=sha256.New()
		h.Write([]byte(blockInfo))
		hashed:=h.Sum(nil)

		//设置计算出来的hash值，给当前的block
		block.Hash = hashed
	return hashed

}

//通过PoW挖矿的方式挖新区块
func GenerateNextBlock(oldBlock *Block) *Block {
	//假设难度系数为4,则挖的区块的hash前边必须有4个0才算挖矿成功
	var newBlock = &Block{oldBlock.Hash,time.Now().Unix(),
	[]byte("bc交易"),nil,0}
	//不断改变nonce值，最终实现当前区块的hash的0的个数与系统中要求的难度系数值一致
	nonce:=1

	for {

		var blockInfo = hex.EncodeToString(newBlock.PreHash)+hex.EncodeToString(newBlock.Data)+
			string(nonce)+string(newBlock.Timestamp)

			h:=sha256.New()
			h.Write([]byte(blockInfo))
			hashed:=h.Sum(nil)
			hashString :=hex.EncodeToString(hashed)
			fmt.Println("挖矿中",hashString)
			if strings.HasPrefix(hashString,"0000") {
				fmt.Println("挖矿成功")

				newBlock.Hash = hashed
				return  newBlock
			}
			nonce++

	}

}



func main () {

	//创建创世区块
	var genesisBlock = CreateGenesisBlock()

	var newBlock = GenerateNextBlock(genesisBlock)

	fmt.Println("挖出的新区快Data为",newBlock.Data)
	fmt.Println("挖出的新区快Data为",hex.EncodeToString(newBlock.Hash))

}


