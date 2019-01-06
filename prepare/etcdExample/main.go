package main

import (
	"go.etcd.io/etcd/clientv3"
	// 为了不因为版本变了改代码，通用好，可以考虑重新命名这个包
	"context"
	"fmt"
	"time"
	//"github.com/coreos/etcd/clientv3"换地方了，加入go的包了
	"encoding/json"
)

const (
	EtcdKey = "/oldboy/backend/logagent/config/192.168.76.137"
)

type LogConf struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
	//sengQps int
}

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

	var logConfArr []LogConf
	logConfArr = append(logConfArr,
		LogConf{
			Path:  "../main/logs/nginx/access.log",
			Topic: "nginx_log",
		},
	)
	logConfArr = append(logConfArr,
		LogConf{
			Path:  "../main/logs/nginx/err.log",
			Topic: "nginx_log_err",
		},
	)
	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("put failed err=", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed,err", err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed,err", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}

}
func main() {
	//SetLogConfToEtcd()
	EtcdExample()

}

func EtcdExample() {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "/home/alan8254402/go/src/LCS/logAgent/conf/", "sampleValue")
	cancel()
	if err != nil {
		fmt.Println("put failed,err", err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/home/alan8254402/go/src/LCS/logAgent/conf/")
	cancel()
	if err != nil {
		fmt.Println("get failed,err", err)
		return
	}
	fmt.Println("resp的值", resp)
	for _, ev := range resp.Kvs {
		fmt.Printf("resp.kVs:%s:%s\n", ev.Key, ev.Value)
	}
}
