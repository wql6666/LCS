package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	filename := "./a" //指定读哪个日志文件
	//实例化一个tail
	tails, err := tail.TailFile(filename, tail.Config{
		ReOpen: true, //读完一个日志文件后（可能按大小一个G或者时间来算）然后挪开，需要打开另一个日志文件
		Follow: true, //文件关闭或挪开后会读新的文件
		//Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //记录读到哪个位置了，出现异常时可以定位
		MustExist: false, //日志文件不存在也监控，当日志文件存在时就收集
		Poll:      true,  //不断查询，有的日志时间间隔长，不断查询是否有日志
	})
	if err != nil {
		fmt.Println("tail file err=", err)
		return
	}
	var (
		msg *tail.Line //定义一行数据
		ok  bool
	)
	for {
		//ok，用来判断管段是否关闭，关闭了ok为false
		msg, ok = <-tails.Lines //从tails.lines管道中读一行
		if !ok {
			fmt.Printf("tail file close reopen,filename%s\n",
				tails.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("msg", msg)
	}
}
