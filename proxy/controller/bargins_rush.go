package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/wx7217242/bargains-rush/proxy/service"
)

type BarginsRushController struct {
	beego.Controller
}

func (p *BarginsRushController) Rush() {

	result := make(map[string]interface{})

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	result["code"] = 0
	result["message"] = "success"
}

func (p *BarginsRushController) Info() {

	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid product_id"

		logs.Error("invalid request, get product_id failed, err:%v", err)
		return
	}

	data, code, err := service.Info(productId)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()

		logs.Error("invalid request, get product_id failed, err:%v", err)
		return
	}

	result["data"] = data
}
