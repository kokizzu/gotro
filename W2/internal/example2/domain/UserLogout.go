package domain

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file UserLogout.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type UserLogout.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type UserLogout.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type UserLogout.go
//go:generate farify doublequote --file UserLogout.go

type (
	UserLogoutIn struct {
		RequestCommon
	}
	UserLogoutOut struct {
		ResponseCommon
		LogoutAt int64 `json:"loggedOut" form:"loggedOut" query:"loggedOut" long:"loggedOut" msg:"loggedOut"`
	}
)

const (
	UserLogoutAction = `user/logout`

	ErrUserSessionRemovalFailed = `user session removal failed`
)

func (d *Domain) UserLogout(in *UserLogoutIn) (out UserLogoutOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	out.LogoutAt = d.ExpireSession(in.SessionToken, &out.ResponseCommon)
	return
}
