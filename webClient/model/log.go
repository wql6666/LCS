package model

import (
	"LCS/logAgent/tailf"
	"context"
	"encoding/json"
	"time"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go.etcd.io/etcd/clientv3"
)

var (
	etcdClient *clientv3.Client
)

type LogInfo struct {
	AppId      int    `db:"app_id"`
	AppName    string `db:"app_name"`
	LogId      int    `db:"log_id"`
	Status     int    `db:"status"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
}

func InitEtcd(client *clientv3.Client) {
	etcdClient = client
}

//获取全部所需要搜集的日志的配置信息
func GetAllLogInfo() (loglist []LogInfo, err error) {
	err = Db.Select(&loglist,
		//注意：这里的a,b表怎么区分的？事先有声明吗
		"select a.app_id, b.app_name, a.create_time, a.topic, a.log_id, a.status, a.log_path from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func CreateLog(info *LogInfo) (err error) {

	tx, err := Db.Begin()
	if err != nil {
		logs.Warn("CreateApp failed, Db.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	var appId []int
	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", info.AppName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}
	//注意：info.AppId = appId[0]这是啥意思？
	info.AppId = appId[0]
	result, err := tx.Exec("insert into tbl_log_info(app_id, log_path, topic)values(?, ?, ?)",
		info.AppId, info.LogPath, info.Topic)

	if err != nil {
		logs.Warn("CreateApp failed, Db.Exec error:%v", err)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		logs.Warn("CreateApp failed, Db.LastInsertId error:%v", err)
		return
	}

	return
}

//将日志的配置信息发送到etcd存储，这里是通过数据库的信息发送给etcd，数据库的信息来源于网页端的输入。
func SetLogConfToEtcd(etcdKey string, info *LogInfo) (err error) {
	//
	var logConfArr []tailf.CollectConf
	logConfArr = append(
		logConfArr,
		tailf.CollectConf{
			LogPath: info.LogPath,
			Topic:   info.Topic,
		},
	)

	data, err := json.Marshal(logConfArr)
	if err != nil {
		logs.Warn("marshal failed, err:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	_, err = etcdClient.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		logs.Warn("Put failed, err:%v", err)
		return
	}

	logs.Debug("put etcd succ, data:%v", string(data))
	return
}
