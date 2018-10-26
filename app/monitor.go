package app

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
)

var (
	validatorAddresses []string
	interval           time.Duration
)

func init() {
	validatorAddresses = strings.Split(beego.AppConfig.String("validator::address"), ",")
	beego.Debug("config validator::address ", validatorAddresses)

	inter, err := beego.AppConfig.Int64("interval")
	if nil != err {
		panic(err)
	}

	interval = time.Second * time.Duration(inter)
}

func Watch() {
	go func() {
		var waitTime = interval
		for {
			err := monitorNodes[0].CheckValidator(validatorAddresses)
			if nil != err {
				waitTime = waitTime * 2
				time.Sleep(waitTime)
				continue
			}

			waitTime = interval
			time.Sleep(waitTime)
		}
	}()
}
