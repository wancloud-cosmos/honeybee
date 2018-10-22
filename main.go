package main

import (
	"validator-monitor/app"
	_ "validator-monitor/routers"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetLogFuncCall(true)
}

func main() {
	app.Start()
	beego.Run()
}
