package Sc

import (
	"github.com/gocql/gocql"
	"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"time"
)

const TABLE = `kvs`

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

func intStrOf(val string) string {
	int_val := `0`
	if val != `` && C.IsDigit(val[0]) {
		int_val = I.ToS(S.ToI(val))
	}
	return int_val
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
	if ksm.Tables[TABLE] == nil {
		err = sess.Query(`CREATE TABLE ` + TABLE + ` ("type" TEXT, "key" TEXT, "inv_value" BIGINT, "value" TEXT, "expired_at" BIGINT, PRIMARY KEY("type","key"))`).Exec()
		L.PanicIf(err, `Failed create session table`)
	}
	return &ScyllaSession{
		Name:    `sc://` + user + `:` + pass + `@` + ip + `/` + keyspace,
		Prefix:  prefix,
		Cluster: clust,
		Session: sess,
	}
}

func (sess ScyllaSession) Expiry(key string) int64 {
	// L.Print(`Expiry`)
	res := M.SX{}
	cql := `SELECT "expired_at" FROM ` + TABLE + sess.Where(key)
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
	int_val := intStrOf(val)
	cql := `UPDATE ` + TABLE + ` USING TTL ` + I.ToS(sec) + ` SET "value" = ` + val + `, "int_value" = ` + int_val + `, "expired_at" = ` + T.EpochAfterStr(time.Second*time.Duration(sec)) + sess.Where(key)
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
	cql := `SELECT "int_value","value" FROM ` + TABLE + sess.Where(key)
	iter := sess.Session.Query(cql).Iter()
	defer iter.Close()
	iter.MapScan(res)
	int_val := res.GetStr(`int_value`)
	val := res.GetStr(`value`)
	if int_val == `` || int_val == `0` {
		return val
	}
	return int_val
}

func (sess ScyllaSession) GetInt(key string) int64 {
	// L.Print(`GetStr`)
	return S.ToI(sess.GetStr(key))
}

func (sess ScyllaSession) GetMSX(key string) M.SX {
	// L.Print(`GetMSX`)
	return S.JsonToMap(sess.GetStr(key))
}

func (sess ScyllaSession) Inc(key string) int64 {
	// workaround, since no data type conversion supported: http://stackoverflow.com/questions/44109822/how-to-convert-cassandra-scylladb-text-to-bigint-and-vice-versa
	cql := `UPDATE ` + TABLE + ` SET int_value = int_value+1 ` + sess.Where(key)
	err := sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
	int_val := sess.GetStr(key)
	cql = `UPDATE ` + TABLE + ` SET value = ` + int_val + ` ` + sess.Where(key) + ` IF int_value = ` + int_val
	err = sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
	return 0
}

func (sess ScyllaSession) Where(key string) string {
	return ` WHERE "type" = ` + Z(sess.Prefix) + ` AND "key" = ` + Z(key)
}

func (sess ScyllaSession) setStr(key, val string) {
	int_val := intStrOf(val)
	cql := `UPDATE ` + TABLE + ` SET "value" = ` + val + `, "int_value" = ` + int_val + sess.Where(key)
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

func (sess ScyllaSession) Del(key string) {
	// L.Print(`Del`)
	now := T.EpochStr()
	cql := `UPDATE ` + TABLE + ` USING TTL 0 SET "expired_at" = ` + now + sess.Where(key)
	err := sess.Session.Query(cql).Exec()
	L.IsError(err, cql)
}
