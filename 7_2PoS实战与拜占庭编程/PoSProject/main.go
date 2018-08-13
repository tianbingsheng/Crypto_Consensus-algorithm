package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"fmt"
	"net"
	"bufio"
	"io"
	"strconv"
	"math/rand"
)

//创建区块
type Block struct {
	Index int
	Timestamp string
	Prehash string
	Hash string
	Data int

	//终端地址
	Validator string
}


//计算某个字符串的hash
func calculateHash(record string)string {
	h:=sha256.New()
	h.Write([]byte(record))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//计算block的hash
func calculateBlockHash(block Block) string {
	record:=block.Timestamp+string(block.Data)+
		block.Prehash+string(block.Index)
		hashCode :=calculateHash(record)
		return hashCode
}

//生成新区块
var Blockchain []Block

func generateNextBlock(oldBlock Block,data int ,vald string ) Block {
	var newBlock Block
	//设置区块高度
	newBlock.Index = oldBlock.Index+1
	newBlock.Timestamp = time.Now().String()
	newBlock.Prehash = oldBlock.Hash
	newBlock.Data =data
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = vald
	//添加到区块链
	Blockchain=append(Blockchain,newBlock)
	return newBlock
}

//创建创世区块
func genesisBlock() Block{
	var genesisBlock =Block{0,time.Now().String(),"",
	"",0,""}
	//计算genesisBlock的hash值
	genesisBlock.Hash = calculateBlockHash(genesisBlock)
	//将创世区块添加到数组
	Blockchain=append(Blockchain,genesisBlock)
	return genesisBlock
}

//创建conn终端连接的数组
var connAddr []net.Conn




//创建节点类型
type Node struct {
	//终端的地址
	Address string
	//币领
	Coins int
}

//保存终端的对象
var nodes []Node

//通道实现线程通信
var announcements=make(chan string)

func main() {

	//测试代码
	//genesisBlock:=genesisBlock()
	//newBlock:=generateNextBlock(genesisBlock,1)
	//generateNextBlock(newBlock,2)
	//区块链中是否有三个区块
	//fmt.Println(Blockchain)

	genesisBlock:=genesisBlock()

	//如何通过终端链接代码上
	//1,如何mac　首先安装telnet命令　，　brew install telnet(1-2分钟左右）
	//2,创建监听
	netListen,_:=net.Listen("tcp","127.0.0.1:1234")
	defer netListen.Close()

	go func() {
		//通过旷工实现区块的挖矿
		for {
			//此代码会将for循环卡死,
			w:=<-announcements

			//将新的区块,利用w旷工，添加到数组中
			generateNextBlock(genesisBlock,100,w)

			//UTXO ,未花费设计模式

		}

	}()



	go func() {
		//每隔开10S选择一次旷工
		for {
			time.Sleep(10*time.Second)
			winner:=pickWinner()
			fmt.Println("系统通过PoS帮您选出的旷工为",winner)
			//将旷工放入到通到中
			announcements<-winner
		}

	}()

	//3,等待连接
	for {
		conn,_:=netListen.Accept()
		//将所有的链接保存到数组
		connAddr=append(connAddr,conn)
		//扫描终端
		scanbalance:=bufio.NewScanner(conn)

		io.WriteString(conn,"请输入币领")

		//扫描一下终端里边写了什么
		go func() {
			//在子线程中扫描终端中的内容
			for scanbalance.Scan(){
				txt:=scanbalance.Text()
				//打印终端输入的信息
				fmt.Println("您刚才从终端输入的币领为：",txt)

				//通过时间戳创建地址
				addr:=calculateHash(time.Now().String())
				cons ,_:=strconv.Atoi(txt)
				node:=Node{addr,cons}
				//将链接终端对象存放到数组
				nodes=append(nodes,node)
				fmt.Println(nodes)

			}
		}()


	}

}


//通过PoS共识算法选择旷工
func pickWinner() string {
	//选择旷工，利用PoS共识算法选择旷工
	var lottyPool []string

	//根据币领把对应的旷工地址，存放到数组中
	for i:=0;i<len(nodes);i++{
		node:=nodes[i]
		for j:=0;j<node.Coins;j++{
			lottyPool=append(lottyPool,node.Address)
		}
	}

	if len(lottyPool)!=0 {
		//通过随机值，找到准备挖矿的旷工
		rand.Seed(time.Now().Unix())
		r:=rand.Intn(len(lottyPool))
		workerAddress:=lottyPool[r]
		//返回旷工地址
		return workerAddress
	}

	return ""
}
