package tailf

import (
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

//main包的东西放到其他包中去
type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

//很多需要tail，创建一个tail的结构体，即TailFile的对象
type TailObj struct {
	tail     *tail.Tail
	conf     CollectConf
	status   int
	exitChan chan int
}

//创建一个管道，两个字段，一个字段表示这条文本，另一个表示这个文本写入哪个topic
type TextMsg struct {
	Msg   string
	Topic string
}

//用来管理TailObj
type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

//从管道取消息
// msg=<-  加了个等于就不报错了，为什么?取出一个值赋值给msg？
func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	//fmt.Printf("msg=%T,%v", msg, msg)
	return
}

func InitTail(conf []CollectConf, chanSize int) (err error) {
	//定义了的全局变量再用个：=，就会重新赋值，出现空指针！！！
	tailObjMgr = &TailObjMgr{
		//初始化，管道chansize不用写死，通过配置文件传进来
		msgChan: make(chan *TextMsg, chanSize),
	}
	if len(conf) == 0 {
		logs.Error("invalid config for log collect,conf:%v", conf)
		return
	}

	//配置中的日志直接开始日志读取任务，
	// UpdateConfig（），是etcd中更新的，然后开始日志读取任务
	for _, v := range conf {
		createNewTask(v)
	}
	return
}

//更新etcd中动态读取到的配置信息 ,信息指的是要搜集的日志的路径和topic
func UpdateConfig(confs []CollectConf) (err error) {
	//fmt.Printf("tailObjMgr是什么鬼%v,%T\n",tailObjMgr,tailObjMgr)
	//实例化一个TailObjMgr{}
	tailObjMgr = &TailObjMgr{}

	//注意：lock这里需要强化下。之前的有点忘了，表示只能一条一条的状态进行更新配置信息？
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
	//判断新增加的配置是否已经在运行
	for _, oneConf := range confs {
		var isRunning = false
		for _, obj := range tailObjMgr.tailObjs {
			if oneConf.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}
		//如果配置路径已经存在，则不需要再开一个任务搜集日志
		if isRunning {
			continue
		}
		//开启新的收集日志任务，tailf搜集
		createNewTask(oneConf)
	}
	//创建一个tailobj的数组存储更新的obj的信息
	var tailobjs []*TailObj

	for _, obj := range tailObjMgr.tailObjs {

		//添加obj的信息前先进行一次判断，判断之前的配置是否还存在。
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			//注意：这里是啥意思，channel这一部分需要加强
			obj.exitChan <- 1
			continue
		}
		tailobjs = append(tailobjs, obj)
	}
	//更新tailObjMgr.tailObjs的内容
	tailObjMgr.tailObjs = tailobjs

	return
}

func createNewTask(conf CollectConf) {
	//创建&TailObj实例
	tailobj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}
	//开始读写日志（根据配置的路径来读写日志）
	tails, errTail := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen: true, //写完一个日志文件后（可能按大小一个G或者时间来算）然后挪开，需要打开另一个日志文件
		Follow: true, //文件关闭或挪开后会读新的文件
		//Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //记录读到哪个位置了，出现异常时可以定位
		MustExist: false, //日志文件不存在也监控，当日志文件存在时就收集
		Poll:      true,  //不断查询，有的日志时间间隔长，不断查询是否有日志
	})
	//下方return的err有歧义？
	if errTail != nil {
		//fmt.Println("tail file err=", err)
		logs.Error("collect  [%s] log failed,err%v", conf.LogPath, errTail)
		return
	}
	tailobj.tail = tails
	//将obj添加到tailObjMgr结构体中的数组当中去
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, tailobj)

	//因为后边要返回主程序执行，但是这边的tail需要一直进行，
	//所以需要起一个goroutine 一直进行tail
	//来一个配置就起一个goroutine读文件，传什么进去？tails？obj？,根据你下边写的函数来
	//而下边函数，为什么用tailobj，因为还会用到topic这个信息。
	go readFromTail(tailobj)

}

func readFromTail(tailObj *TailObj) {
	//for循环先写死了，按道理写一个配置，给信号，，信号为false，程序退出了，我们就退出
	for true {
		select {
		case line, ok := <-tailObj.tail.Lines: //从tails.lines管道中读一行:
			//管道被关了，ok就为false
			if !ok {
				logs.Warn("tail file close reopen,filename%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			//实例化一个对象，&TextMsg
			textMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}

			tailObjMgr.msgChan <- textMsg
			//注意：此处的tailObj.exitChan怎么理解
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited,conf:%v", tailObj.conf)
			return
		}
	}
}
