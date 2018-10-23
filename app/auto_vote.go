package app

import (
	"fmt"
	"time"
	"validator-monitor/app/http"

	"github.com/astaxie/beego"
)

var (
	GovVoter string
)

func init() {
	GovVoter = beego.AppConfig.String("gov::voter")
	if "" == GovVoter {
		panic("gov::voter invalid")
	}
}

func AutoVote(id int64) error {
	emailTitle := fmt.Sprintf("vote proposal-id:%d ", id)
	emailBody := fmt.Sprintf("vote proposal-id:%d ", id)

	err := http.Vote(id, GovVoter, http.OptionYes)
	if nil != err {
		beego.Error(err)

		emailTitle += " failed"
		emailBody += " failed,err:" + err.Error()
	} else {
		emailTitle += " success"
		emailBody += " success"
	}

	//try 3 times
	for i := 0; i < 3; i++ {
		err := SendMail(emailTos, emailTitle, emailBody)
		if nil != err {
			time.Sleep(time.Second * 10)
			continue
		}

		break
	}

	return nil
}
