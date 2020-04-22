package helper

import (
	"gopkg.in/gomail.v2"
	"one/infra/base_c"
)

// 发送邮件
func SendEmail(from, cc, ccName, subject, body, attach string, to interface{}) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	if to.(string) != "" {
		m.SetHeader("To", to.(string))
	} else {
		m.SetHeader("To", to.([]string)...)
	}
	if cc != "" {
		m.SetAddressHeader("Cc", cc, ccName)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attach != "" {
		m.Attach(attach)
	}
	// Send the email to Bob, Cora and Dan.
	if err := base_c.Email().DialAndSend(m); err != nil {
		return err
	}
	return nil
}
