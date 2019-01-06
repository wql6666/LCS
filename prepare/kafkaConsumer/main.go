package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	//创建一个consumer
	consumer, err := sarama.NewConsumer(strings.Split("192.168.76.138:9092", ","), nil)
	if err != nil {
		fmt.Println("failed to start consumer err:%s", err)
		return
	}
	//获得topic分区的数量
	partitionList, err := consumer.Partitions("nginx_log")
	if err != nil {
		fmt.Println("failed to get the list of "+
			"partitions err:", err)
		return
	}
	fmt.Println(partitionList)

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition("nginx_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition"+
				"%d:%s\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			wg.Add(1)
			//起一个goroutine就加1，里边数据是几就加几，实际上里边是一个计数器
			for msg := range pc.Messages() {
				fmt.Printf("partition:%d,offset:%d,Key:%s,Value:%s",
					msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				fmt.Println()
			}
			wg.Done() //完成一个就减一
		}(pc)
	}
	//time.Sleep(time.Hour)
	wg.Wait() //等于0的时候就不阻塞
	consumer.Close()

}
