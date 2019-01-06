package main

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelInfo
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

//初始化日志
func initLogger(logPath string, logLevel string) (err error) {
	config := make(map[string]interface{})
	config["filename"] = logPath                //实际上从配置文件库中读出文件路径
	config["level"] = convertLogLevel(logLevel) //配置log级别,为什么要转换？
	//json把map序列化后是一个byte数组
	// 因为后边setlogger，创建日志，需要的日志的配置config需要的类型是string
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger configStr marshal err", err)
		return
	}
	//初始化log
	logs.SetLogger(logs.AdapterFile, string(configStr))

	return

}
