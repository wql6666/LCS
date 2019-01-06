package main

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

func main() {
	config := make(map[string]interface{})
	config["filename"] = "./logs/logcollect.log" //实际上从配置文件库中读出文件路径
	config["level"] = logs.LevelInfo             //配置log级别

	configStr, err := json.Marshal(config) //json把map序列化后是一个byte数组
	if err != nil {
		fmt.Println("config marshal err", err)
		return
	}
	//初始化log
	logs.SetLogger(logs.AdapterFile, string(configStr))
	//使用日志，调试的日志
	logs.Debug("this is a test,my name is %s", "stu01")
	logs.Trace("this is a test,my name is %s", "stu02")
	logs.Warn("this is a test,my name is %s", "stu03")

}
