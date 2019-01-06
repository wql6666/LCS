package main

import (
	"fmt"

	"LCS/logAgent/tailf"

	"LCS/logAgent/kafka"

	"github.com/astaxie/beego/logs"
)

func main() {
	filename := "/home/alan8254402/go/src/LCS/config/logAgent.conf"
	//加载配置文件,啦啦试验下git
	err := loadConf("ini", filename)
	if err != nil {
		fmt.Println("load conf failed,err=", err)
		panic("load conf failed")
		return
	}
	fmt.Println("loadConf succ")

	err = initLogger()
	if err != nil {
		fmt.Printf("initLogger err%v\n", err)
		panic("load logger failed ")
		return
	}
	logs.Debug("load conf success,conf:%v", appConfig)
	fmt.Println("initLogger succ")

	collectConf, err := initEtcd(appConfig.EtcdAddr, appConfig.EtcdKey)
	if err != nil {
		logs.Error("init etcd failed,err%v", err)
	}
	fmt.Println("init etcd success", collectConf)
	logs.Debug("init etcd success")

	err = tailf.InitTail(collectConf, appConfig.ChanSize)
	if err != nil {
		logs.Error("init tail failed err%v", err)
		return
	}
	logs.Debug("initTailf success")
	fmt.Println("initTailf succ")

	err = kafka.InitKafka(appConfig.KafkaAddr)
	if err != nil {
		logs.Error("initKafka failed err%v", err)
		return
	}

	logs.Debug("init all success")
	fmt.Println("init all succ")

	//开个goroutine写日志到终端测试。
	//go func() {
	//	var count int
	//	for {
	//		count++
	//		logs.Debug("test for logger %d", count)
	//		time.Sleep(time.Second)
	//	}
	//}()
	//time.Sleep(time.Second * 3)

	err = serverRun()
	if err != nil {
		logs.Error("serverRun failed err%v", err)
		return
	}
	logs.Info("program exit")

}
