package AppController

import (
	"LCS/webClient/model"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type AppController struct {
	beego.Controller
}

//项目列表
func (p *AppController) AppList() {
	//注意：这里日志去哪里查看？
	logs.Debug("enter applist index controller")
	//网页的框架
	p.Layout = "layout/layout.html"
	appList, err := model.GetAllAppInfo()
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("服务器繁忙")
		p.TplName = "app/error.html"

		logs.Warn("get app list failed, err:%v", err)
		return
	}

	logs.Debug("get app list succ, data:%v", appList)
	p.Data["applist"] = appList

	p.TplName = "app/index.html"
}

//项目申请
func (p *AppController) AppApply() {
	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"
	p.TplName = "app/apply.html"
}

//创建项目
func (p *AppController) AppCreate() {

	logs.Debug("enter index controller")
	appName := p.GetString("app_name")
	appType := p.GetString("app_type")
	developPath := p.GetString("develop_path")
	ipstr := p.GetString("iplist")

	p.Layout = "layout/layout.html"

	if len(appName) == 0 || len(appType) == 0 || len(developPath) == 0 || len(ipstr) == 0 {
		p.Data["Error"] = fmt.Sprintf("非法参数")
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}
	//创建一个model.AppInfo{}对象
	appInfo := &model.AppInfo{}
	appInfo.AppName = appName
	appInfo.AppType = appType
	appInfo.DevelopPath = developPath
	//注意：ip这里再一次用“，”分开，为啥？几台服务器搜集日志，直接给ip地址就行，就能定位到，服务器？
	appInfo.IP = strings.Split(ipstr, ",")

	err := model.CreateApp(appInfo)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("创建项目失败，数据库繁忙")
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}
	//插入数据完成后重新定位到app/list页面
	p.Redirect("/app/list", 302)
}
