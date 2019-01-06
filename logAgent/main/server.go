package main

import (
	"LCS/logAgent/kafka"
	"LCS/logAgent/tailf"
	"time"

	"github.com/astaxie/beego/logs"
)

func serverRun() (err error) {
	//server一直运行，tailf中读取的消息不断的发给kafka
	for {
		msg := tailf.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send to kafka failed ,err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

func sendToKafka(msg *tailf.TextMsg) (err error) {

	//fmt.Printf("read msg:%s,topic：%s\n", msg.Msg, msg.Topic)
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return

}
