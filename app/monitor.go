package app

import (
	"time"

	"validator-monitor/app/config"
)

func Watch() {
	MonitorValJailOrNot()
	MonitorValMissBlock()
}

//monitor if the validator continuously miss block
func MonitorValMissBlock() {
	go func() {
		var waitTime = config.MillBlockCheckInterval
		for {
			monitorNode.DidMissBlock(config.ValidatorAddresses[0])
			time.Sleep(waitTime)
		}
	}()
}

//monitor if the validator is jailed
func MonitorValJailOrNot() {
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
