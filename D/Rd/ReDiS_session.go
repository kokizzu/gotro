package Rd

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"gopkg.in/redis.v5"
	"strings"
	"time"
)

const DEFAULT_HOST = `127.0.0.1:6379`

type RedisSession struct {
	Pool   *redis.Client
	Prefix string
}

func NewRedisSession(host, pass string, db_num int, prefix string) *RedisSession {
	host = S.IfEmpty(host, DEFAULT_HOST)
	conn := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       db_num,
	})
	return &RedisSession{
		Pool:   conn,
		Prefix: prefix,
	}
}

func (sess RedisSession) Del(key string) {
	err := sess.Pool.Del(sess.Prefix + key).Err()
	L.IsError(err, `failed to DEL`, key)
}

func (sess RedisSession) Expiry(key string) int64 {
	val := sess.Pool.TTL(sess.Prefix + key).Val()
	return int64(val.Seconds())
}

func (sess RedisSession) FadeStr(key, val string, sec int64) {
	err := sess.Pool.Set(sess.Prefix+key, val, time.Second*time.Duration(sec)).Err()
	L.IsError(err, `failed to SETEX`, key, sec, val)
}

func (sess RedisSession) FadeInt(key string, val int64, sec int64) {
	sess.FadeStr(key, I.ToS(val), sec)
}

func (sess RedisSession) FadeMSX(key string, val M.SX, sec int64) {
	sess.FadeStr(key, M.ToJson(val), sec)
}

func (sess RedisSession) GetStr(key string) string {
	val, err := sess.Pool.Get(sess.Prefix + key).Result()
	if err != nil && err.Error() != `redis: nil` {
		L.IsError(err, `failed to GET`, key)
	}
	return strings.Trim(val, `"`) // remove quotes, khusus redis
}

func (sess RedisSession) GetInt(key string) int64 {
	return S.ToI(sess.GetStr(key))
}

func (sess RedisSession) GetMSX(key string) M.SX {
	return S.JsonToMap(sess.GetStr(key))
}

func (sess RedisSession) Inc(key string) int64 {
	val, err := sess.Pool.Incr(sess.Prefix + key).Result()
	L.IsError(err, `failed to INCR`, key)
	return val
}

func (sess RedisSession) SetStr(key, val string) {
	err := sess.Pool.Set(sess.Prefix+key, val, 0).Err()
	L.IsError(err, `failed to SET`, key, val)
}

func (sess RedisSession) SetInt(key string, val int64) {
	sess.SetStr(key, I.ToS(val))
}

func (sess RedisSession) SetMSX(key string, val M.SX) {
	sess.SetStr(key, M.ToJson(val))
}
