package xMailer

import (
	"github.com/kokizzu/gotro/L"

	"example2/conf"
)

type SendMailFunc func(toEmailName map[string]string, subject, text, html string) error

type Mailer struct {
	SendMailFunc SendMailFunc
}

func (m *Mailer) SendRegistrationEmail(email string, verifyEmailUrl string) error {
	if conf.IsDebug() {
		L.Print(`SendRegistrationEmail`, email, verifyEmailUrl)
	}
	return m.SendMailFunc(
		map[string]string{email: ``},
		`Verify Registration Link`,
		`Hi `+email+`, 

please click this link to verify your registration: 
  `+verifyEmailUrl+`

please ignore this email if you didn't register`,
		`Hi `+email+`, <br><br>
please click this link to verify your registration: <br/>
  <a href="`+verifyEmailUrl+`">`+verifyEmailUrl+`</a><br/><br/>
please ignore this email if you didn't register<br/>`,
	)
}

func (m *Mailer) SendResetPasswordEmail(email string, resetPassUrl string) error {
	if conf.IsDebug() {
		L.Print(`SendResetPasswordEmail`, email, resetPassUrl)
	}
	return m.SendMailFunc(
		map[string]string{email: ``},
		`Reset Password Link`,
		`Hi `+email+`,

please click this link to reset your password: 
`+resetPassUrl+`

please ignore this email if you didn't request reset password link`,
		`Hi `+email+`, <br><br>
please click this link to reset your password: <br/>
<a href="`+resetPassUrl+`">`+resetPassUrl+`</a><br/><br/>
please ignore this email if you didn't request reset password link<br/>`,
	)
}
