package app

import (
	"fmt"
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
	err := http.Vote(id, GovVoter, http.OptionYes)
	if nil != err {
		beego.Error(err)
		return err
	}
	//TODO notify admin by sending email

	emailTitle := fmt.Sprintf("vote proposal-id:%d success", id)
	emailBody := fmt.Sprintf("vote proposal-id:%d success", id)
	beego.Error(emailBody)
	SendMail(emailTos, emailTitle, emailBody)
	return nil
}
