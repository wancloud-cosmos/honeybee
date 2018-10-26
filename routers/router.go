package routers

import (
	"validator-monitor/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/vote", &controllers.VoteController{})
}
