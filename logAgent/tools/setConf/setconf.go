package main

import (
	"go.etcd.io/etcd/clientv3"
	// 为了不因为版本变了改代码，通用好，可以考虑重新命名这个包
	"context"
	"fmt"
	"time"
	//"github.com/coreos/etcd/clientv3"换地方了，加入go的包了
	"LCS/logAgent/tailf"
	"encoding/json"
)

const (
	EtcdKey = "/oldboy/backend/logagent/config/192.168.76.139"
)

func SetLogConfToEtcd() {
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

	var logConfArr []tailf.CollectConf
	logConfArr = append(logConfArr,
		tailf.CollectConf{
			LogPath: "/home/alan8254402/go/src/LCS/logsForTest/logAgent.log",
			Topic:   "nginx_log",
		},
	)
	logConfArr = append(logConfArr,
		tailf.CollectConf{
			LogPath: "../main/logs/nginx/err.log",
			Topic:   "nginx_log_err2",
		},
	)
	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("marshal data failed err=", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	//cli.Delete(ctx, EtcdKey)
	//return
	//将序列化后的要收集的日志的配置信息存入etcd
	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed,err", err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//将刚才存入的信息取出来并打印查看是否正确？
	//注意：了解下get出来的resp，结构体的字段的含义
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed,err", err)
		return
	}
	for _, keyValue := range resp.Kvs {
		fmt.Printf("range resp.kvs：%s:%s\n", keyValue.Key, keyValue.Value)
	}

}
func main() {
	SetLogConfToEtcd()
}
