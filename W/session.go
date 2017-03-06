package W

import (
	"github.com/kokizzu/gotro/M"
	"gitlab.com/kokizzu/gokil/S"
	"time"
)

type SessionConnector interface {
	// delete key
	Del(key string)
	// get remaining lifespan in seconds
	Expiry(key string) int64
	// set string with remaining lifespan in seconds
	FadeStr(key, val string, ttl int64)
	// set integer with remaining lifespan in seconds
	FadeInt(key string, val int64, ttl int64)
	// set json with remaining lifespan in seconds
	FadeMSX(key string, val M.SX, ttl int64)
	// get string
	GetStr(key string) string
	// get integer
	GetInt(key string) int64
	// get string
	GetMSX(key string) M.SX
	// increment
	Inc(key string) int64
	// set string
	SetStr(key, val string)
	// set integer
	SetInt(key string, val int64)
	// set json
	SetMSX(key string, val M.SX)
}

var SESS_KEY = `SK`
var EXPIRE_NS time.Duration
var RENEW_NS time.Duration

type Session struct {
	UserId    int64
	AppId     int64
	Level     M.SX
	Email     string
	UserAgent string
	IpAddr    string
	// TODO: continue this
}

func (s *Session) Logout() {
	// TODO: continue this
}

func (s *Session) Login(id int64, email string) {
	// TODO: continue this
}

func (s *Session) StateCSRF() string {
	// TODO: continue this
	return ``
}

func (s *Session) Touch() {
	// TODO: update ttl
}

func InitSession(sess_key string, expire_ns, renew_ns time.Duration, conn SessionConnector) {
	SESS_KEY = S.IfEmpty(sess_key, SESS_KEY)
	EXPIRE_NS = expire_ns
	RENEW_NS = renew_ns
	Sessions = conn
}
