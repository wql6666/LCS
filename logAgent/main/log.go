package main

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

//转换日志类型，原来定义为string，需要转换成int
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
func initLogger() (err error) {
	//beego的日志配置，传的是个jason串，怕拼串出错，用到一个map来存，然后marshal
	config := make(map[string]interface{})
	config["filename"] = appConfig.LogPath                //实际上从配置文件库中读出文件路径
	config["level"] = convertLogLevel(appConfig.LogLevel) //配置log级别,为什么要转换？
	//json把存储日志配置config的map序列化后是一个byte数组，
	// 因为后边setlogger，创建日志，需要的日志的配置config需要的类型是string
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger configStr marshal err", err)
		return
	}
	//初始化log,,logs.AdapterFile表示日志打到文件中去了,
	// AdapterConsole就表示日志写到终端里去,还可以直接写到es里边。

	logs.SetLogger(logs.AdapterFile, string(configStr))

	return

}
