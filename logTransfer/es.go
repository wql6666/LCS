package main

import (
	"fmt"

	"gopkg.in/olivere/elastic.v2"
)

var (
	esClient *elastic.Client
)

type LogMessae struct {
	App     string
	Topic   string
	Message string
}

func initES(addr string) (err error) {

	//注意：后边的参数的含义，弄清楚。创建客户端，通过客户端来操作es
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		fmt.Println("connect es error", err)
		return
	}
	esClient = client
	return
	/*
		fmt.Println("connect success!")
		for i:=0;i<10000;i++{
			tweet := Tweet{User: "olivere", Message: "Take Five"}
			_,err = client.Index().
				Index("twitter"). //类似数据库的概念
				Type("tweet").//类似表格的概念
				Id(fmt.Sprintf("%d",i)).
				BodyJson(tweet).Do()
			//调用do表示命令正式执行
			if err != nil {
				//handle error
				panic(err)
				return
			}
		}

		fmt.Println("insert success")
	*/

}

func sendToES(topic string, data []byte) (err error) {
	//创建一个LogMessae{}对象
	msg := &LogMessae{}
	msg.Topic = topic
	msg.Message = string(data)
	//将msg传到ES，链式操作，看上边知道含义。
	_, err = esClient.Index().Index(topic).Type(topic).BodyJson(msg).Do()
	if err != nil {
		panic(err)
		return
	}
	return

}
