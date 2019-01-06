package main

import (
	"fmt"

	"github.com/astaxie/beego/config"
)

func main() {
	//创建一个config的实例,参数是配置文件的格式和具体是哪一个配置文件
	conf, err := config.NewConfig("ini", "./logAgent.conf")
	if err != nil {
		fmt.Println("new config failed err=", err)
		return //有错误别忘了return
	}
	//配置项和名字用两个冒号隔开进行读取,然后读各个配置项
	port, err := conf.Int("server::listen_port")
	if err != nil {
		fmt.Println("read server port failed err", err)
		return
	}
	fmt.Println("prot:", port)
	log_level := conf.String("logs::log_level")
	//返回值里没有错误判断的话，可以自己做一个判断
	if len(log_level) == 0 { //没读出来，则长度为0
		log_level = "debug"
	}
	fmt.Println("log_level", log_level)
	log_path := conf.String("logs::log_path")
	fmt.Println("log_path", log_path)

}
