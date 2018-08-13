package main

import (
	"strconv"
	"crypto/sha256"
	"time"
	"encoding/hex"
	"fmt"
	"strings"
)

//在区块链中应用PoW挖矿算法

//定义难度系数
const difficulty =4

//创建节点类型
type Block struct {
	Index int
	Timestamp string
	Data int
	Hash string
	Prehash string
	Nonce int
	Difficulty int
}

//创建区块链，可以用数组或链表维护区块链，比特币和以太坊都是用的数组维护区块链
var Blockchain []Block


//创建生成新区块,PoW挖矿
func generatNextBlock(oldBlock Block,data int ) Block {
	var newBlock Block
	newBlock.Index = oldBlock.Index +1
	newBlock.Data= data
	newBlock.Difficulty = difficulty
	newBlock.Prehash = oldBlock.Hash
	newBlock.Nonce = 0
	//挖矿
	for {
		//使得当前区块的Hash值的前边0的个数与难度系数值相同
		newBlock.Timestamp= time.Now().String()
		//计算当前区块的hash值
		newBlockHash:=hex.EncodeToString(calculateHash(&newBlock))
		fmt.Println("挖矿中",newBlockHash)
		//判断挖矿结束
		if isBlackValid(newBlockHash,difficulty) {
			//则挖矿成功
			fmt.Println("挖矿成功")
			newBlock.Hash = newBlockHash
			//需要将newBlock校验后，然后将挖出的区块添加区块链
			if newBlockVerify(newBlock,oldBlock) {

				//可以添加到数组
				Blockchain=append(Blockchain,newBlock)
				return newBlock
			}

		}

		newBlock.Nonce++

	}
}

//计算区块的hash值
func calculateHash (block *Block) []byte {
	//利用当前区块数据，计算完成hash散列
	record:=strconv.Itoa(block.Nonce)+strconv.Itoa(block.Difficulty)+
		strconv.Itoa(block.Data)+strconv.Itoa(block.Index)+block.Prehash+
			block.Timestamp

			//通过sha256实现hash散列
			h:=sha256.New()
			h.Write([]byte(record))
			hashed:=h.Sum(nil)

			return hashed
}


//判断是否挖矿结束
func isBlackValid(hash string ,diff int) bool {
	hashPrefix:=strings.Repeat("0",diff)//"0000"
	return strings.HasPrefix(hash,hashPrefix)

}



//校验新的区块是否合法
func newBlockVerify(newblock Block ,oldoblock Block) bool {
	if oldoblock.Index +1 !=newblock.Index {
		return false
	}
	if newblock.Prehash !=oldoblock.Hash {
		return false
	}
	return true

}


//创建传世区块
func genesisBlock() Block {
	var genesisBlock = Block{0,time.Now().String(),0,"",
	"0",0,difficulty}
	//计算创世区块的hash
	calculateHash(&genesisBlock)

	//首先将创世区块放入到数组
	Blockchain=append(Blockchain,genesisBlock)

	return genesisBlock
}


func main() {

	var genesisBlock = genesisBlock()


	generatNextBlock(genesisBlock,1)


}
