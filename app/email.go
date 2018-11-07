package app

import (
	"net/smtp"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

var (
	host, user, password, to string
)

func init() {
	host = beego.AppConfig.String("email::host")
	user = beego.AppConfig.String("email::user")
	password = beego.AppConfig.String("email::password")
	to = beego.AppConfig.String("receiptor::email")
}

func SendMail3Times(subject, body string) {
	for i := 0; i < 3; i++ {
		err := SendMail(subject, body)
		if nil != err {
			time.Sleep(time.Second * 5)
			continue
		}

		return
	}
}

func SendMail(subject, body string) error {
	err := doSendMail(user, password, host, to, subject, body, "")
	if nil != err {
		beego.Error(err)
		return err
	}

	beego.Info("send to:", to, "email success")

	return nil
}

func doSendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
