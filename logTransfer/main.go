package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

func main() {
	//加载配置
	err := initConfig("ini", "/home/alan8254402/go/src/LCS/config/log_tansfer.conf")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(appConfig)
	//logs.Debug("init config succ")//日志还没初始化，还不能写到日志中去，只能写到终端

	err = initLogger(appConfig.LogPath, appConfig.LogLevel)
	if err != nil {
		panic(err)
		return
	}
	logs.Debug("init logger succ")

	err = initKafka(appConfig.KafkaAddr, appConfig.KafkaTopic)
	if err != nil {
		logs.Error("init kafka failed,err=", err)
		return
	}
	logs.Debug("initKafka succ")

	err = initES(appConfig.EsAddr)
	if err != nil {
		logs.Error("init es failed ,err=", err)
		return
	}
	logs.Debug("init ES succ")
	fmt.Println("init ES succ !!")

	err = run()
	if err != nil {
		logs.Error("run failed ,err=", err)
		return
	}
	logs.Warn("warning,logtransfer is exited")

}
