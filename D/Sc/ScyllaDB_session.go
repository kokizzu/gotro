package Sc

import (
	"github.com/gocql/gocql"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"time"
)

const KV_TABLE = `kvs`
const COUNTER_TABLE = `counters`

type KeyValue struct {
	Type      string
	Key       string
	Value     string
	ExpiredAt string
}

type ScyllaSession struct {
	Name    string
	Prefix  string
	Cluster *gocql.ClusterConfig
	Session *gocql.Session
}

// CREATE KEYSPACE "replace_this" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
// note that INSERT and UPDATE in C*/Scylla is an UPSERT
func NewScyllaSession(ip, keyspace, prefix, user, pass string) *ScyllaSession {
	// L.Print(`NewScyllaSession`)
	clust := gocql.NewCluster(ip)
	clust.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	clust.Timeout = 8 * time.Second
	clust.ConnectTimeout = 15 * time.Second
	clust.Keyspace = keyspace
	if user != `` {
		clust.Authenticator = gocql.PasswordAuthenticator{
			Username: user,
			Password: pass,
		}
	}
	sess, err := clust.CreateSession()
	L.PanicIf(err, `Failed create session Sc`)
	ksm, err := sess.KeyspaceMetadata(keyspace)
	L.PanicIf(err, `Failed retrieve meta`)
	if ksm.Tables[KV_TABLE] == nil {
		err = sess.Query(`CREATE TABLE ` + KV_TABLE + ` ("type" TEXT, "key" TEXT, "value" TEXT, "expired_at" BIGINT, PRIMARY KEY("type","key"))`).Exec()
		L.PanicIf(err, `Failed create kv table`)
	}
	if ksm.Tables[COUNTER_TABLE] == nil {
		err = sess.Query(`CREATE TABLE ` + COUNTER_TABLE + ` ("key" TEXT PRIMARY KEY, "value" COUNTER)`).Exec()
		L.PanicIf(err, `Failed create counter table`)
	}
	return &ScyllaSession{
		Name:    `sc://` + user + `:` + pass + `@` + ip + `/` + keyspace,
		Prefix:  prefix,
		Cluster: clust,
		Session: sess,
	}
}

func (sess ScyllaSession) Product() string {
	return D.SCYLLA
}

func (sess ScyllaSession) Expiry(key string) int64 {
	// L.Print(`Expiry`)
	res := M.SX{}
	cql := `SELECT "expired_at" FROM ` + KV_TABLE + sess.Where(key)
	iter := sess.Session.Query(cql).Iter()
	defer iter.Close()
	iter.MapScan(res)
	expired_at := res.GetInt(`expired_at`)
	if expired_at < 1 {
		return 0
	}
	return expired_at - T.Epoch()
}

func (sess ScyllaSession) fadeStr(key, val string, sec int64) {
	cql := `UPDATE ` + KV_TABLE + ` USING TTL ` + I.ToS(sec) + ` SET "value" = ` + val + `, "expired_at" = ` + T.EpochAfterStr(time.Second*time.Duration(sec)) + sess.Where(key)
	err := sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
}

func (sess ScyllaSession) FadeStr(key, val string, sec int64) {
	// L.Print(`FadeStr`)
	sess.fadeStr(key, Z(val), sec)
}

func (sess ScyllaSession) FadeInt(key string, val int64, sec int64) {
	// L.Print(`FadeInt`)
	sess.fadeStr(key, I.ToS(val), sec)
}

func (sess ScyllaSession) FadeMSX(key string, val M.SX, sec int64) {
	// L.Print(`FadeMSX`)
	sess.fadeStr(key, S.ZJ(val.ToJson()), sec)
}

func (sess ScyllaSession) GetStr(key string) string {
	// L.Print(`GetStr`)
	res := M.SX{}
	cql := `SELECT "value" FROM ` + KV_TABLE + sess.Where(key)
	iter := sess.Session.Query(cql).Iter()
	defer iter.Close()
	iter.MapScan(res)
	return res.GetStr(`value`)
}

func (sess ScyllaSession) GetInt(key string) int64 {
	// L.Print(`GetStr`)
	return S.ToI(sess.GetStr(key))
}

func (sess ScyllaSession) GetMSX(key string) M.SX {
	// L.Print(`GetMSX`)
	return S.JsonToMap(sess.GetStr(key))
}

func (sess ScyllaSession) Inc(key string) (ret int64) {
	// workaround, since no data type conversion supported: http://stackoverflow.com/questions/44109822/how-to-convert-cassandra-scylladb-text-to-bigint-and-vice-versa
	cql := `UPDATE ` + COUNTER_TABLE + ` SET value = value+1 WHERE "key" = ` + Z(key)
	err := sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
	cql = `SELECT "value" FROM ` + COUNTER_TABLE + `WHERE "key" = ` + Z(key)
	iter := sess.Session.Query(cql).Iter()
	defer iter.Close()
	iter.Scan(&ret)
	return
}

func (sess ScyllaSession) Where(key string) string {
	return ` WHERE "type" = ` + Z(sess.Prefix) + ` AND "key" = ` + Z(key)
}

func (sess ScyllaSession) setStr(key, val string) {
	cql := `UPDATE ` + KV_TABLE + ` SET "value" = ` + val + sess.Where(key)
	err := sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
}

func (sess ScyllaSession) SetStr(key, val string) {
	// L.Print(`SetStr`)
	sess.setStr(key, Z(val))
}

func (sess ScyllaSession) SetInt(key string, val int64) {
	// L.Print(`SetInt`)
	sess.setStr(key, Z(I.ToS(val)))
}

func (sess ScyllaSession) SetMSX(key string, val M.SX) {
	// L.Print(`SetMSX`)
	sess.setStr(key, S.ZJ(val.ToJson()))
}

func (sess ScyllaSession) SetMSS(key string, val M.SS) {
	// L.Print(`SetMSX`)
	sess.setStr(key, val.ToScylla())
}

func (sess ScyllaSession) Del(key string) {
	// L.Print(`Del`)
	now := T.EpochStr()
	cql := `UPDATE ` + KV_TABLE + ` USING TTL 0 SET "expired_at" = ` + now + sess.Where(key)
	err := sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
}
