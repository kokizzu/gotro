package conf

import (
	"os"

	"github.com/kokizzu/gotro/S"
)

type MailhogConf struct {
	MailerConf
	MailhogHost string
	MailhogPort int
}

func EnvMailhog() MailhogConf {
	return MailhogConf{
		MailhogHost: os.Getenv("MAILHOG_HOST"),
		MailhogPort: S.ToInt(os.Getenv("MAILHOG_PORT")),
		MailerConf:  EnvMailer(),
	}
}
