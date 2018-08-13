# 002 etcd 多台服务器集群

## 下载安装 etcd

```
$ mkdir -p $GOPATH/src/github.com/coreos
$ cd !$
$ git clone https://github.com/coreos/etcd.git
$ cd etcd
$ ./build
```

## 启动服务

每个节点都要执行以下配置，HOST_1、HOST_2、HOST_3 分别设置为多台服务器的IP

```
TOKEN=token-03
CLUSTER_STATE=new
NAME_1=machine-1
NAME_2=machine-2
NAME_3=machine-3
HOST_1=192.168.1.105
HOST_2=192.168.1.143
HOST_3=192.168.1.103
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
```

**machine 1** 执行如下命令

```
$ cd $GOPATH/src/github.com/coreos/etcd/bin

# For machine 1
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
./etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

```

**machine 2** 执行如下命令

```
$ cd $GOPATH/src/github.com/coreos/etcd/bin

# For machine 2
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
./etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

```

**machine 3** 执行如下命令

```
$ cd $GOPATH/src/github.com/coreos/etcd/bin

# For machine 3
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
./etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

检测服务器运行是否正常

```
$ cd $GOPATH/src/github.com/coreos/etcd/bin

export ETCDCTL_API=3
HOST_1=192.168.1.105
HOST_2=192.168.1.143
HOST_3=192.168.1.103
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379

./etcdctl --endpoints=$ENDPOINTS member list
```

## 存储数据

```
./etcdctl --endpoints=$ENDPOINTS put foo "Hello World!"

./etcdctl --endpoints=$ENDPOINTS get foo
./etcdctl --endpoints=$ENDPOINTS --write-out="json" get foo
```

## 根据前缀查询

```
./etcdctl --endpoints=$ENDPOINTS put web1 value1
./etcdctl --endpoints=$ENDPOINTS put web2 value2
./etcdctl --endpoints=$ENDPOINTS put web3 value3

./etcdctl --endpoints=$ENDPOINTS get web --prefix

web1
value1
web2
value2
web3
value3
```

## 删除

```
./etcdctl --endpoints=$ENDPOINTS put key myvalue
./etcdctl --endpoints=$ENDPOINTS del key

./etcdctl --endpoints=$ENDPOINTS put k1 value1
./etcdctl --endpoints=$ENDPOINTS put k2 value2
./etcdctl --endpoints=$ENDPOINTS del k --prefix
```

## 事务写入

```
$ ./etcdctl --endpoints=$ENDPOINTS put user1 bad
OK

$ ./etcdctl --endpoints=$ENDPOINTS txn --interactive

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


$ ./etcdctl --endpoints=$ENDPOINTS txn --interactive
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

$ ./etcdctl --endpoints=$ENDPOINTS get user1
user1
2222
```

## watch 

```
// 当 stock1 的数值改变（ put 方法）的时候，watch 会收到通知
./etcdctl --endpoints=$ENDPOINTS watch stock1

// 新打开终端
$ cd $GOPATH/src/github.com/coreos/etcd/bin

export ETCDCTL_API=3
HOST_1=192.168.1.126
HOST_2=192.168.1.119
HOST_3=192.168.1.103
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379
./etcdctl --endpoints=$ENDPOINTS put stock1 1000

./etcdctl --endpoints=$ENDPOINTS watch stock --prefix

./etcdctl --endpoints=$ENDPOINTS put stock1 10
./etcdctl --endpoints=$ENDPOINTS put stock2 20
```

## 更多操作

https://coreos.com/etcd/docs/latest/demo.html

