package domain

import (
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserChangePassword.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserChangePassword.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserChangePassword.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserChangePassword.go
//go:generate farify doublequote --file UserChangePassword.go

type (
	UserChangePasswordIn struct {
		RequestCommon
		OldPass string `json:"oldPass" form:"oldPass" query:"oldPass" long:"oldPass" msg:"oldPass"`
		NewPass string `json:"newPass" form:"newPass" query:"newPass" long:"newPass" msg:"newPass"`
	}
	UserChangePasswordOut struct {
		ResponseCommon
		ok bool
	}
)

const (
	UserChangePasswordAction = `user/changePassword`

	ErrUserChangePasswordUserNotFound    = `user to change password not found`
	ErrUserChangePasswordWrongOldPass    = `old password does not match`
	ErrUserChangePasswordNewPassTooShort = `new password too short`
	ErrUserChangePasswordSaveUserFailed  = `failed saving user`
)

func (d *Domain) UserChangePassword(in *UserChangePasswordIn) (out UserChangePasswordOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	sess := d.MustLogin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}
	out.refId = sess.UserId

	if len(in.NewPass) < minPassLength {
		out.SetError(400, ErrUserChangePasswordNewPassTooShort)
		return
	}

	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Id = sess.UserId
	if !user.FindById() {
		out.SetError(400, ErrUserChangePasswordUserNotFound)
		return
	}

	if err := user.CheckPassword(in.OldPass); err != nil {
		out.SetError(400, ErrUserChangePasswordWrongOldPass)
		return
	}

	user.SetEncryptedPassword(in.NewPass, in.UnixNow())
	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpdateById() {
		out.SetError(500, ErrUserChangePasswordSaveUserFailed)
		return
	}

	out.ok = true
	return
}
