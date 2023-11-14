package domain

import (
	"github.com/kokizzu/gotro/S"

	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestVerifyEmail.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestVerifyEmail.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestVerifyEmail.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestVerifyEmail.go
//go:generate farify doublequote --file GuestVerifyEmail.go

type (
	GuestVerifyEmailIn struct {
		RequestCommon
		SecretCode string `json:"secretCode" form:"secretCode" query:"secretCode" long:"secretCode" msg:"secretCode"`
		Hash       string `json:"hash" form:"hash" query:"hash" long:"hash" msg:"hash"`
	}
	GuestVerifyEmailOut struct {
		ResponseCommon
		Ok    bool   `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
		Email string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
	}
)

const (
	GuestVerifyEmailAction = `guest/verifyEmail`

	ErrGuestVerifyEmailInvalidHash        = `invalid hash`
	ErrGuestVerifyEmailUserNotFound       = `user not found`
	ErrGuestVerifyEmailSecretCodeMismatch = `secret code mismatch`
	ErrGuestVerifyEmailModificationFailed = `failed modifying user`
)

func (d *Domain) GuestVerifyEmail(in *GuestVerifyEmailIn) (out GuestVerifyEmailOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	userId, ok := S.DecodeCB63[uint64](in.Hash)
	out.refId = userId

	if !ok {
		out.SetError(400, ErrGuestVerifyEmailInvalidHash)
		return
	}
	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Id = userId
	if !user.FindById() {
		out.SetError(400, ErrGuestVerifyEmailUserNotFound)
		return
	}
	out.actor = userId

	out.Email = user.Email

	if user.VerifiedAt != 0 { // already verified
		out.Ok = true
		return
	}
	if user.SecretCode != in.SecretCode {
		out.SetError(400, ErrGuestVerifyEmailSecretCodeMismatch)
		return
	}
	user.SetVerifiedAt(in.UnixNow())
	user.SetSecretCode(``)
	user.SetSecretCodeAt(0)

	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpdateById() {
		out.SetError(500, ErrGuestVerifyEmailModificationFailed)
		return
	}
	out.Ok = true
	return
}
