package config

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
)

var (
	Interval           time.Duration
	ValidatorAddresses []string
	NodeAddr           string
)

func init() {
	inter, err := beego.AppConfig.Int64("interval")
	if nil != err {
		panic(err)
	}
	Interval = time.Second * time.Duration(inter)

	ValidatorAddresses = strings.Split(beego.AppConfig.String("validator::address"), ",")
	beego.Debug("config validator::address ", ValidatorAddresses)

	NodeAddr = beego.AppConfig.String("node::address")
	if "" == NodeAddr {
		panic("node::address empty")
	}
	beego.Debug("config node::address ", NodeAddr)
}
