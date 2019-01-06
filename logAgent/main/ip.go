package main

import (
	"fmt"
	"net"
	//"github.com/astaxie/beego/logs"
)

var (
	localIpArray []string //多个，然后就存到切片里边，后边可以遍历
)

//获取本地ip
func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("get local ip failed,%v", err))
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localIpArray = append(localIpArray, ipNet.IP.String())
			}
		}

	}
	fmt.Println("localIpArry", localIpArray)
}
