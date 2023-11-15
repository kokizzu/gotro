package domain

import (
	"github.com/kokizzu/gotro/M"

	"example2/model/mAuth/rqAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserProfile.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserProfile.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserProfile.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserProfile.go
//go:generate farify doublequote --file UserProfile.go

type (
	UserProfileIn struct {
		RequestCommon
	}
	UserProfileOut struct {
		ResponseCommon
		User *rqAuth.Users `json:"user" form:"user" query:"user" long:"user" msg:"user"`

		Segments M.SB `json:"segments" form:"segments" query:"segments" long:"segments" msg:"segments"`
	}
)

const (
	UserProfileAction = `user/profile`

	ErrUserProfileNotFound = `user not found`
)

func (d *Domain) UserProfile(in *UserProfileIn) (out UserProfileOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	sess := d.MustLogin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}
	out.refId = sess.UserId

	user := rqAuth.NewUsers(d.AuthOltp)
	user.Id = sess.UserId
	if !user.FindById() {
		out.SetError(403, ErrUserProfileNotFound)
		return
	}
	out.actor = user.Id

	user.CensorFields()
	out.User = user
	out.Segments = sess.Segments
	return
}
