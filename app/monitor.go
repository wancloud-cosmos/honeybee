package app

import (
	"time"

	"validator-monitor/app/config"
)

func Watch() {
	go func() {
		var waitTime = config.Interval
		for {
			err := monitorNode.CheckValidator(config.ValidatorAddresses)
			if nil != err {
				waitTime = waitTime * 2
				time.Sleep(waitTime)
				continue
			}

			waitTime = config.Interval
			time.Sleep(waitTime)
		}
	}()
}
