package base_c

import (
	"gopkg.in/gomail.v2"
	"one/infra"
)

var email *gomail.Dialer

func Email() *gomail.Dialer {
	return email
}

type EmailStarter struct {
	infra.BaseStarter
}

func (s *EmailStarter) Init(ctx infra.StarterContext) {
	user := ctx.Props().GetDefault("email.user", "")
	pwd := ctx.Props().GetDefault("email.pwd", "")
	host := ctx.Props().GetDefault("email.host", "192.016.10.10")
	port := ctx.Props().GetIntDefault("email.port", 2525)
	email = gomail.NewDialer(host, port, user, pwd)
}
