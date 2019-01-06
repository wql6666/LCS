package main

import (
	"LCS/logAgent/tailf"
	"errors"
	"fmt"

	"github.com/astaxie/beego/config"
)

//下边的loadConf 函数里new出来的config的类型Configer，查看源代码可知道
func loadCollectConf(conf config.Configer) (err error) {
	var collectConf tailf.CollectConf
	collectConf.LogPath = conf.String("collect::logPath")
	if collectConf.LogPath == "" {
		err = errors.New("invalid collect::logPath")
		return
	}
	collectConf.Topic = conf.String("collect::topic")
	if collectConf.Topic == "" {
		err = errors.New("invalid collect::topic")
		return
	}
	appConfig.CollectConf = append(appConfig.CollectConf, collectConf)
	return
}

var (
	appConfig *Config
)

type Config struct {
	//定义一个struct来存放配置项
	LogLevel    string
	LogPath     string
	ChanSize    int
	KafkaAddr   string
	CollectConf []tailf.CollectConf //所需要搜集的日志的配置信息
	EtcdAddr    string
	EtcdKey     string
}

func loadConf(confType, filename string) (err error) {
	//创建一个config配置的实例,参数是配置文件的格式和具体是哪一个配置文件的路径
	//返回一个配置项的实例，是一个接口，就可以读取配置了
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed err=", err)
		return //有错误别忘了return
	}
	appConfig = &Config{}
	//appConfig是从conf实例中读取出来的,读取配置
	appConfig.LogLevel = conf.String("logs::logLevel")
	if len(appConfig.LogLevel) == 0 {
		appConfig.LogLevel = "debug"
	}
	//logPath指定的是这个程序的日志写到哪个地方。
	appConfig.LogPath = conf.String("logs::logPath")
	if len(appConfig.LogPath) == 0 {
		appConfig.LogPath = "/home/alan8254402/go/src/LCS/logsForTest/logAgent.log"
	}
	appConfig.ChanSize, err = conf.Int("collect::chanSize")
	if err != nil {
		//err用来判断是否读取到配置
		appConfig.ChanSize = 100
	}
	appConfig.KafkaAddr = conf.String("kafka::serverAddr")
	if len(appConfig.KafkaAddr) == 0 {
		//读取不了地址，没法干活，直接返回
		err = fmt.Errorf("invalid kafka addr")
		return
	}
	appConfig.EtcdAddr = conf.String("etcd::etcdAddr")
	if len(appConfig.EtcdAddr) == 0 {
		err = fmt.Errorf("invalid etcd addr")
		return
	}
	appConfig.EtcdKey = conf.String("etcd::configKey")
	if len(appConfig.EtcdKey) == 0 {
		err = fmt.Errorf("invalid EtcdKey")
		return
	}
	//加载需要搜集的日志的配置信息
	err = loadCollectConf(conf)
	if err != nil {
		fmt.Println("loadCollectConf err=", err)
	}

	return
}
