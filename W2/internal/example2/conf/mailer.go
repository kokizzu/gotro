package conf

import (
	"os"

	"github.com/kokizzu/gotro/X"
)

type MailerConf struct {
	DefaultFromEmail string
	DefaultFromName  string
	ReplyToEmail     string
	UseBcc           bool
	DefaultMailer    string
}

func EnvMailer() MailerConf {
	return MailerConf{
		DefaultFromEmail: os.Getenv("MAILER_DEFAULT_FROM_EMAIL"),
		DefaultFromName:  os.Getenv("MAILER_DEFAULT_FROM_NAME"),
		ReplyToEmail:     os.Getenv("MAILER_REPLY_TO_EMAIL"),
		UseBcc:           X.ToBool(os.Getenv("MAILER_USE_BCC")),
		DefaultMailer:    os.Getenv("MAILER_DEFAULT"),
	}
}
