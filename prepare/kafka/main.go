package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama" //kafka技术库
)

func main() {
	//初始化kafka的配置
	config := sarama.NewConfig()                              //实例化一个配置
	config.Producer.RequiredAcks = sarama.WaitForAll          //确保日志发到kafka
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机分区
	config.Producer.Return.Successes = true                   //表示写成功了

	//实例化一个同步producer客户端
	client, err := sarama.NewSyncProducer([]string{"192.168.76.137:9092"}, config)
	//指定kafka的ip端口，和配置信息，配置信息是从上边传入
	//未指定端口号也会报错：client has run out of available brokers to talk to
	// (Is your cluster reachable?)
	if err != nil {
		fmt.Println("producer colse err", err)
		return
	}

	defer client.Close()
	for {

		msg := &sarama.ProducerMessage{} //创建一个消息实例
		msg.Topic = "test2"              //队列名字
		//会编码消息
		msg.Value = sarama.StringEncoder("this is a good test,一个分区？")
		//调用client发送消息，返回一个分区id，和偏移量offset
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			fmt.Println("send message failed,err", err)
			return
		}

		fmt.Printf("pid:%v offset:%v\n", pid, offset)
		time.Sleep(10 * time.Millisecond)
	}
}
