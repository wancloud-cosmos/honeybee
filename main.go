package main

import (
	_ "validator-monitor/routers"

	"validator-monitor/app"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetLogFuncCall(true)
}

func main() {
	app.Start()
	beego.Run()
}
