package domain

import (
	"github.com/kokizzu/gotro/S"

	"example2/conf"
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestResetPassword.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestResetPassword.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestResetPassword.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestResetPassword.go
//go:generate farify doublequote --file GuestResetPassword.go

type (
	GuestResetPasswordIn struct {
		RequestCommon
		SecretCode string `json:"secretCode" form:"secretCode" query:"secretCode" long:"secretCode" msg:"secretCode"`
		Hash       string `json:"hash" form:"hash" query:"hash" long:"hash" msg:"hash"`
		Password   string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	GuestResetPasswordOut struct {
		ResponseCommon
		Ok bool `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
	}
)

const (
	GuestResetPasswordAction = `guest/resetPassword`

	ErrGuestResetPasswordInvalidHash        = `invalid hash`
	ErrGuestResetPasswordTooShort           = `password too short`
	ErrGuestResetPasswordUserNotFound       = `user not found`
	ErrGuestResetPasswordWrongSecret        = `wrong secret code`
	ErrGuestResetPasswordExpiredLink        = `expired link`
	ErrGuestResetPasswordModificationFailed = `reset password modification failed`
)

func (d *Domain) GuestResetPassword(in *GuestResetPasswordIn) (out GuestResetPasswordOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	userId, ok := S.DecodeCB63[uint64](in.Hash)
	out.refId = userId
	if !ok {
		out.SetError(400, ErrGuestResetPasswordInvalidHash)
		return
	}
	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Id = userId
	if !user.FindById() {
		out.SetError(400, ErrGuestResetPasswordUserNotFound)
		return
	}
	out.actor = user.Id

	if len(in.Password) < minPassLength {
		out.SetErrorf(400, ErrGuestResetPasswordTooShort)
		return
	}
	if user.SecretCode == `` ||
		user.SecretCodeAt == 0 ||
		in.UnixNow()-user.SecretCodeAt > conf.ForgotPasswordExpireMinute*60 {
		out.SetError(404, ErrGuestResetPasswordExpiredLink)
		return
	}
	if user.SecretCode != in.SecretCode {
		out.SetError(400, ErrGuestResetPasswordWrongSecret)
		return
	}
	// also verify the user if not verified yet
	if user.VerifiedAt == 0 {
		user.SetVerifiedAt(in.UnixNow())
	}
	user.SetSecretCode(``)
	user.SetSecretCodeAt(0)
	user.SetEncryptedPassword(in.Password, in.UnixNow())

	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpdateById() {
		out.SetError(500, ErrGuestResetPasswordModificationFailed)
		return
	}
	out.Ok = true
	return
}
