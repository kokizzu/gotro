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

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestResendVerificationEmail.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestResendVerificationEmail.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestResendVerificationEmail.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestResendVerificationEmail.go
//go:generate farify doublequote --file GuestResendVerificationEmail.go

type (
	GuestResendVerificationEmailIn struct {
		RequestCommon
		Email string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
	}
	GuestResendVerificationEmailOut struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`

		verifyEmailUrl string
	}
)

const (
	GuestResendVerificationEmailAction = `guest/resendVerificationEmail`

	ErrGuestResendVerificationEmailUserNotFound           = `user not found`
	ErrGuestResendVerificationEmailTriggeredTooFrequently = `resend verification triggered to frequently`
	ErrGuestResendVerificationEmailUserAlreadyVerified    = `user already verified`
	ErrGuestResendVerificationEmailModificationFailed     = `resend verification modification failed`
)

var guestResendVerificationEmailLock = nsync.NewNamedMutex()

func (d *Domain) GuestResendVerificationEmail(in *GuestResendVerificationEmailIn) (out GuestResendVerificationEmailOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Email = in.Email

	if !guestResendVerificationEmailLock.TryLock(in.Email) {
		out.SetError(400, ErrGuestResendVerificationEmailTriggeredTooFrequently)
		return
	}

	if !user.FindByEmail() {
		guestResendVerificationEmailLock.Unlock(in.Email)
		out.SetError(400, ErrGuestResendVerificationEmailUserNotFound)
		return
	}
	out.actor = user.Id
	out.refId = user.Id

	if user.VerifiedAt > 0 {
		guestResendVerificationEmailLock.Unlock(in.Email)
		out.SetError(400, ErrGuestResendVerificationEmailUserAlreadyVerified)
		return
	}

	recently := in.TimeNow().Add(-conf.ResendVerificationEmailThrottleMinute * time.Minute).Unix()
	if user.SecretCodeAt >= recently {
		guestResendVerificationEmailLock.Unlock(in.Email)
		out.SetError(400, ErrGuestResendVerificationEmailTriggeredTooFrequently)
		return
	}

	secretCode := id64.SID() + S.RandomCB63(1)
	user.SetSecretCode(secretCode)
	user.SetSecretCodeAt(in.UnixNow())
	hash := S.EncodeCB63(user.Id, 8)

	out.verifyEmailUrl = in.Host + `/` + GuestVerifyEmailAction + `?secretCode=` + secretCode + `&hash=` + hash
	d.runSubtask(func() {
		defer guestResendVerificationEmailLock.Unlock(in.Email)
		err := d.Mailer.SendRegistrationEmail(user.Email, out.verifyEmailUrl)
		L.IsError(err, `SendRegistrationEmail`)
		// TODO: insert failed event to clickhouse
	})

	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpdateById() {
		guestResendVerificationEmailLock.Unlock(in.Email)
		out.SetError(500, ErrGuestResendVerificationEmailModificationFailed)
		return
	}

	out.Ok = true
	return
}
