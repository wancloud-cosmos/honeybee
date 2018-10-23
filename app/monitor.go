package app

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
)

var (
	validatorAddresses []string
	interval           int64
)

func init() {
	validatorAddresses = strings.Split(beego.AppConfig.String("validator::address"), ",")
	beego.Debug("config validator::address ", validatorAddresses)

	var err error
	interval, err = beego.AppConfig.Int64("interval")
	if nil != err {
		panic(err)
	}
}

func Watch() {
	go func() {
		for {
			for _, n := range monitorNodes {
				n.CheckValidator(validatorAddresses)
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	}()
}
