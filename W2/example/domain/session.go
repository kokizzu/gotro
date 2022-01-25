package domain

import (
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
	"github.com/kpango/fastime"
	"time"
)

func (d *Domain) expireSession(sessionToken string) bool {
	if sessionToken == `` {
		return false
	}
	session := wcAuth.NewSessionsMutator(d.Taran)
	session.SessionToken = sessionToken
	now := fastime.UnixNow()
	if session.FindBySessionToken() {
		if session.ExpiredAt > now {
			session.SetExpiredAt(now)
			session.DoUpdateBySessionToken()
		}
		return true
	}
	return false
}

func (d *Domain) createSession(userId uint64, email, userAgent string) *wcAuth.SessionsMutator {
	session := wcAuth.NewSessionsMutator(d.Taran)
	session.UserId = userId
	sess := conf.Session{
		UserId:    userId,
		Email:     email,
		ExpiredAt: time.Now().AddDate(0, 0, conf.CookieDays).Unix(),
	}
	session.SessionToken = sess.Encrypt(userAgent)
	session.ExpiredAt = sess.ExpiredAt
	return session
}
