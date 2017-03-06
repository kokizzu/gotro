package Pg

import (
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/T"
	"gitlab.com/kokizzu/gokil/I"
	"gitlab.com/kokizzu/gokil/L"
	"time"
)

var SESSION_DBACTOR_ID = int64(1)
var SESSION_VALUE_KEY = `value`
var SESSION_EXPIRY_KEY = `expired_at`

type PostgreSession struct {
	Pool  *RDBMS
	Table string
}

func NewPostgreSession(conn *RDBMS, table string) *PostgreSession {
	sess := &PostgreSession{
		Pool:  conn,
		Table: table,
	}
	sess.Pool.CreateBaseTable(table)
	return sess
}

func (sess *PostgreSession) Del(key string) {
	sess.Pool.DoTransaction(func(tx *Tx) string {
		dm := tx.QBaseUniq(sess.Table, key)
		if dm.Id < 1 || dm.IsDeleted {
			return ``
		}
		dm.Delete(SESSION_DBACTOR_ID)
		return ``
	})
}

func (sess *PostgreSession) Expiry(key string) int64 {
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

func (sess *PostgreSession) FadeVal(key string, val interface{}, sec int64) {
	sess.Pool.DoTransaction(func(tx *Tx) string {
		dm := tx.QBaseUniq(sess.Table, key)
		dm.SetVal(SESSION_VALUE_KEY, val)
		dm.SetVal(SESSION_EXPIRY_KEY, T.EpochAfter(time.Second*time.Duration(sec)))
		dm.Save(SESSION_DBACTOR_ID)
		return ``
	})
}

func (sess *PostgreSession) FadeStr(key, val string, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess *PostgreSession) FadeInt(key string, val int64, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess *PostgreSession) FadeMSX(key string, val M.SX, sec int64) {
	sess.FadeVal(key, val, sec)
}

func (sess *PostgreSession) GetStr(key string) string {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted || dm.XData == nil {
		return ``
	}
	return dm.GetStr(SESSION_VALUE_KEY)
}

func (sess *PostgreSession) GetInt(key string) int64 {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted || dm.XData == nil {
		return 0
	}
	return dm.GetInt(SESSION_VALUE_KEY)
}

func (sess *PostgreSession) GetMSX(key string) M.SX {
	dm := sess.Pool.QBaseUniq(sess.Table, key)
	if dm.Id < 1 || dm.IsDeleted || dm.XData == nil {
		return M.SX{}
	}
	return dm.GetMSX(SESSION_VALUE_KEY)
}

func (sess *PostgreSession) Inc(key string) (ival int64) {
	k2 := ZZ(SESSION_VALUE_KEY)
	k1 := Z(SESSION_VALUE_KEY)
	table := ZZ(sess.Table)
	uniq := Z(key)
	sess.Pool.DoTransaction(func(tx *Tx) string {
		res := tx.DoExec(`
			INSERT INTO ` + table + `(unique_id,data) VALUES(` + uniq + `,'{` + k2 + `:1}')
			ON CONFLICT(unique_id) DO
			UPDATE SET data=jsonb_merge(` + table + `.data,('{` + k2 + `:' || COALESCE((` + table + `.data->>` + k1 + `)::BIGINT+1,1) || '}')::JSONB)
				, ` + table + `.updated_by=` + I.ToS(SESSION_DBACTOR_ID) + ` 
			WHERE ` + table + `.unique_id = ` + uniq + `
			RETURNING (` + table + `.data->>` + k1 + `)::BIGINT
		`)
		var err error
		ival, err = res.LastInsertId()
		if L.IsError(err, `Inc failed RETURNING`) {
			return err.Error()
		}
		return ``
	})
	return ival
}

func (sess *PostgreSession) SetVal(key string, val interface{}) {
	sess.Pool.DoTransaction(func(tx *Tx) string {
		dm := tx.QBaseUniq(sess.Table, key)
		dm.SetVal(SESSION_VALUE_KEY, val)
		dm.Save(SESSION_DBACTOR_ID)
		return ``
	})
}

func (sess *PostgreSession) SetStr(key, val string) {
	sess.SetVal(key, val)
}

func (sess *PostgreSession) SetInt(key string, val int64) {
	sess.SetVal(key, val)
}

func (sess *PostgreSession) SetMSX(key string, val M.SX) {
	sess.SetVal(key, val)
}
