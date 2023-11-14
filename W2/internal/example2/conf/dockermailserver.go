package conf

import (
	"os"

	"github.com/kokizzu/gotro/S"
)

type DockermailserverConf struct {
	MailerConf
	DockermailserverHost string
	DockermailserverPort int
	DockermailserverUser string
	DockermailserverPass string
}

func EnvDockermailserver() DockermailserverConf {
	return DockermailserverConf{
		DockermailserverHost: os.Getenv("DOCKERMAILSERVER_HOST"),
		DockermailserverPort: S.ToInt(os.Getenv("DOCKERMAILSERVER_PORT")),
		DockermailserverUser: os.Getenv("DOCKERMAILSERVER_USER"),
		DockermailserverPass: os.Getenv("DOCKERMAILSERVER_PASS"),
		MailerConf:           EnvMailer(),
	}
}
