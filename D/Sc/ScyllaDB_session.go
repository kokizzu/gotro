package Sc

import "github.com/kokizzu/gotro/M"

type ScyllaSession struct {
}

func (sess ScyllaSession) Expiry(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess ScyllaSession) FadeStr(key, val string, sec int64) {
	// TODO: continue this
}

func (sess ScyllaSession) FadeInt(key string, val int64, sec int64) {
	// TODO: continue this
}

func (sess ScyllaSession) FadeMSX(key string, val M.SX, sec int64) {
	// TODO: continue this
}

func (sess ScyllaSession) GetStr(key string) string {
	// TODO: continue this
	return ``
}

func (sess ScyllaSession) GetInt(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess ScyllaSession) GetMSX(key string) M.SX {
	// TODO: continue this
	res := M.SX{}
	return res
}

func (sess ScyllaSession) Inc(key string) int64 {
	// TODO: continue this
	return 0
}

func (sess ScyllaSession) SetStr(key, val string) {
	// TODO: continue this
}

func (sess ScyllaSession) SetInt(key string, val int64) {
	// TODO: continue this
}

func (sess ScyllaSession) SetMSX(key string, val M.SX) {
	// TODO: continue this
}

func (sess ScyllaSession) Del(key string) {
	// TODO: continue this
}
