package main

import (
	"fmt"
	"flag"
	"math/rand"
	"net"
	"strconv"
	"time"
	"strings"
	"net/http"
)

const (
	LEADER = iota
	CANDIDATE
	FOLLOWER
)

//声明地址信息
type Addr struct {
	Host string //IP
	Port int    //Port
	Addr string //地址描述
}

//实现Raft算法中选举的过程
type RaftServer struct {
	//选票
	Votes int
	//角色 ,follower ,candidate ,leader 三种角色
	Role int
	//存放每个节点的地址信息
	Node []Addr
	//判断当前节点现在是否处于选举中
	isElecting bool
	//设置选举的间隔时间，也为超时时间
	Timeout int
	//端口号
	Port int

	//设计通道实现接收到数网页数据的传参
	CostomerMsg chan string

	//设置两个通道变量
	ElectChan chan bool
	//控制leader的心跳信号
	HeartBeatChan chan bool

}

//设置节点对象的状态
func (rs * RaftServer)changeRole(role int) {
	switch role {
	case LEADER:
		fmt.Println("I Become Leader")
	case CANDIDATE:
		fmt.Println("I Become Candidate")
	case FOLLOWER:
		fmt.Println("I Become Follower")
	}
	rs.Role = role
}

func (rs * RaftServer)resetTIMEOut()  {
	//Raft系统中一般为1500-3000MilliSecond选一次
	rs.Timeout=  1500+rand.Intn(1500)

}

func main() {

	//获取参数的方法
	port:=flag.Int("p",5000,"port")
	flag.Parse()
	fmt.Println("您传入的端口号",*port)



	//创建新对象
	rs:=RaftServer{}

	rs.CostomerMsg = make(chan string)

	//监听http协议
	go rs.setHettServers()


	//设置默认票数
	rs.Votes = 0
	//3
	rs.isElecting = true
	//当创建一个节点对象，首先就是Follower状态
	rs.Role= FOLLOWER
	//ElectChan控制是否开始投票
	rs.ElectChan = make(chan bool)
	//心跳信号
	rs.HeartBeatChan = make(chan bool)
	rs.resetTIMEOut()
	//设置节点数组
	rs.Node = []Addr{
		{"127.0.0.1",5000,"5000"},
		{"127.0.0.1",5001,"5001"},
	}

	//设置服务端口
	rs.Port= *port
	//运行rs
	rs.Run()

	for{;}
}

//运行服务器程序
func(rs *RaftServer)Run() {
	//rs监听是否有别分给我投票了
	//通过tcp协议监听
	netListen,_:=net.Listen("tcp",":"+strconv.Itoa(rs.Port))

	//给别人投票
	go rs.elect()
	//控制投票的间隔时间
	go rs.electTimeDuration()

	//每个一秒打印一次当前对象的状态
	go rs.printRole()


	//设置头结点可以发送心跳
	go rs.sendHeartBeat()

	//主节点给其他子节点发送数据
	go rs.sendDataToOtherNodes()




	for{
		//等待别人的链接
		conn,_:=netListen.Accept()
		//监听别人发送的消息
		go func() {
			for {
				bts:=make([]byte,1024)
				n,_:=conn.Read(bts)
				fmt.Println("我收到了别人的投票",string(bts[:n]))
				if strings.HasPrefix( string(bts[:n]),"IVote1") {
					//就说明这包数据是投票数据，就是别人给我投票了
					rs.Votes++
					data:="我是"+strconv.Itoa(rs.Port)+" ,我当前的票数为"+strconv.Itoa(rs.Votes)
					fmt.Println(data)
					//判断如何票数等于指定的值，则leader选举成功
					if VotesSuccess(rs.Votes,len(rs.Node)/2) {
						msg:="我是"+strconv.Itoa(rs.Port)+"我被选举成了leader,我非常荣幸"
						fmt.Println(msg)
						//当前选leader成功了，当前的rs 则改为leader类型
						//１，刚创建的对象Follow
						//２，处于选举状态为Candidate
						//３，谁选中了leader，则为leader,其他的节点全部退回Follow

						//通知其他节点，停止选举工作,并且其他节点应退回follow状态
						rs.writesToOthers("StopVotes")
						rs.isElecting = false
						//修改了自己状态为leader
						rs.changeRole(LEADER)

					}
				}

				//就是接收到了leader　发来的停止投票命令
				if strings.HasPrefix(string(bts[:n]),"StopVotes"){

					//1停止给别人投票
					rs.isElecting = false
					//设置自己为Follow状态
					rs.changeRole(FOLLOWER)

				}
			}

		}()
	}
}

func VotesSuccess(votes int ,target int ) bool {
	if votes ==target {
		return true
	}
	return false
}



//发送数据
func (rs *RaftServer)writesToOthers(data string ){
	//给别人发送数据
	for _,k:=range rs.Node {
		if k.Port!=rs.Port {

			netAddr,err:=net.ResolveTCPAddr("tcp4",":"+strconv.Itoa(k.Port))
			if err!=nil {
				fmt.Println("您想连接的服务器没有运行")
			}else {
				conn,er:=net.DialTCP("tcp",nil,netAddr)
				if er !=nil {
					fmt.Println("DialTCP链接错误")
				}else {
					data = data +" 发送者为 "+strconv.Itoa(rs.Port)
					conn.Write([]byte(data))
				}

			}

		}
	}

}


//给别人投票
func (rs * RaftServer)elect() {
	for {
		//通过通道，确定现在可以投票，才给别人投票
		<-rs.ElectChan
		rs.writesToOthers("IVote1")
		//rs.ElectChan<-false

		if rs.Role != LEADER {
			//设置选举状态为Candidate
			rs.changeRole(CANDIDATE)
		}else {
			//是Leader的情况下

		}


	}
}

//时间选举时间间隔
func (rs * RaftServer)electTimeDuration() {
	for {
		//2
		if rs.isElecting {
			time.Sleep(time.Duration(rs.Timeout)*time.Millisecond)
			rs.ElectChan<-true
		}


	}

}


//打印当前对象的角色
func (rs *RaftServer)printRole() {
	for {
		time.Sleep(1*time.Second)
		fmt.Println(strconv.Itoa(rs.Port)+"  状态为  ",rs.Role)
	}
}


//主节点发送心跳信号给其他子节点
func (rs *RaftServer)sendHeartBeat() {
	//每个１秒发送一次心跳
	for {
		time.Sleep(1 * time.Second)
		if rs.Role == LEADER {
			//发送消息
			rs.writesToOthers("heart beating xintiaoxinhao......")
		}
	}
}


//通过leader给其他所有子节点发送数据
func (rs *RaftServer)sendDataToOtherNodes() {

	for {
		//每个10秒由主节点给子节点发送一包数据
		msg:=<-rs.CostomerMsg
		if rs.Role == LEADER{

			rs.writesToOthers(msg)
			//创建http协议的服务器，实现通过网页将数据创给leader，然后由leader转发给子节点，确保数据的一致性

		}

	}

}


//设置服务器
func (rs *RaftServer)setHettServers() {


	http.HandleFunc("/req",rs.request)
	if err:=http.ListenAndServe("127.0.0.1:12345",nil);err!=nil {
		fmt.Println(err)
	}

}

func (rs *RaftServer)request(writer http.ResponseWriter,request *http.Request) {
	request.ParseForm()
	//接收网页的消息
	if len(request.Form["data"])>0 {
		fmt.Println(request.Form["data"][0])
		rs.CostomerMsg<-request.Form["data"][0]
	}

}
