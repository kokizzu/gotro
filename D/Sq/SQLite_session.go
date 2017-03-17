package Sq

import (
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/T"
	"time"
)

var SESSION_DBACTOR_ID = int64(1)
var SESSION_VALUE_KEY = `value`
var SESSION_EXPIRY_KEY = `expired_at`

type SqliteSession struct {
	Pool  *RDBMS
	Table string
}

func NewSession(conn *RDBMS, table string) *SqliteSession {
	sess := &SqliteSession{
		Pool:  conn,
		Table: table,
	}
	InitFunctions(sess.Pool)
	sess.Pool.CreateBaseTable(table)
	return sess
}

func (sess SqliteSession) Del(key string) {
	sess.Pool.DoTransaction(func(tx *Tx) string {
		dm := tx.QBaseUniq(sess.Table, key)
		if dm.Id < 1 || dm.IsDeleted {
			return ``
		}
		dm.Delete(SESSION_DBACTOR_ID)
		return ``
	})
}

func (sess SqliteSession) Expiry(key string) int64 {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted {
		return 0
	}
	expired_at := dm.GetInt(SESSION_EXPIRY_KEY)
	if expired_at < 1 {
		return 0
	}
	return expired_at - T.Epoch()
}

func (sess SqliteSession) FadeVal(key string, val interface{}, sec int64) {
	sess.Pool.DoTransaction(func(tx *Tx) string {
		dm := tx.QBaseUniq(sess.Table, key)
		dm.SetVal(SESSION_VALUE_KEY, val)
		dm.SetVal(SESSION_EXPIRY_KEY, T.EpochAfter(time.Second*time.Duration(sec)))
		dm.Save(SESSION_DBACTOR_ID)
		return ``
	})
}
func (sess SqliteSession) FadeStr(key, val string, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess SqliteSession) FadeInt(key string, val int64, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess SqliteSession) FadeMSX(key string, val M.SX, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess SqliteSession) GetStr(key string) string {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted || dm.XData == nil {
		return ``
	}
	return dm.GetStr(SESSION_VALUE_KEY)
}

func (sess SqliteSession) GetInt(key string) int64 {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted || dm.XData == nil {
		return 0
	}
	return dm.GetInt(SESSION_VALUE_KEY)
}

func (sess SqliteSession) GetMSX(key string) M.SX {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted || dm.XData == nil {
		return M.SX{}
	}
	return dm.GetMSX(SESSION_VALUE_KEY)
}

func (sess SqliteSession) Inc(key string) (ival int64) {
	/* Couldn't use SQLite's JSON functions
	k2 := ZZ(SESSION_VALUE_KEY)
	// k1 := Z(SESSION_VALUE_KEY)
	table := sess.Table
	uniq := Z(key)
	sess.Pool.DoTransaction(func(tx *Tx) string {
		tx.DoExec(`
			INSERT OR IGNORE INTO ` + table + `(unique_id, data) VALUES(` + uniq + `,'{` + k2 + `:0}')
		`)
		res := tx.DoExec(`
			UPDATE ` + table + `
			SET data = (SELECT json_set(json(` + table + `.data), '$.` + k2 + `' , 12) FROM ` + table + `)
			WHERE unique_id =` + uniq + `
		`)
		res := tx.DoExec(`
			SELECT JSON_EXTRACT(` + table + `.data, '$.` + k2 + `') FROM ` + table + `
		`)
		var err error
		ival, err = res.LastInsertId()
		if L.IsError(err, `Inc failed RETURNING`) {
			return err.Error()
		}
		return ``
	})
	*/
	return ival
}

func (sess SqliteSession) SetVal(key string, val interface{}) {
	sess.Pool.DoTransaction(func(tx *Tx) string {
		dm := tx.QBaseUniq(sess.Table, key)
		dm.SetVal(SESSION_VALUE_KEY, val)
		dm.Save(SESSION_DBACTOR_ID)
		return ``
	})
}
func (sess SqliteSession) SetStr(key, val string) {
	sess.SetVal(key, val)
}

func (sess SqliteSession) SetInt(key string, val int64) {
	sess.SetVal(key, val)
}

func (sess SqliteSession) SetMSX(key string, val M.SX) {
	sess.SetVal(key, val)
}
