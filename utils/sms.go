package utils

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
)

const (
	SMS_HOST = `http://gateway.iems.net.cn/GsmsHttp`
)

var (
	smsUser, smsPassword, smsFrom, smsTo, smsProject string
	smsEnable                                        bool
)

func init() {
	smsUser = beego.AppConfig.String("sms::user")
	smsPassword = beego.AppConfig.String("sms::password")
	smsFrom = beego.AppConfig.String("sms::from")
	smsEnable = beego.AppConfig.DefaultBool("sms::enable", false)
	smsTo = beego.AppConfig.String("sms::to")
	smsProject = beego.AppConfig.String("sms::project")
}

func SendSMS(content string) error {
	subj := fmt.Sprintf("[%s] %s", smsProject, content)

	return doSendSMS(smsTo, subj)
}

func doSendSMS(to string, content string) error {
	if !smsEnable {
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, SMS_HOST, nil)
	if nil != err {
		beego.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("username", smsUser)
	req.Header.Set("password", smsPassword)
	req.Header.Set("from", smsFrom)
	req.Header.Set("to", to)
	req.Header.Set("content", content)

	beego.Debug("send to:", to, "sms")

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		beego.Error(err)
		return err
	}

	beego.Debug("resp of sending sms:", resp)
	return nil
}
