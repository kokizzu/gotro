package xMailer

import (
	"fmt"

	"github.com/kokizzu/gotro/L"
	"github.com/wneessen/go-mail"

	"example2/conf"
)

type Mailhog struct {
	conf.MailhogConf
	client *mail.Client
}

func NewMailhog(cfg conf.MailhogConf) (*Mailhog, error) {
	res := &Mailhog{
		MailhogConf: cfg,
	}
	err := res.Connect()
	return res, err
}

func (m *Mailhog) Connect() error {
	if m.client != nil {
		err := m.client.Close()
		L.IsError(err, `Mailhog) Connect.Close`)
	}
	var err error
	// if without mailhog, adds
	// mail.WithSMTPAuth(mail.SMTPAuthPlain),
	// mail.WithUsername("user")
	// mail.WithPassword("pwd"))
	m.client, err = mail.NewClient(m.MailhogHost,
		mail.WithPort(m.MailhogPort),
		mail.WithTLSPolicy(mail.NoTLS),
	)
	return err

}

var ErrMailhogSendingEmail = fmt.Errorf(`Mailhog) SendEmail`)

func (m *Mailhog) SendEmail(
	toEmailName map[string]string,
	subject, text, html string) error {
	msg := mail.NewMsg()
	if err := msg.FromFormat(m.DefaultFromName, m.DefaultFromEmail); err != nil {
		return fmt.Errorf("%w: FromFormat: %v", ErrMailhogSendingEmail, err)
	}
	if m.UseBcc {
		if err := msg.AddToFormat(m.DefaultFromName, m.DefaultFromEmail); err != nil {
			return fmt.Errorf("%w: AddToFormat: %v", ErrMailhogSendingEmail, err)
		}
	}
	if err := msg.ReplyToFormat(m.DefaultFromName, m.ReplyToEmail); err != nil {
		return fmt.Errorf("%w: ReplyToFormat: %v", ErrMailhogSendingEmail, err)
	}
	for email, name := range toEmailName {
		if m.UseBcc {
			if err := msg.AddBccFormat(name, email); err != nil {
				return fmt.Errorf("%w: AddBccFormat: %v", ErrMailhogSendingEmail, err)
			}
		} else {
			if err := msg.AddToFormat(name, email); err != nil {
				return fmt.Errorf("%w: AddToFormat: %v", ErrMailhogSendingEmail, err)
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
		return fmt.Errorf("%w: DialAndSend: %v", ErrMailhogSendingEmail, err)
	}
	return nil
}
