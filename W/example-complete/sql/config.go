package sql

import (
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
)

var WEBMASTER_EMAILS = M.SS{
	`kiswono@gmail.com`: `Kiswono Prayogo`,
	// TODO: TODO_CHANGE_THIS
}

const SUPPORT_EMAIL = `TODO_CHANGE_THIS`
const DEBUGGER_EMAIL = `kiswono+gotro-example-complete@gmail.com`

var PROJECT_NAME string
var DOMAIN string

var PG *Pg.RDBMS

func init() {
	PG = Pg.NewConn(`TODO_CHANGE_DB`, `TODO_CHANGE_DB`)
	W.Mailers = map[string]*W.SmtpConfig{
		``: {
			Name:     `Mailer Daemon`,
			Username: `TODO_CHANGE_THIS`,
			Password: `TODO_CHANGE_THIS`,
			Hostname: `smtp.gmail.com`,
			Port:     587,
		},
	}
}
