package domain

import (
	"github.com/kokizzu/gotro/S"

	"example2/model/mAuth/rqAuth"
	"example2/model/mAuth/wcAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserUpdateProfile.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserUpdateProfile.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserUpdateProfile.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserUpdateProfile.go
//go:generate farify doublequote --file UserUpdateProfile.go

type (
	UserUpdateProfileIn struct {
		RequestCommon
		UserName string `json:"userName" form:"userName" query:"userName" long:"userName" msg:"userName"`
		FullName string `json:"fullName" form:"fullName" query:"fullName" long:"fullName" msg:"fullName"`
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Country  string `json:"country" form:"country" query:"country" long:"country" msg:"country"`
		Language string `json:"language" form:"language" query:"language" long:"language" msg:"language"`
	}

	UserUpdateProfileOut struct {
		ResponseCommon
		User *rqAuth.Users `json:"user" form:"user" query:"user" long:"user" msg:"user"`
	}
)

const (
	UserUpdateProfileAction = `user/updateProfile`

	ErrUpdateProfileUsernameAlreadyUsed = `username already used`
	ErrUpdateProfileEmailAlreadyUsed    = `email already used`
	ErrUpdateProfileFailed              = `update profile failed`
)

func (d *Domain) UserUpdateProfile(in *UserUpdateProfileIn) (out UserProfileOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	sess := d.MustLogin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}
	out.refId = sess.UserId

	user := wcAuth.NewUsersMutator(d.AuthOltp)
	user.Id = sess.UserId
	if !user.FindById() {
		out.SetError(400, ErrUserProfileNotFound)
		return
	}

	if in.Email != `` && user.Email != in.Email {
		dup := rqAuth.NewUsers(d.AuthOltp)
		dup.Email = S.ValidateEmail(in.Email)
		if dup.FindByEmail() && dup.Id != user.Id {
			out.SetError(400, ErrUpdateProfileEmailAlreadyUsed)
			return
		}
		user.SetEmail(dup.Email)
		user.SetVerifiedAt(0) // must also unset verifiedAt
	}

	if in.FullName != `` && user.FullName != in.FullName {
		user.SetFullName(in.FullName)
	}
	if !user.DoUpdateById() {
		user.HaveMutation()
		out.SetError(400, ErrUpdateProfileFailed)
		return
	}

	out.AddDbChangeLogs(user.Logs())

	user.CensorFields()
	out.User = &user.Users
	return
}
