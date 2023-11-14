package xMailer

import (
	"fmt"

	"github.com/kokizzu/gotro/L"
	"github.com/wneessen/go-mail"

	"example2/conf"
)

type Dockermailserver struct {
	conf.DockermailserverConf
	client *mail.Client
}

func (m *Dockermailserver) Connect() error {
	if m.client != nil {
		err := m.client.Close()
		L.IsError(err, `Dockermailserver) Connect.Close`)
	}
	var err error
	m.client, err = mail.NewClient(m.DockermailserverHost,
		mail.WithPort(m.DockermailserverPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(m.DockermailserverUser),
		mail.WithPassword(m.DockermailserverPass),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	return err

}

var ErrDockermailserverSendingEmail = fmt.Errorf(`Dockermailserver) SendEmail`)

func (m *Dockermailserver) SendEmail(
	toEmailName map[string]string,
	subject, text, html string) error {
	msg := mail.NewMsg()
	if err := msg.FromFormat(m.DefaultFromName, m.DefaultFromEmail); err != nil {
		return fmt.Errorf("%w: FromFormat: %v", ErrDockermailserverSendingEmail, err)
	}
	if m.UseBcc {
		if err := msg.AddToFormat(m.DefaultFromName, m.DefaultFromEmail); err != nil {
			return fmt.Errorf("%w: AddToFormat: %v", ErrDockermailserverSendingEmail, err)
		}
	}
	if err := msg.ReplyToFormat(m.DefaultFromName, m.ReplyToEmail); err != nil {
		return fmt.Errorf("%w: ReplyToFormat: %v", ErrDockermailserverSendingEmail, err)
	}
	for email, name := range toEmailName {
		if m.UseBcc {
			if err := msg.AddBccFormat(name, email); err != nil {
				return fmt.Errorf("%w: AddBccFormat: %v", ErrDockermailserverSendingEmail, err)
			}
		} else {
			if err := msg.AddToFormat(name, email); err != nil {
				return fmt.Errorf("%w: AddToFormat: %v", ErrDockermailserverSendingEmail, err)
			}
		}
	}
	msg.Subject(subject)
	if text != `` {
		msg.SetBodyString(mail.TypeTextPlain, text)
	}
	if html != `` {
		msg.SetBodyString(mail.TypeTextHTML, html)
	}
	if err := m.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("%w: DialAndSend: %v", ErrDockermailserverSendingEmail, err)
	}
	return nil
}
