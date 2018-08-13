# etcd - 一个分布式一致性键值存储系统

etcd是一个分布式一致性键值存储系统，用于共享配置和服务发现，专注于：

* 简单:良好定义的，面向用户的API (gRPC)

* 安全： 带有可选客户端证书认证的自动TLS

* 快速:测试验证，每秒10000写入

* 可靠:使用Raft适当分布

etcd是Go编写，并使用Raft一致性算法来管理高可用复制日志，架构如下图所示：


## 下载安装

```
$ mkdir -p $GOPATH/src/github.com/coreos
$ cd !$
$ git clone https://github.com/coreos/etcd.git
$ cd etcd
$ ./build
$ ./bin/etcd
```

另外一种下载安装的方法：

直接下载etcd二进制 （包含etcd、etcdctl）
https://github.com/coreos/etcd/releases

## 测试

```
$ cd $GOPATH
$ ./bin/etcd

$ cd $GOPATH
$ ETCDCTL_API=3 ./bin/etcdctl put foo bar

# 输出结果显示OK，表示安装成功
OK
```

## 搭建本地集群

```
$ go get github.com/mattn/goreman

$ cd $GOPATH/src/github.com/coreos/etcd
$ goreman -f Procfile start
```


查看本地集群的服务器列表

```
$ cd $GOPATH/src/github.com/coreos/etcd

$ ./bin/etcdctl member list

# 显示结果

8211f1d0f64f3269: name=infra1 peerURLs=http://127.0.0.1:12380 clientURLs=http://127.0.0.1:2379 isLeader=false
91bc3c398fb3c146: name=infra2 peerURLs=http://127.0.0.1:22380 clientURLs=http://127.0.0.1:22379 isLeader=true
fd422379fda50e48: name=infra3 peerURLs=http://127.0.0.1:32380 clientURLs=http://127.0.0.1:32379 isLeader=false

```

### 存储数据

```
export ETCDCTL_API=3

$ ./bin/etcdctl put foo "Hello World!"

OK

$ ./bin/etcdctl get foo

foo
Hello World!


$ ./bin/etcdctl  --write-out="json" get foo

{"header":{"cluster_id":17237436991929493444,"member_id":9372538179322589801,"revision":2,"raft_term":2},"kvs":[{"key":"Zm9v","create_revision":2,"mod_revision":2,"version":1,"value":"SGVsbG8gV29ybGQh"}],"count":1}
```

### 根据前缀查询

```

$ ./bin/etcdctl put web1 value1
$ ./bin/etcdctl put web2 value2
$ ./bin/etcdctl put web3 value3

$ ./bin/etcdctl get web --prefix

web1
value1
web2
value2
web3
value3
```

### 删除数据

```
$ ./bin/etcdctl put key myvalue
$ ./bin/etcdctl del key
1
$ ./bin/etcdctl get key
// 查询结果为空

$ ./bin/etcdctl put k1 value1
$ ./bin/etcdctl put k2 value2
$ ./bin/etcdctl del k --prefix
2
$ ./bin/etcdctl get k --prefix
// 查询结果为空
```

### 事务写入

```
$ ./bin/etcdctl put user1 bad
OK
$ ./bin/etcdctl txn --interactive

compares:
// 输入以下内容，输入结束按 两次回车
value("user1") = "bad"      

//如果 user1 = bad，则执行 get user1 
success requests (get, put, del):
get user1

//如果 user1 != bad，则执行 put user1 good
failure requests (get, put, del):
put user1 good

// 运行结果，执行 success
SUCCESS

user1
bad



$ ./bin/etcdctl txn --interactive
compares:
value("user1") = "111"  

// 如果 user1 = 111，则执行 get user1 
success requests (get, put, del):
get user1

//如果 user1 != 111，则执行 put user1 2222
failure requests (get, put, del):
put user1 2222

// 运行结果，执行 failure
FAILURE

OK

$ ./bin/etcdctl get user1
user1
2222
```

### watch 

```
// 当 stock1 的数值改变（ put 方法）的时候，watch 会收到通知
$ ./bin/etcdctl watch stock1

// 新打开终端
$ export ETCDCTL_API=3
$ ./bin/etcdctl put stock1 1000

//在watch 终端显示
PUT
stock1
1000


$ ./bin/etcdctl watch stock --prefix
$ ./bin/etcdctl put stock1 10
$ ./bin/etcdctl put stock2 20
```

### lease

```
$ ./bin/etcdctl lease grant 300
# lease 326963a02758b527 granted with TTL(300s)

$ ./bin/etcdctl put sample value --lease=326963a02758b527
OK

$ ./bin/etcdctl get sample

$ ./bin/etcdctl lease keep-alive 326963a02758b520
$ ./bin/etcdctl lease revoke 326963a02758b527
lease 326963a02758b527 revoked

# or after 300 seconds
$ ./bin/etcdctl get sample
```

### Distributed locks

```
//第一终端
$ ./bin/etcdctl lock mutex1
mutex1/326963a02758b52d

# 第二终端
$ ./bin/etcdctl lock mutex1

// 当第一个终端结束了，第二个终端会显示
mutex1/326963a02758b531

```

### Elections

```
$ ./bin/etcdctl elect one p1

one/326963a02758b539
p1


# another client with the same name blocks
$ ./bin/etcdctl elect one p2

//结束第一终端，第二终端显示
one/326963a02758b53e
p2
```

### Cluster status

集群状态

```
$ ./bin/etcdctl --write-out=table endpoint status

$ ./bin/etcdctl endpoint health
```

### Snapshot


```
./bin/etcdctl snapshot save my.db

Snapshot saved at my.db

./bin/etcdctl --write-out=table snapshot status my.db
```

### Member

```
./bin/etcdctl member list -w table
```




