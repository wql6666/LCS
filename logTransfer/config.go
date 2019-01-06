package main

import (
	"fmt"

	"github.com/astaxie/beego/config"
)

type logConfig struct {
	KafkaAddr  string
	EsAddr     string
	LogPath    string
	LogLevel   string
	KafkaTopic string
}

var (
	appConfig *logConfig
)

func initConfig(confTpey string, filename string) (err error) {
	//创建一个configer的接口,参数是配置文件的格式和具体是哪一个配置文件
	conf, err := config.NewConfig("ini", filename)
	if err != nil {
		fmt.Println("new config failed err=", err)
		return //有错误别忘了return
		//time.Sleep(time.Second)
	}
	//创建一个logconfig的实例
	appConfig = &logConfig{}
	//appConfig是从conf实例中读取出来的,读取配置
	appConfig.LogLevel = conf.String("logs::logLevel")
	if len(appConfig.LogLevel) == 0 {
		appConfig.LogLevel = "debug"
	}
	appConfig.LogPath = conf.String("logs::logPath")
	if len(appConfig.LogPath) == 0 {
		appConfig.LogPath = "/home/alan8254402/go/src/LCS/logsForTest/log_transfer.log"
	}
	appConfig.KafkaAddr = conf.String("kafka::serverAddr")
	if len(appConfig.KafkaAddr) == 0 {
		//读取不了地址，没法干活，直接返回
		err = fmt.Errorf("invalid kafka addr")
		return
	}
	appConfig.EsAddr = conf.String("es::esAddr")
	if len(appConfig.EsAddr) == 0 {
		err = fmt.Errorf("invalid EsAddr addr")
		return
	}
	appConfig.KafkaTopic = conf.String("kafka::topic")
	if len(appConfig.KafkaTopic) == 0 {
		err = fmt.Errorf("invalid kafka topic")
		return
	}
	return
}
