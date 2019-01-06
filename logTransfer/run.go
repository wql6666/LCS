package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

func run() (err error) {
	//获得topic分区的数量
	partitionList, err := kafkaClient.consumer.Partitions(kafkaClient.topic)
	if err != nil {
		logs.Error("failed to get the list of "+
			"partitions err:", err)
		fmt.Println("failed to get the list of "+
			"partitions err:", err)
		return
	}
	fmt.Println(partitionList)

	for partition := range partitionList {
		//得到每一个分区的消费者,OffsetNewest,从常量-1，表示从头开始消费，
		// OffsetOldest int64 = -2表示从上次的地方开始消费。这里确认下。
		partitionConsumer, errRet := kafkaClient.consumer.ConsumePartition(
			kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			logs.Error("failed to start consumer for partition"+
				"%d:%s\n", partition, err)
			return
		}
		//注意：这里的关闭啥意思？
		defer partitionConsumer.AsyncClose()
		//每个分区开启一个goroutine消费，匿名函数,
		// kafka是一个消息队列，之前logAgent有存入消息到kafka，现在就可以消费消息
		go func(partitionConsumer sarama.PartitionConsumer) {
			//起一个goroutine就加1，里边数据是几就加几，实际上里边是一个计数器
			kafkaClient.wg.Add(1)
			//消费信息
			for msg := range partitionConsumer.Messages() {
				logs.Debug("partition:%d,offset:%d,Key:%s,Value:%s",
					msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//fmt.Println()
				//信息发送到ES中去,两个信息，一个topic，一个消息的value值
				err = sendToES(kafkaClient.topic, msg.Value)
				if err != nil {
					logs.Warn("send to ES failed,err=", err)
				}
			}
			kafkaClient.wg.Done() //完成一个就减一
		}(partitionConsumer)
	}
	//time.Sleep(time.Hour)
	kafkaClient.wg.Wait() //等于0的时候就不阻塞
	logs.Debug("run succ")
	return

}
