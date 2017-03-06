package W

import (
	"github.com/kokizzu/gotro/M"
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
	GetJson(key string) M.SX
	// increment
	Inc(key string) int64
	// set string
	SetStr(key, val string)
	// set integer
	SetInt(key string, val int64)
	// set json
	SetJson(key string, val M.SX)
}

type Session struct {
	Conn      *SessionConnector
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

func InitSession(expire_ns, renew_ns time.Duration) {
	// TODO: continue this

}
