# 003 etcd 本机手动测试集群及 member 操作

## 集群搭建

下面只用同一台服务器进行三个成员节点的开启

**节点1**

```
./etcd --name cd0 --initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://127.0.0.1:2380 --listen-client-urls http://127.0.0.1:2379 --advertise-client-urls http://127.0.0.1:2379 --initial-cluster-token etcd-cluster-1 --initial-cluster cd0=http://127.0.0.1:2380,cd1=http://127.0.0.1:2480,cd2=http://127.0.0.1:2580 --initial-cluster-state new
```

**节点2**

```
./etcd --name cd1 --initial-advertise-peer-urls http://127.0.0.1:2480 --listen-peer-urls http://127.0.0.1:2480 --listen-client-urls http://127.0.0.1:2479 --advertise-client-urls http://127.0.0.1:2479 --initial-cluster-token etcd-cluster-1 --initial-cluster cd0=http://127.0.0.1:2380,cd1=http://127.0.0.1:2480,cd2=http://127.0.0.1:2580 --initial-cluster-state new
```

**节点3**

```
./etcd --name cd2 --initial-advertise-peer-urls http://127.0.0.1:2580 --listen-peer-urls http://127.0.0.1:2580 --listen-client-urls http://127.0.0.1:2579 --advertise-client-urls http://127.0.0.1:2579 --initial-cluster-token etcd-cluster-1 --initial-cluster cd0=http://127.0.0.1:2380,cd1=http://127.0.0.1:2480,cd2=http://127.0.0.1:2580 --initial-cluster-state new
```

## 查询 member 列表

```
export ETCDCTL_API=3
ENDPOINTS=127.0.0.1:2379,127.0.0.1:2479,127.0.0.1:2579

./etcdctl --endpoints=$ENDPOINTS member list
```

运行结果：

```
98f0c6bf64240842, started, cd2, http://127.0.0.1:2580, http://127.0.0.1:2579
bf9071f4639c75cc, started, cd0, http://127.0.0.1:2380, http://127.0.0.1:2379
e3ba87c3b4858ef1, started, cd1, http://127.0.0.1:2480, http://127.0.0.1:2479
```

## 添加 member 节点

**member add** 添加节点

```
./etcdctl --endpoints=$ENDPOINTS member add cd3 --peer-urls=http://127.0.0.1:2180
Member b9057cfdc8ff17ce added to cluster 9da8cd75487bd6dc
```

运行结果：

```
ETCD_NAME="cd3"
ETCD_INITIAL_CLUSTER="cd2=http://127.0.0.1:2580,cd3=http://127.0.0.1:2180,cd0=http://127.0.0.1:2380,cd1=http://127.0.0.1:2480"
ETCD_INITIAL_ADVERTISE_PEER_URLS="http://127.0.0.1:2180"
ETCD_INITIAL_CLUSTER_STATE="existing"
```

**查询 member 节点列表信息**

```
./etcdctl --endpoints=$ENDPOINTS member list
```

运行结果：

```
98f0c6bf64240842, started, cd2, http://127.0.0.1:2580, http://127.0.0.1:2579
b9057cfdc8ff17ce, unstarted, , http://127.0.0.1:2180, 
bf9071f4639c75cc, started, cd0, http://127.0.0.1:2380, http://127.0.0.1:2379
e3ba87c3b4858ef1, started, cd1, http://127.0.0.1:2480, http://127.0.0.1:2479
```

通过查询结果可以发现：`http://127.0.0.1:2180` 显示状态为：`unstarted` 

**启动新节点**

```
./etcd --name cd3 --listen-client-urls http://127.0.0.1:2179 --advertise-client-urls http://127.0.0.1:2179 --listen-peer-urls http://127.0.0.1:2180 --initial-advertise-peer-urls http://127.0.0.1:2180 --initial-cluster-state existing --initial-cluster cd2=http://127.0.0.1:2580,cd0=http://127.0.0.1:2380,cd3=http://127.0.0.1:2180,cd1=http://127.0.0.1:2480 --initial-cluster-token etcd-cluster-1
```

**查询 member 节点列表信息**

```
./etcdctl --endpoints=$ENDPOINTS member list
```

运行结果：

```
98f0c6bf64240842, started, cd2, http://127.0.0.1:2580, http://127.0.0.1:2579
b9057cfdc8ff17ce, started, cd3, http://127.0.0.1:2180, http://127.0.0.1:2179
bf9071f4639c75cc, started, cd0, http://127.0.0.1:2380, http://127.0.0.1:2379
e3ba87c3b4858ef1, started, cd1, http://127.0.0.1:2480, http://127.0.0.1:2479
```

## 删除 member

```
./etcdctl --endpoints=$ENDPOINTS member remove b9057cfdc8ff17ce
```

运行结果：

```
Member b9057cfdc8ff17ce removed from cluster 9da8cd75487bd6dc
```

**查询 member 节点列表信息**

```
./etcdctl --endpoints=$ENDPOINTS member list
```

运行结果：

```
98f0c6bf64240842, started, cd2, http://127.0.0.1:2580, http://127.0.0.1:2579
bf9071f4639c75cc, started, cd0, http://127.0.0.1:2380, http://127.0.0.1:2379
e3ba87c3b4858ef1, started, cd1, http://127.0.0.1:2480, http://127.0.0.1:2479
```

## 代码实现 member 的管理

### 添加节点

```
func addMember(cli *clientv3.Client)  {
    peerURLs := []string{"http://127.0.0.1:2180"}

    mresp, err := cli.MemberAdd(context.Background(), peerURLs)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("added member.PeerURLs:", mresp.Member.PeerURLs)
    resp, err := cli.MemberList(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("添加后 members:", resp.Members)
}
```

执行完 添加节点，需要打开终端启动节点服务器

```
./etcd --name cd3 --listen-client-urls http://127.0.0.1:2179 --advertise-client-urls http://127.0.0.1:2179 --listen-peer-urls http://127.0.0.1:2180 --initial-advertise-peer-urls http://127.0.0.1:2180 --initial-cluster-state existing --initial-cluster cd2=http://127.0.0.1:2580,cd0=http://127.0.0.1:2380,cd3=http://127.0.0.1:2180,cd1=http://127.0.0.1:2480 --initial-cluster-token etcd-cluster-1
```

**如果启动失败，需要删除 cd3 的信息**

![](http://olgjbx93m.bkt.clouddn.com/20180129-104838.png)

### 删除节点

```
// 删除节点
func delMember (cli *clientv3.Client, memberId uint64) {

	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.MemberRemove(context.Background(), memberId)
	if err != nil {
		log.Fatal(err)
	}

	resp, err = cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("删除后 members:", resp.Members)
}
```

### 全部代码

```
package main

import (
	"github.com/coreos/etcd/clientv3"
	"log"
	"fmt"
	"context"
	"time"
)
var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
	endpoints      = []string{"127.0.0.1:2379"}
)
func main()  {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("members:", resp.Members)


	//添加member
	//addMember(cli)


	// 删除节点
	//delMember(cli,uint64(7438291228984697304))

}

func addMember(cli *clientv3.Client)  {
	peerURLs := []string{"http://127.0.0.1:2180"}

	mresp, err := cli.MemberAdd(context.Background(), peerURLs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added member.PeerURLs:", mresp.Member.PeerURLs)
	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("添加后 members:", resp.Members)
}

// 删除节点
func delMember (cli *clientv3.Client, memberId uint64) {

	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.MemberRemove(context.Background(), memberId)
	if err != nil {
		log.Fatal(err)
	}

	resp, err = cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("删除后 members:", resp.Members)
}
```


