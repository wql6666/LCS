package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	//接口不用指针，结构体用指针
	client sarama.SyncProducer
)

func InitKafka(addr string) (err error) {
	//初始化kafka的配置
	config := sarama.NewConfig()                              //实例化一个kafka配置
	config.Producer.RequiredAcks = sarama.WaitForAll          //确保日志发到kafka
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机分区,让负载均衡
	config.Producer.Return.Successes = true                   //表示写成功了

	//实例化一个同步producer客户端
	//指定kafka的ip端口，和配置信息，配置信息是从上边传入

	client, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Error("init kafka produce failed ,err", err)
		return
	}
	logs.Debug("init kafka success")
	//defer client.Close()不用关了
	return
}

func SendToKafka(data, topic string) (err error) {
	//封装成一个msg然后发出去就可以了。
	msg := &sarama.ProducerMessage{} //创建一个消息生产者实例
	//配了个topic，就会自动生成这个topic，然后就可以根据这个topic消费了。
	msg.Topic = topic //队列topic的名字,，也就是要搜集项目日志的名字，消费端匹配上topic才能消费消息。
	//编码消息，形成了消息
	msg.Value = sarama.StringEncoder(data)
	//调用client发送消息到kafka，返回一个分区id，分区后可以扩容，和分区中的偏移量offset,
	_, _, err = client.SendMessage(msg) //这样就可以将消息发送到kafka了
	if err != nil {
		logs.Error("send message failed ,err:%v,data:%v,topic:%v", err, data, topic)
		return
	}
	//写一个地方，读同样的地方，死循环
	//logs.Debug("send success,pid:%v,offset:%v,topic:%v\n", pid, offset, topic)
	//fmt.Println(msg)
	return
}
