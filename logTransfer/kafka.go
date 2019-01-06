package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

//var (
//	wg sync.WaitGroup
//
//)
//用来保存client，不然当client出现错误的时候都没有实例来查错误
//一般结构体用大写，实例小写？
type KafkaClient struct {
	consumer sarama.Consumer
	addr     string
	topic    string
	wg       sync.WaitGroup
}

var (
	kafkaClient *KafkaClient
)

func initKafka(addr string, topic string) (err error) {
	//创建一个kafkaClient结构体的实例
	kafkaClient = &KafkaClient{}
	//创建一个consumer;注意：这里用了一个strings包有什么作用，config怎么是nil
	//注意：newConsumer addr的参数需要传一个切片，这里用的strings包，为什么
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("init kafka failed ,err=%v", err)
		fmt.Println("failed to start consumer err:%s", err)
		return
	}
	//将信息保存到kafkaClient对象中
	kafkaClient.consumer = consumer
	kafkaClient.addr = addr
	kafkaClient.topic = topic
	return
	/*下边干活的代码，干活和初始化的分开
		//获得topic分区的数量
		partitionList, err := consumer.Partitions(topic)
		if err != nil {
			logs.Error("failed to get the list of "+
			"partitions err:", err)
			fmt.Println("failed to get the list of "+
				"partitions err:", err)
			return
		}
		fmt.Println(partitionList)

		for _, partition := range partitionList {
			pc, err := consumer.ConsumePartition("nginx_log", int32(partition), sarama.OffsetNewest)
			if err != nil {
				logs.Error("failed to start consumer for partition"+
					"%d:%s\n", partition, err)
				return
			}
			defer pc.AsyncClose()
			go func(sarama.PartitionConsumer) {
				//wg.Add(1)
				//起一个goroutine就加1，里边数据是几就加几，实际上里边是一个计数器
				for msg := range pc.Messages() {
					logs.Debug("partition:%d,offset:%d,Key:%s,Value:%s",
						msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
					fmt.Println()
				}
				//wg.Done() //完成一个就减一
			}(pc)
		}
		//time.Sleep(time.Hour)
		//wg.Wait() //等于0的时候就不阻塞
		consumer.Close()
	return
	*/
}
