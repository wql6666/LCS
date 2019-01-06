package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379",
			"localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed,err", err)
		return
	}
	fmt.Println("conn succ")
	defer cli.Close()

	cli.Put(context.Background(), "/logAgent/conf/", "12345")
	//监控节点的变化
	for {
		rch := cli.Watch(context.Background(), "/logAgent/conf/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s,%q:%q\n", ev.Type, ev.Kv.Key,
					ev.Kv.Value)
			}
		}
	}
}
