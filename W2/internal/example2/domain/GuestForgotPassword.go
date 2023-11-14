package domain

import (
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/id64"
	"github.com/vburenin/nsync"

	"example2/conf"
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestForgotPassword.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestForgotPassword.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestForgotPassword.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestForgotPassword.go
//go:generate farify doublequote --file GuestForgotPassword.go

type (
	GuestForgotPasswordIn struct {
		RequestCommon
		Email string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
	}

	GuestForgotPasswordOut struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`

		resetPassUrl string
	}
)

const (
	GuestForgotPasswordAction = `guest/forgotPassword`

	ErrGuestForgotPasswordEmailNotFound          = `forgot password email not found`
	ErrGuestForgotPasswordTriggeredTooFrequently = `forgot password triggered to frequently`
	ErrGuestForgotPasswordModificationFailed     = `forgot password modification failed`
)

var guestForgotPasswordLock = nsync.NewNamedMutex()

func (d *Domain) GuestForgotPassword(in *GuestForgotPasswordIn) (out GuestForgotPasswordOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Email = in.Email

	if !guestForgotPasswordLock.TryLock(in.Email) {
		out.SetError(400, ErrGuestForgotPasswordTriggeredTooFrequently)
		return
	}
	defer guestForgotPasswordLock.Unlock(in.Email)

	if !user.FindByEmail() {
		out.SetError(400, ErrGuestForgotPasswordEmailNotFound)
		return
	}
	out.actor = user.Id
	out.refId = user.Id

	recently := in.TimeNow().Add(-conf.ForgotPasswordThrottleMinute * time.Minute).Unix()
	if user.SecretCodeAt >= recently {
		out.SetError(400, ErrGuestForgotPasswordTriggeredTooFrequently)
		return
	}

	secretCode := id64.SID() + S.RandomCB63(1)
	user.SetSecretCode(secretCode)
	user.SetSecretCodeAt(in.UnixNow())
	hash := S.EncodeCB63(user.Id, 8)

	out.resetPassUrl = in.Host + `/` + GuestResetPasswordAction + `?secretCode=` + secretCode + `&hash=` + hash
	d.runSubtask(func() {
		err := d.Mailer.SendResetPasswordEmail(user.Email, out.resetPassUrl)
		L.IsError(err, `SendResetPasswordEmail`)
		// TODO: insert failed event to clickhouse
	})

	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpdateById() {
		out.SetError(500, ErrGuestForgotPasswordModificationFailed)
		return
	}

	out.Ok = true
	return
}
