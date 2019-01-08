package config

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
)

var (
	Interval               time.Duration
	ValidatorAddresses     []string
	NodeAddr               string
	MillBlockCheckInterval time.Duration
	MissBlockLimit         int
	MissBlockWindowSize    int
)

func init() {
	inter, err := beego.AppConfig.Int64("interval")
	if nil != err {
		panic(err)
	}
	Interval = time.Second * time.Duration(inter)

	inter, err = beego.AppConfig.Int64("missblock::interval")
	if nil != err {
		panic(err)
	}
	MillBlockCheckInterval = time.Second * time.Duration(inter)

	MissBlockLimit, err = beego.AppConfig.Int("missblock::miss")
	if nil != err {
		panic(err)
	}

	MissBlockWindowSize, err = beego.AppConfig.Int("missblock::size")
	if nil != err {
		panic(err)
	}

	ValidatorAddresses = strings.Split(beego.AppConfig.String("validator::address"), ",")
	beego.Debug("config validator::address ", ValidatorAddresses)

	NodeAddr = beego.AppConfig.String("node::address")
	if "" == NodeAddr {
		panic("node::address empty")
	}
	beego.Debug("config node::address ", NodeAddr)
}
