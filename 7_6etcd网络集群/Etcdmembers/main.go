package main

import (
	"time"
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"context"
)

//通过代码如何添加和删除新节点

var (
	dialTimeout = 5 *time.Second
	requestTimeout = 2 *time.Second
	endpints = []string{"127.0.0.1:2379"}
)

func main() {

	cli,err:=clientv3.New(clientv3.Config{
		Endpoints:endpints,
		DialTimeout:dialTimeout,
	})
	if err!=nil {
		fmt.Println(err)
	}
	defer cli.Close()

	//查看系统中的menber
	//background可理解为后台进行，todo可理解为前台进行
	//前台进行在主线程中执行，后台进行在子线程中运行
	resp,err:=cli.MemberList(context.Background())
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("members:",resp.Members[0].Name)


	//添加新节点
	//addMember(cli)

	delMember(cli,0xad2af19e12799db9)

}

//向集群中添加节点
func addMember(cli *clientv3.Client) {
	peerURLs:=[]string{"http://127.0.0.1:2180"}

	mresp,err:=cli.MemberAdd(context.Background(),peerURLs)
	if err!=nil {
		fmt.Println(err)
	}

	fmt.Println("added member PeerURLs",mresp.Member.PeerURLs)

	resp,err:=cli.MemberList(context.Background())
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("添加后的members:",resp.Members)


}



func delMember(cli *clientv3.Client,memberId uint64) {
	_,err:=cli.MemberRemove(context.Background(),memberId)
	if err!=nil {
		fmt.Println(err)
	}

}



