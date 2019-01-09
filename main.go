package main

import (
	"validator-monitor/app"
	_ "validator-monitor/routers"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetLogFuncCall(true)
	beego.SetLevel(utils.LogLevel())
}

func main() {
	app.Start()
	beego.Run()
}
