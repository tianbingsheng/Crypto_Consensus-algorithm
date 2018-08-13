package main

import (
	"time"
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"context"
)

//etcd网络中角色和用户的使用
//角色中可以设置功能（只读、读写），然后将用户添加到角色中，那么用户则具备了角色中设置的功能


//设置连接服务器的参数
var (
	dialTimeout =5*time.Second
	requestTimeout = 2 *time.Second
	endPoints = []string{"127.0.0.1:2379"}
)


func main() {
	//链接服务器
	//cli ,err:=clientv3.New(clientv3.Config{
	//	//设置参数
	//	Endpoints:endPoints,
	//	DialTimeout:dialTimeout,
	//})
	//if err!=nil {
	//	fmt.Println(err)
	//}
	//defer cli.Close()
	//
	//
	////添加角色，必须首先root角色
	//_,err=cli.RoleAdd(context.TODO(),"root")
	//if err!=nil {
	//	fmt.Println(err)
	//}
	//
	//
	////给刚才的root角色添加读写权限
	//_,err=cli.RoleGrantPermission(context.TODO(),
	//	"root",
	//		"foo",
	//			"foo",
	//				clientv3.PermissionType(clientv3.PermReadWrite),
	//		)
	//if err!=nil {
	//	fmt.Println(err)
	//}



	//测试，打印root角色中的功能
	//roleInfo,err:=cli.RoleGet(context.TODO(),"root")
	//if err!=nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(roleInfo.Perm)
	//}


	//创建用户，主要，第一个用户名必须为root
	//以后在添加用户，用户名就可以随便写了
	//_,err = cli.UserAdd(context.TODO(),"root","123")
	//if err!=nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println("创建用户成功")
	//}


	//添加新用户abc
	//_,err=cli.UserAdd(context.TODO(),"abc","123")
	//if err!=nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("创建用户成功")
	//}
	//
	//_,err = cli.UserGrantRole(context.TODO(),"abc","root")
	//if err!=nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("abc 授权成功")
	//}



	//设置用户的角色
	//_,err= cli.UserGrantRole(context.TODO(),"root","root")
	//if err!=nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println("给root用户授予root角色成功")
	//}
	//同过以上这段代码的授权，则root用户对foo键值对具备了readwrite的权限



	//获取用户信息
	//userInfo,err:=cli.UserGet(context.TODO(),"root")
	//if err!=nil {
	//	fmt.Println(err)
	//} else {
	//	//打印root用户的角色
	//	fmt.Println(userInfo.Roles)
	//}







	//用root用户，去链接服务器
	cliAuth,err:=clientv3.New(clientv3.Config{
		Endpoints:endPoints,
		DialTimeout:dialTimeout,
		Username:"abc",
		Password:"123",
	})

	if err!=nil {
		fmt.Println(err)
	}

	defer cliAuth.Close()




	//向服务器中保存数据
	_,err = cliAuth.Put(context.TODO(),"foo","helloff")
	if err!=nil {
		fmt.Println("保存失败",err)
	}


	//从服务器中读取数据
	resp,err:=cliAuth.Get(context.TODO(),"foo")
	if err !=nil {
		fmt.Print(err)
	} else {
		//读
		if len(resp.Kvs)>=1 {
			fmt.Println(string(resp.Kvs[0].Key))
			fmt.Println(string(resp.Kvs[0].Value))
		}
	}




}