package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type BarginsRushController struct {
	beego.Controller
}

func (c *BarginsRushController) Info() {
	c.Data["json"] = "sec kill"
	c.ServeJSON()

	logs.Debug("recv request ")
}
