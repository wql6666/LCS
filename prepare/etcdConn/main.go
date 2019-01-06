package main

import (
	//etcdClient "go.etcd.io/etcd/clientv3"
	// 为了不因为版本变了改代码，通用好，可以考虑重新命名这个包
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	//new生成一个客户端，连接etcd
	cli, err := clientv3.New(clientv3.Config{
		//Endpoints是etcd的ip端口，也可以写域名，是个数组，是个集群可能有多个
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second, //etcd连接超时
	})
	if err != nil {
		fmt.Println("connect failed ,err", err)
		return
	}
	fmt.Println("connect success")
	defer cli.Close()
}
