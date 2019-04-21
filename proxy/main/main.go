package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/wx7217242/bargains-rush/proxy/router"
)

func main() {

	err := initConf()
	if err != nil {
		panic(err)
		return
	}

	err = initApp()
	if err != nil {
		panic(err)
		return
	}

	beego.Run()

	logs.Info("run over...")
}
