package main

import (
	"time"
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	//"log"
)

//用代码实现单机版的增删查

var (
	//连接超时时间
	dialTimeout =5 *time.Second
	//请求超时时间
	requestTimeout = 2*time.Second
	//设置访问的IP
	endpoints = []string{"127.0.0.1:2379"}

)


//插入数据
func putValue(cli *clientv3.Client,key ,value string) string {
	fmt.Println("向服务中存储数据")
	_,err:=cli.Put(context.TODO(),key,value)
	if err!=nil {
		return "FAIL"
	}else {
		return  "OK"
	}
}


//获取数据
func getValue(cli *clientv3.Client,key string )[] *mvccpb.KeyValue{
	fmt.Println("获取值")
	resp,err:=cli.Get(context.TODO(),key)
	if err!=nil {
		fmt.Println(err)
	} else {
		//返回对应key 的结果
		return resp.Kvs
	}
	return nil
}

func delValue(cli *clientv3.Client,key string ) int64 {
	fmt.Println("按键删除键值对")
	delRes,err:=cli.Delete(context.TODO(),key)
	if err!=nil {
		fmt.Println(err)
	}else {
		return delRes.Deleted
	}
	return 0
}

//按照前缀获取数据
func getValueByPrefix(cli *clientv3.Client,key string ) []*mvccpb.KeyValue{
	fmt.Println("按照key的前缀获取数据")
	resp,err:=cli.Get(context.TODO(),key,clientv3.WithPrefix())
	if err!=nil {
		fmt.Println(err)
	}else {
		return resp.Kvs
	}
	return  nil

}


//按前缀删除某写数据
func delValueByPrefix(cli *clientv3.Client,key string)int64 {
	fmt.Println("按前缀删除数据")
	resp,err:=cli.Delete(context.TODO(),key,clientv3.WithPrefix())
	if err!=nil {
		fmt.Println(err)
	}else {
		return  resp.Deleted
	}
	return  0
}

func main() {

	//创建访问etcd的对象
	cli ,err:=clientv3.New(clientv3.Config{
		Endpoints:endpoints,
		DialTimeout:dialTimeout,
	})
	if err!=nil {
		fmt.Println(err)
	}
	//延时关闭
	defer  cli.Close()


	//完成加入数据的功能
	//putResult:=putValue(cli,"a","abc")
	//fmt.Println(putResult)


	//删除数据
	//delResult:=delValue(cli,"a")
	//fmt.Println(delResult)

	//查询数据
	getResult:=getValue(cli,"a")
	if len(getResult)>=1 {
		for _,item:=range getResult {
			fmt.Printf("value为: %s",string(item.Value))
		}
	}


	putValue(cli,"abc","1000")
	putValue(cli,"abb","2000")
	putValue(cli,"acc","3000")

	//安前缀删除
	relt := delValueByPrefix(cli,"a")
	fmt.Println(relt)

	getResult:=getValueByPrefix(cli,"a")
	if len(getResult)>=1{
		for _,item:=range getResult {
			fmt.Println(string(item.Value))
		}
	}


	//事务处理
	putValue(cli,"user1","bad1")
	_,err1:=cli.Txn(context.TODO()).
	If(clientv3.Compare(clientv3.Value("user1"),"=","bad1")).
		Then(clientv3.OpPut("user1","good")).
			Else(clientv3.OpPut("user1","bad-bad")).Commit()

			if err1!=nil {
				fmt.Println(err)
			}
	rlt:=getValue(cli,"user1")
	log.Printf("获取值的key: %s ，value：%s", string(rlt[0].Key),string(rlt[0].Value))





	//设置租期时间5s
	resp,err1:=cli.Grant(context.TODO(),5)
	_,err1=cli.Put(context.TODO(),"foo","bar",clientv3.WithLease(resp.ID))

	if err1 != nil {
		fmt.Println(err1)
	}

	getResp,err2:=cli.Get(context.TODO(),"foo")
	if err2!=nil {
		fmt.Println(err2)
	}
	fmt.Println("这是六秒前：",string(getResp.Kvs[0].Key),string(getResp.Kvs[0].Value))


	time.Sleep(6*time.Second)


	getResp,err2=cli.Get(context.TODO(),"foo")
	if err2!=nil {
		fmt.Println(err2)
	}
	if len(getResp.Kvs)>=1 {
		fmt.Println("这是六秒后:", getResp.Kvs[0].Key, getResp.Kvs[0].Value)
	}





	putValue(cli,"foo","abc")
	go func() {
		rch:=cli.Watch(context.Background(),"foo")
		for wresp:=range rch {
			for _,ev:=range wresp.Events {
				fmt.Printf("Watch foo----%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()
	putValue(cli,"foo","bcd")


	/*
	go func() {
		rch:=cli.Watch(context.Background(),"fooo",clientv3.WithPrefix())
		for wresp:=range rch {
			for _,ev:=range wresp.Events {
				fmt.Printf("Watch foo----%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	putValue(cli,"fooo1","abc")

	putValue(cli,"fooo2","bcd")

	putValue(cli,"fooo3","def")
	*/


	//监听某个范围的键
	go func() {
		rch:=cli.Watch(context.Background(),"fo0",clientv3.WithRange("fo3"))
		for wresp:=range rch {
			for _,ev:=range wresp.Events {
				fmt.Printf("Watch foo----%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	putValue(cli,"fo0","a")
	putValue(cli,"fo1","b")
	putValue(cli,"fo2","c")
	//fo3监听不到
	putValue(cli,"fo3","d")

}
