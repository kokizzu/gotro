package Rd

import (
	"context"
	"strings"

	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"

	"github.com/rueian/rueidis"
)

const DEFAULT_HOST = `127.0.0.1:6379`

type RedisSession struct {
	Pool   rueidis.Client
	Prefix string
}

func (sess RedisSession) Product() string {
	return D.REDIS
}

// TryRedisSession non panic version, returns error if failed to connect
func TryRedisSession(host, pass string, dbNum int, prefix string) (*RedisSession, error) {
	host = S.IfEmpty(host, DEFAULT_HOST)
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{host},
		Password:    pass,
		SelectDB:    dbNum,
	})
	return &RedisSession{
		Pool:   conn,
		Prefix: prefix,
	}, err
}

// NewRedisSession panic version
func NewRedisSession(host, pass string, dbNum int, prefix string) *RedisSession {
	host = S.IfEmpty(host, DEFAULT_HOST)
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{host},
		Password:    pass,
		SelectDB:    dbNum,
	})
	L.PanicIf(err, `redis.NewClient`)
	return &RedisSession{
		Pool:   conn,
		Prefix: prefix,
	}
}

func (sess RedisSession) Del(key string) {

	err := sess.Pool.Do(context.Background(), sess.Pool.B().Del().Key(sess.Prefix+key).Build()).Error()
	L.IsError(err, `failed to DEL`, key)
}

//Expiry check the expiry time in second
func (sess RedisSession) Expiry(key string) int64 {

	val, err := sess.Pool.Do(context.Background(), sess.Pool.B().Ttl().Key(sess.Prefix+key).Build()).AsInt64()
	if err != nil {
		return -1
	}
	return val
}

func (sess RedisSession) FadeStr(key, val string, sec int64) {
	err := sess.Pool.Do(context.Background(), sess.Pool.B().Setex().Key(sess.Prefix+key).Seconds(sec).Value((val)).Build()).Error()
	L.IsError(err, `failed to SETEX`, key, sec, val)
}

func (sess RedisSession) FadeInt(key string, val int64, sec int64) {
	sess.FadeStr(key, I.ToS(val), sec)
}

func (sess RedisSession) FadeMSX(key string, val M.SX, sec int64) {
	sess.FadeStr(key, M.ToJson(val), sec)
}

func (sess RedisSession) GetStr(key string) string {
	val, err := sess.Pool.Do(context.Background(), sess.Pool.B().Get().Key(sess.Prefix+key).Build()).ToString()
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

	val, err := sess.Pool.Do(context.Background(), sess.Pool.B().Incr().Key(sess.Prefix+key).Build()).AsInt64()
	L.IsError(err, `failed to INCR`, key)
	return val
}

func (sess RedisSession) Dec(key string) int64 {
	val, err := sess.Pool.Do(context.Background(), sess.Pool.B().Decr().Key(sess.Prefix+key).Build()).AsInt64()
	L.IsError(err, `failed to DECR`, key)
	return val
}

func (sess RedisSession) SetStr(key, val string) {
	err := sess.Pool.Do(context.Background(), sess.Pool.B().Set().Key(sess.Prefix+key).Value(val).Build()).Error()
	L.IsError(err, `failed to SET`, key, val)
}

func (sess RedisSession) SetInt(key string, val int64) {
	sess.SetStr(key, I.ToS(val))
}
func (sess RedisSession) SetMSX(key string, val M.SX) {
	sess.SetStr(key, val.ToJson())
}
func (sess RedisSession) SetMSS(key string, val M.SS) {
	sess.SetStr(key, val.ToJson())
}

func (sess RedisSession) Lpush(key string, val string) {
	sess.Pool.Do(context.Background(), sess.Pool.B().Lpush().Key(key).Element(val).Build())
}
func (sess RedisSession) Rpush(key string, val string) {
	sess.Pool.Do(context.Background(), sess.Pool.B().Rpush().Key(key).Element(val).Build())
}
func (sess RedisSession) Lrange(key string, start, end int64) []string {
	res, err := sess.Pool.Do(context.Background(), sess.Pool.B().Lrange().Key(key).Start(start).Stop(end).Build()).ToArray()
	if err != nil {
		return []string{}
	}
	str := make([]string, 0)
	for _, v := range res {
		v, _ := v.ToString()
		str = append(str, v)
	}
	return str
}
