package Du

import (
	"github.com/OneOfOne/cmap"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/X"
	"time"
)

// in memory data session (used when no database installed), gone when program exit
type DummySession struct {
	Pool *cmap.CMap
}

type DummyRecord struct {
	Value     interface{}
	ExpiredAt int64
}

func NewDummy() *DummySession {
	return &DummySession{
		Pool: cmap.New(),
	}
}

func (sess *DummySession) Product() string {
	return D.DUMMY
}

func (sess DummySession) Del(key string) {
	sess.Pool.Delete(key)
}

func (sess DummySession) Expiry(key string) int64 {
	val := sess.Pool.Get(key)
	if val == nil {
		return 0
	}
	if rec, ok := val.(DummyRecord); ok {
		if rec.ExpiredAt < 1 {
			return 0
		}
		return rec.ExpiredAt - T.Epoch()
	}
	return 0
}

func (sess DummySession) FadeVal(key string, val interface{}, sec int64) {
	sess.Pool.Set(key, DummyRecord{
		Value:     val,
		ExpiredAt: T.EpochAfter(time.Second * time.Duration(sec)),
	})
}

func (sess DummySession) FadeStr(key, val string, sec int64) {
	sess.Pool.Set(key, DummyRecord{
		Value:     val,
		ExpiredAt: T.EpochAfter(time.Second * time.Duration(sec)),
	})
}

func (sess DummySession) FadeInt(key string, val int64, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess DummySession) FadeMSX(key string, val M.SX, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess DummySession) GetVal(key string) interface{} {
	val := sess.Pool.Get(key)
	if val == nil {
		return nil
	}
	if rec, ok := val.(DummyRecord); ok {
		if rec.ExpiredAt < 1 {
			return nil
		}
		return rec.Value
	}
	return nil
}

func (sess DummySession) GetStr(key string) string {
	return X.ToS(sess.GetVal(key))
}

func (sess DummySession) GetInt(key string) int64 {
	return X.ToI(sess.GetVal(key))
}

func (sess DummySession) GetMSX(key string) M.SX {
	val := sess.GetVal(key)
	if val == nil {
		return M.SX{}
	}
	if m, ok := val.(M.SX); ok {
		return m
	}
	return M.SX{}
}

func (sess DummySession) Inc(key string) int64 {
	val := sess.GetInt(key) + 1
	// TODO: Protect from concurrent access
	sess.SetInt(key, val)
	return val
}

func (sess DummySession) SetStr(key, val string) {
	sess.Pool.Set(key, DummyRecord{Value: val})
}

func (sess DummySession) SetInt(key string, val int64) {
	sess.Pool.Set(key, DummyRecord{Value: val})
}

func (sess DummySession) SetMSX(key string, val M.SX) {
	sess.Pool.Set(key, DummyRecord{Value: val})
}
