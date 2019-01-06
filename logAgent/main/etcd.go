package main

import (
	"LCS/logAgent/tailf"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type EtcdClient struct {
	client *clientv3.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
)

func initEtcd(addr string, key string) (CollectConf []tailf.CollectConf, err error) {
	//new生成一个连接etcd的客户端
	cli, err := clientv3.New(clientv3.Config{
		//Endpoints是etcd的ip端口，也可以写域名，是个数组，是个集群可能有多个
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second, //etcd连接超时控制
	})
	if err != nil {
		logs.Error("connect etcd failed ,err", err)
		return
	}
	logs.Debug("connect etcd succ")

	etcdClient = &EtcdClient{
		client: cli,
	}
	//从etcd中获取配置，优化的话，就把上边和下边的拆成2个函数
	//先获取etcd中的key，这个key是之前setconf函数存进去的

	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}

	for _, localIp := range localIpArray {
		etcdKey := fmt.Sprintf("%s%s", key, localIp)
		//将取出来的key存到keys的切片中
		etcdClient.keys = append(etcdClient.keys, etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("client get from etcd failed,err%v", err)
			continue
		}
		cancel()
		logs.Debug("resp from etcd:%v", resp.Kvs)
		for _, keyValue := range resp.Kvs {
			if string(keyValue.Key) == etcdKey {
				// 将反序列话后的要搜集的日志的配置信息存入&CollectConf
				err = json.Unmarshal(keyValue.Value, &CollectConf)
				if err != nil {
					logs.Error("unmarshal failed,err", err)
					continue
				}
				logs.Debug("log config is %v", CollectConf)

			}
		}
	}
	initEtcdWatcher()
	//将etcd中的信息return回去
	return

}

func initEtcdWatcher() {
	for _, key := range etcdClient.keys {
		//需要时刻监控？那么就开一个goroutine一直监控就好了
		go watchKey(key)
	}
}

func watchKey(key string) {
	//new生成一个客户端，连接etcd，
	// 每次需要用到etcd的时候都需要根据ip或域名，生成一个客户端连接etcd
	cli, err := clientv3.New(clientv3.Config{
		//Endpoints是etcd的ip端口，也可以写域名，是个数组，是个集群可能有多个
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second, //etcd连接超时
	})
	if err != nil {
		logs.Error("connect etcd failed ,err", err)
		return
	}
	logs.Debug("begin watch key:%s", key)
	fmt.Printf("begin watch key:%s", key)
	//死循环watch key
	for {
		//watch的结果是一个通道
		watchChan := cli.Watch(context.Background(), key)
		//把变更的配置保存起来
		//创建一个collectConf的实例
		var collectConf []tailf.CollectConf
		var getConfSucc = true //设定一个标识符
		for watchResponce := range watchChan {
			for _, event := range watchResponce.Events {
				fmt.Printf("%s,%q:%q\n", event.Type, event.Kv.Key, event.Kv.Value)
				//根据watch的结果，动态获取etcd中存入，和删除的信息的动作
				if event.Type == mvccpb.DELETE {
					logs.Warn("dynamic key[%s]'s config deleted", key)
					continue

				}

				if event.Type == mvccpb.PUT && string(event.Kv.Key) == key {
					//将增加的配置信息进行存储到&collectConf
					err = json.Unmarshal(event.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key[%s],unmarshal ,err:%v", err)
						//有错误就不将错误的配置信息传到tailf中进行更新
						getConfSucc = false
						continue
					}

				}
				logs.Debug("dynamic get config from etcd,%s,%q:%q\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
			//时刻保持更新taif中的配置信息
			if getConfSucc {
				logs.Debug("get config from etcd succ,%v", collectConf)
				//更新配置，有新的配置就开始读取日志,存到msgChan里边
				tailf.UpdateConfig(collectConf)
			}

		}

	}
}
