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
		var waitTime = time.Duration(interval)
		for {
			err := monitorNodes[0].CheckValidator(validatorAddresses)
			if nil != err {
				waitTime = waitTime * 2
				time.Sleep(time.Second * waitTime)
				continue
			}

			waitTime = time.Duration(interval)
			time.Sleep(time.Second * waitTime)
		}
	}()
}
