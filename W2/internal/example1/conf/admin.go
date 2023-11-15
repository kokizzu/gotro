package conf

import "github.com/kpango/fastime"

const SuperAdmin = `admin@localhost`

var Admins = map[string]bool{
	SuperAdmin: true,
}

var AdminTestSessionToken string
var AdminTestExpiredSession string

func init() {
	sess := Session{
		UserId:    1,
		Email:     SuperAdmin,
		ExpiredAt: fastime.UnixNow() + 60*60*24*365,
	}
	AdminTestSessionToken = sess.Encrypt(``)
	sess.ExpiredAt = fastime.UnixNow() - 60*60*24*365
	AdminTestExpiredSession = sess.Encrypt(``)
}
