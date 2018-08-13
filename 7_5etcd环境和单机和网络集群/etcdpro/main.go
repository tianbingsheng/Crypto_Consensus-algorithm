package main

import (
	"time"
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"context"
	//"github.com/coreos/etcd/mvcc"
	"github.com/coreos/etcd/mvcc/mvccpb"

)

//通过代码测试单机版etcd框架的使用
//添加，删除，查找，前缀，延时......
var (
	//设置请求链接时间
	dialTimeout = 5*time.Second
	//请求超时时间
	requestTimeout = 2 *time.Second
	//设置IP
	endPoints = []string{"127.0.0.1:2379"}
)

//添加键值对
func putValue(cli *clientv3.Client, key ,value string ) string{
	//向服务器添加键值对
	_,err:=cli.Put(context.TODO(),key ,value)
	if err!=nil {
		return "failed"
	}else {
		return "ok"
	}
}

//查询键值对
func getValue(cli *clientv3.Client,key string )[]*mvccpb.KeyValue{
	resp,err:=cli.Get(context.TODO(),key)

	if err !=nil {
		fmt.Println(err)
		return  nil
	} else {
		return resp.Kvs
	}
}

//返回的结果为删除了几条数据
func delValue(cli *clientv3.Client,key string ) int64{
	delRes,err:=cli.Delete(context.TODO(),key)
	if err!=nil {
		fmt.Println(err)
		return 0
	}else {
		//返回删除的个数
		return delRes.Deleted
	}

}


//按照前缀删除
func delValueByPrefix(cli *clientv3.Client,prefixKey string) int64{
	delRes,err:=cli.Delete(context.TODO(),prefixKey,clientv3.WithPrefix())
	if err!=nil {
		fmt.Println(err)
		return  0
	}else {
		return delRes.Deleted
	}
}


//按键的前缀查询系统中的键值对
func getValueByPrefix(cli *clientv3.Client,keyPrefix string)[]*mvccpb.KeyValue{

	resp,err:=cli.Get(context.TODO(),keyPrefix,clientv3.WithPrefix())
	if err!=nil {
		fmt.Println(err)
		return nil
	}else {

		return resp.Kvs
	}
}


func main() {

	//创建连接etcd的对象
	cli,err:=clientv3.New(clientv3.Config{
		Endpoints:endPoints,
		DialTimeout:dialTimeout,
	})
	if err!=nil {
		fmt.Println(err)
	}


	/*


	//插入数据
	putResult:=putValue(cli,"at","hello kongyi")
	fmt.Println(putResult)

	//getResult为集合
	getResult:=getValue(cli,"a")
	//打印键值对
	for _,item:=range getResult {
		fmt.Println(string(item.Key),string(item.Value))
	}


	//删除键值对
	cnt:=delValue(cli,"at")
	fmt.Println("删除了",cnt)


	//添加两个带a前缀的键值对
	putValue(cli,"aaa1","aaaa")
	putValue(cli,"aaa2","bbbb")

	//按照前缀删除
	//relt:=delValueByPrefix(cli,"a")
	//fmt.Println("按前缀删除的个数为",relt)


	//查询带前缀
	gets:=getValueByPrefix(cli,"a")
	for _,item:=range gets {
		//fmt.Sprintf("您输入key:%s　　值%s",item.Key,item.Value)
		fmt.Printf("您输入key:%s　　值%s\n",item.Key,item.Value)
	}


	//事务处理
	putValue(cli,"user1","bad")
	_,erro:=cli.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.Value("user1"),"=","bad")).
			Then(clientv3.OpPut("user1","good")).
				Else(clientv3.OpPut("user1","bad-good")).Commit()

				if erro!=nil{
					fmt.Println(err)
				}

				rlt:=getValue(cli,"user1")
				fmt.Println(rlt)
	            fmt.Printf("获取值的key: %s ，value：%s", string(rlt[0].Key),string(rlt[0].Value))






	//lease 设置有效期时间
	resp,err:=cli.Grant(context.TODO(),5)
	_,err=cli.Put(context.TODO(),"foo","bar",clientv3.WithLease(resp.ID))
	if err!=nil {
		fmt.Println(err)
	}
	rsp,e:=cli.Get(context.TODO(),"foo")
	if e!=nil {
		fmt.Println(e)
	}
	fmt.Println(rsp.Kvs[0].Key)
	fmt.Println(rsp.Kvs[0].Value)
	time.Sleep(6*time.Second)

	//6S后重新获取
	rsp,e=cli.Get(context.TODO(),"foo")
	if e!=nil {
		fmt.Println(e)
	}
	if len(rsp.Kvs)>1 {
		fmt.Println(rsp.Kvs[0].Key)
		fmt.Println(rsp.Kvs[0].Value)
	}



	*/





	/*
	//测试watch监听的使用
	putValue(cli,"a1","hello")
	//在子线程中监听
	go func() {
		rch:=cli.Watch(context.Background(),"a1")
		for wresp:=range rch {
			for _,ev:=range wresp.Events {
				fmt.Printf("watch...a1 %s %q %q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
			}
		}
	}()
	putValue(cli,"a1","abc")

	*/



	//监听某个key在一定范围内的value的变化
	//putValue(cli,"fo0","a")

	var th = make(chan bool)
	go func() {
		//左边包含右侧不包含
		rch:=cli.Watch(context.Background(),"foa",clientv3.WithRange("fod"))
		for wresp:=range rch {
			for _,ev:=range wresp.Events {

				fmt.Println(ev.Type,string(ev.Kv.Key),string(ev.Kv.Value))
			}
			<-th
		}

	}()

	putValue(cli,"foa","aa")
	th<-true

	putValue(cli,"fob","aa")
	th<-true

	putValue(cli,"foc","bb")
	th<-true

	//这行代码监听不到【０，３）
	putValue(cli,"fod","cc")
	th<-true


}
