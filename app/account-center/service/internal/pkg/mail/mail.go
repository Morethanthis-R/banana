package mailUtils

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"banana/app/account-center/service/internal/conf"
)

func Send(conf *conf.Mail,token string, target ...string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", conf.Username)
	m.SetHeader("To", target...)

	m.SetHeader("Subject", "PeachCloud官方邮件")
	m.SetBody("text/html", generateBody(token))
	d := gomail.NewDialer(
		conf.SmtpHost, 465,
		conf.Username, conf.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
