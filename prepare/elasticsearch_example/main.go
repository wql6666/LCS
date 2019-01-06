package main

import (
	"fmt"
	"gopkg.in/olivere/elastic.v2"
)

type Tweet struct {
	User    string
	Message string
}

func main() {
	//创建客户端，通过客户端来操作es
	client, err := elastic.NewClient(elastic.SetSniff(false),
		elastic.SetURL("http://127.0.0.1:9200/"))
	if err != nil {
		fmt.Println("connect es error", err)
		return
	}
	fmt.Println("connect success!")
	for i := 0; i < 10000; i++ {
		tweet := Tweet{User: "olivere", Message: "Take Five"}
		_, err = client.Index().
			Index("twitter"). //类似数据库的概念
			Type("tweet").    //类似表格的概念
			Id(fmt.Sprintf("%d", i)).
			BodyJson(tweet).Do()
		//调用do表示命令正式执行
		if err != nil {
			//handle error
			panic(err)
			return
		}
	}

	fmt.Println("insert success")
}
