package app

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

const (
	SMS_HOST = `http://gateway.iems.net.cn/GsmsHttp`
)

var (
	smsUser, smsPassword string
)

func init() {
	smsUser = beego.AppConfig.String("sms::user")
	smsPassword = beego.AppConfig.String("sms::password")
}

func SendSMS(tos []string, content string) error {
	req, err := http.NewRequest(http.MethodGet, SMS_HOST, nil)
	if nil != err {
		beego.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("username", smsUser)
	req.Header.Set("password", smsPassword)
	req.Header.Set("from", "80")

	beego.Debug("send to:", strings.Join(tos, ","), "sms")

	req.Header.Set("to", strings.Join(tos, ","))
	req.Header.Set("content", content)

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		beego.Error(err)
		return err
	}

	beego.Debug("resp of sending sms:", resp)
	return nil
}
