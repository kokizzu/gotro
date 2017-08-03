package As

import (
	"github.com/OneOfOne/cmap"
	aero "github.com/aerospike/aerospike-client-go"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

const DEFAULT_HOST = `127.0.0.1`
const DEFAULT_PORT = 3333

var POLICIES *cmap.CMap

func init() {
	POLICIES = cmap.New()
}

func PolicyByTtl(sec int64) *aero.WritePolicy {
	if sec <= 0 {
		return nil
	}
	idx := I.ToS(sec)
	if any := POLICIES.Get(idx); any != nil {
		wp, ok := any.(*aero.WritePolicy)
		if ok {
			return wp
		}
	}
	wp := aero.NewWritePolicy(0, uint32(sec))
	POLICIES.Set(idx, wp)
	return wp
}

type AerosSession struct {
	Pool      *aero.Client
	Namespace string
	Bucket    string
}

func NewAerosSession(host string, port int, namespace, bucket string) *AerosSession {
	host = S.IfEmpty(host, DEFAULT_HOST)
	port = I.IsZero(port, DEFAULT_PORT)
	conn, err := aero.NewClient(host, port)
	L.PanicIf(err, `Failed to connect to in-memory database`)
	return &AerosSession{
		Pool:      conn,
		Namespace: namespace,
		Bucket:    bucket,
	}
}
func (sess AerosSession) Product() string {
	return D.AEROSP
}

func (sess AerosSession) Key(key string) *aero.Key {
	asKey, err := aero.NewKey(sess.Namespace, sess.Bucket, key)
	if L.IsError(err, `Failed to create compute digest key `+sess.Bucket+` `+key) {
		return nil
	}
	return asKey
}

func (sess AerosSession) Del(key string) {
	if askey := sess.Key(key); askey != nil {
		_, err := sess.Pool.Delete(nil, askey)
		L.IsError(err, `Error deleting CACHE `+askey.String())
	}
}

func (sess AerosSession) Expiry(key string) int64 {
	if askey := sess.Key(key); askey != nil {
		rec, err := sess.Pool.GetHeader(nil, askey)
		if err == nil && rec != nil {
			return int64(rec.Expiration)
		}
	}
	return 0
}

func (sess AerosSession) FadeStr(key, val string, sec int64) {
	if asKey := sess.Key(key); asKey != nil {
		policy := PolicyByTtl(sec)
		err := sess.Pool.Put(policy, asKey, aero.BinMap{
			`value`: val,
		})
		L.IsError(err, `Error putting CACHE `+asKey.String())
	}
}

func (sess AerosSession) FadeInt(key string, val int64, sec int64) {
	sess.FadeStr(key, I.ToS(val), sec)
}

func (sess AerosSession) FadeMSX(key string, val M.SX, sec int64) {
	sess.FadeStr(key, M.ToJson(val), sec)
}

func (sess AerosSession) GetStr(key string) string {
	if asKey := sess.Key(key); asKey != nil {
		rec, err := sess.Pool.Get(nil, asKey)
		if !L.IsError(err, `Error getting CACHE `+asKey.String()) {
			return X.ToS(rec.Bins[`value`])
		}
	}
	return ``
}

func (sess AerosSession) GetInt(key string) int64 {
	val := sess.GetStr(key)
	if val == `` {
		return 0
	}
	return S.ToI(val)
}

func (sess AerosSession) GetMSX(key string) M.SX {
	val := sess.GetStr(key)
	if val == `` {
		return M.SX{}
	}
	return S.JsonToMap(val)
}

func (sess AerosSession) Inc(key string) int64 {
	if askey := sess.Key(key); askey != nil {
		ops := []*aero.Operation{
			aero.AddOp(aero.NewBin(`value`, 1)), // add the value of the bin to the existing value
			aero.GetOp(),                        // get the value of the record after all operations are executed
		}
		rec, err := sess.Pool.Operate(nil, askey, ops...)
		if !L.IsError(err, `Error getting CACHE `+askey.String()) {
			return X.ToI(rec.Bins[`value`])
		}
	}
	return 0
}

func (sess AerosSession) SetStr(key, val string) {
	if asKey := sess.Key(key); asKey != nil {
		err := sess.Pool.Put(nil, asKey, aero.BinMap{
			`value`: val,
		})
		L.IsError(err, `Error putting CACHE `+asKey.String())
	}
}

func (sess AerosSession) SetInt(key string, val int64) {
	sess.SetStr(key, I.ToS(val))
}

func (sess AerosSession) SetMSX(key string, val M.SX) {
	sess.SetStr(key, M.ToJson(val))
}
