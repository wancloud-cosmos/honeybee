package app

import (
	"strings"

	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
)

var (
	host, user, password string
	emailTos             []string
)

func init() {
	host = beego.AppConfig.String("email::host")
	user = beego.AppConfig.String("email::user")
	password = beego.AppConfig.String("email::password")

	emailTos = strings.Split(beego.AppConfig.String("receiptor::email"), ",")
}

func SendMail(tos []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", tos...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	err := gomail.NewDialer(host, 25, user, password).DialAndSend(m)
	if err != nil {
		beego.Error("dial and send email failed,err:", err)
		return err
	}

	beego.Info("send to:", tos, "email success")

	return nil
}
