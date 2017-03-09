package W

import (
	"github.com/jordan-wright/email"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"net/smtp"
)

type SmtpConfig struct {
	Name     string
	Username string
	Password string
	Hostname string
	Port     int
}

func (mc *SmtpConfig) Address() string {
	return mc.Hostname + `:` + I.ToStr(mc.Port)
}

func (mc *SmtpConfig) Auth() smtp.Auth {
	return smtp.PlainAuth(``, mc.Username, mc.Password, mc.Hostname)
}
func (mc *SmtpConfig) From() string {
	return mc.Name + ` <` + mc.Username + `>`
}

// run sendbcc on another goroutine
func (mc *SmtpConfig) SendBCC(bcc []string, subject string, message string) {
	L.Print(`SendBCC started ` + A.StrJoin(bcc, `, `) + `; subject: ` + subject)
	go mc.SendSyncBCC(bcc, subject, message)
}

// run sendAttachbcc on another goroutine
func (mc *SmtpConfig) SendAttachBCC(bcc []string, subject string, message string, files []string) {
	L.Print(`SendAttachBCC started ` + A.StrJoin(bcc, `, `) + `; subject: ` + subject)
	go mc.SendSyncAttachBCC(bcc, subject, message, files)
}

// sendbcc synchronous version, returns error message
func (mc *SmtpConfig) SendSyncBCC(bcc []string, subject string, message string) string {
	return mc.SendSyncAttachBCC(bcc, subject, message, []string{})
}

// sendbcc synchronous version, returns error message
func (mc *SmtpConfig) SendSyncAttachBCC(bcc []string, subject string, message string, files []string) string {
	e := email.NewEmail()
	e.From = mc.From()
	e.To = []string{e.From}
	e.Bcc = bcc
	e.Subject = subject
	attach := A.StrJoin(files, ` `)
	for _, file := range files {
		e.AttachFile(file)
	}
	if attach != `` {
		attach = `; attachments: ` + attach
	}
	e.HTML = []byte(message + `<br/>
<br/>
--<br/>
Sincerely,<br/>
Automated Software<br/>
` + e.From)
	L.Describe(e.Subject, e.Bcc)
	err := e.Send(mc.Address(), mc.Auth())
	if L.IsError(err, `failed to SendBCC`) {
		return err.Error()
	}
	L.Print(`SendAttachBCC completed ` + A.StrJoin(bcc, `, `) + attach + `; subject: ` + subject)
	return ``
}
