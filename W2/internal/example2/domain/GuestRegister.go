package domain

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/id64"

	"example2/model/mAuth/rqAuth"
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestRegister.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestRegister.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestRegister.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestRegister.go
//go:generate farify doublequote --file GuestRegister.go

type (
	GuestRegisterIn struct {
		RequestCommon
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	GuestRegisterOut struct {
		ResponseCommon
		User rqAuth.Users `json:"user" form:"user" query:"user" long:"user" msg:"user"`

		verifyEmailUrl string
	}
)

const (
	GuestRegisterAction = `guest/register`

	ErrGuestRegisterEmailInvalid       = `email must be valid`
	ErrGuestRegisterPasswordTooShort   = `password must be at least 12 characters`
	ErrGuestRegisterEmailUsed          = `email already used`
	ErrGuestRegisterUserCreationFailed = `user creation failed`

	minPassLength = 12
)

func (d *Domain) GuestRegister(in *GuestRegisterIn) (out GuestRegisterOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	in.Email = S.Trim(S.ValidateEmail(in.Email))
	if in.Email == `` {
		out.SetError(400, ErrGuestRegisterEmailInvalid)
		return
	}
	if len(in.Password) < minPassLength {
		out.SetErrorf(400, ErrGuestRegisterPasswordTooShort)
		return
	}
	exists := rqAuth.NewUsers(d.AuthOltp)
	exists.Email = in.Email
	if exists.FindByEmail() {
		out.SetError(400, ErrGuestRegisterEmailUsed)
		return
	}
	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Email = in.Email
	user.SetEncryptedPassword(in.Password, in.UnixNow())
	user.SecretCode = id64.SID() + S.RandomCB63(1)

	user.SetUpdatedAt(in.UnixNow())
	user.SetCreatedAt(in.UnixNow())

	if !d.Superadmins[in.Email] {
		user.SetRole(UserSegment)
	}
	if !user.DoInsert() {
		out.SetError(500, ErrGuestRegisterUserCreationFailed)
		return
	}
	out.actor = user.Id
	out.refId = user.Id

	// send verification link
	hash := S.EncodeCB63(user.Id, 8)
	out.verifyEmailUrl = in.Host + `/` + GuestVerifyEmailAction + `?secretCode=` + user.SecretCode + `&hash=` + hash

	user.CensorFields()
	out.User = user.Users

	d.runSubtask(func() {
		err := d.Mailer.SendRegistrationEmail(user.Email, out.verifyEmailUrl)
		L.IsError(err, `SendRegistrationEmail`)
		// TODO: insert failed event to clickhouse
	})
	return
}
