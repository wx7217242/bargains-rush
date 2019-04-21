package router

import (
	"github.com/astaxie/beego"
	"github.com/wx7217242/bargains-rush/proxy/controller"
)

func init() {
	beego.Router("/bargins-rush/info", &controller.BarginsRushController{}, "*:Info")
	beego.Router("/bargins-rush/rush", &controller.BarginsRushController{}, "*:Rush")
}
