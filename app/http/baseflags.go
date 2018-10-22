package http

import (
	"github.com/astaxie/beego"
)

var (
	flagName     string
	flagPassword string
	flagChainID  string
)

func init() {
	flagName = beego.AppConfig.String("baseflags::name")
	if "" == flagName {
		panic("baseflags::name invalid")
	}

	flagPassword = beego.AppConfig.String("baseflags::password")
	if "" == flagPassword {
		panic("baseflags::password invalid")
	}

	flagChainID = beego.AppConfig.String("baseflags::chaindid")
	if "" == flagChainID {
		panic("baseflags::chaindid invalid")
	}
}
