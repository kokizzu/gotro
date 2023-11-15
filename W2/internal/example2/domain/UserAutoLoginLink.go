package domain

import (
	"fmt"

	"github.com/kokizzu/gotro/S"

	"example2/conf"
	"example2/model/mAuth/rqAuth"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserAutoLoginLink.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserAutoLoginLink.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserAutoLoginLink.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserAutoLoginLink.go
//go:generate farify doublequote --file UserAutoLoginLink.go

type (
	UserAutoLoginLinkIn struct {
		RequestCommon

		Path string `json:"path" form:"path" query:"path" long:"path" msg:"path"`
	}
	UserAutoLoginLinkOut struct {
		ResponseCommon

		Link string `json:"link" form:"link" query:"link" long:"link" msg:"link"`
	}
)

const (
	UserAutoLoginLinkAction = `user/autoLoginLink`

	ErrUserAutoLoginLinkInvalidFor = `autologin link invalid path`
)

func (d *Domain) UserAutoLoginLink(in *UserAutoLoginLinkIn) (out UserAutoLoginLinkOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	sess := d.MustLogin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}

	if !S.StartsWith(in.Path, `/`) {
		out.SetError(400, ErrUserAutoLoginLinkInvalidFor)
		return
	}

	user := rqAuth.NewUsers(d.AuthOltp)
	user.Id = sess.UserId
	user.FindById() // assume always exists

	out.Link = `/` + GuestAutoLoginAction + `?uid=` + S.EncodeCB63(sess.UserId, 1) +
		`&token=` + sess.Encrypt(conf.AutoLoginUA+fmt.Sprint(user.UpdatedAt)+in.Path) +
		`&path=` + in.Path
	return
}
