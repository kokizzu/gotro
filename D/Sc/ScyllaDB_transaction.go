package Sc

import (
	"github.com/gocql/gocql"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
	"time"
)

type Tx struct {
	Trans *gocql.Session
}

// query 2 colums of integer-integer as map
func (tx *Tx) QIntIntMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := int64(0)
	val := key
	for rows.Scan(&key, &val) {
		res[key] = val
	}
	return res
}

// query 2 colums of integer-string as map
func (tx *Tx) QIntStrMap(query string, params ...interface{}) M.IS {
	res := M.IS{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := int64(0)
	val := ``
	for rows.Scan(&key, &val) {
		res[key] = val
	}
	return res
}

// query 1 colums of string as map
// SELECT unique_id
func (tx *Tx) QStrBoolMap(query string, params ...interface{}) M.SB {
	res := M.SB{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := ``
	for rows.Scan(&key) {
		res[key] = true
	}
	return res
}

// query single column int64, return with true value
func (tx *Tx) QIntBoolMap(query string, params ...interface{}) M.IB {
	res := M.IB{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := int64(0)
	for rows.Scan(&key) {
		res[key] = true
	}
	return res
}

// query 2 colums of string-string as map
func (tx *Tx) QStrStrMap(query string, params ...interface{}) M.SS {
	res := M.SS{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := ``
	val := ``
	for rows.Scan(&key, &val) {
		res[key] = val
	}
	return res
}

//// query 1+N colums of string-[]any as map
//func (tx *Tx) QStrArrMap(query string, params ...interface{}) M.SAX {
//	// not implemented, far too inefficient with current gocql
//}

// query 2 colums of string-integer as map
func (tx *Tx) QStrIntMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := ``
	val := int64(0)
	for rows.Scan(&key, &val) {
		res[key] = val
	}
	return res
}

// query 1 colums of integer
func (tx *Tx) QIntCountMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	val := int64(0)
	for rows.Scan(&val) {
		res[val] += 1
	}
	return res
}

// query 1 colums of string
func (tx *Tx) QStrCountMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	val := ``
	for rows.Scan(&val) {
		res[val] += 1
	}
	return res
}

// query 1 colums of integer
func (tx *Tx) QIntArr(query string, params ...interface{}) []int64 {
	res := []int64{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	val := int64(0)
	for rows.Scan(&val) {
		res = append(res, val)
	}
	return res
}

// query one cell [1,2,3,...] and return array of int64
func (tx *Tx) QJsonIntArr(query string, params ...interface{}) []int64 {
	str := tx.QStr(query, params...)
	ax := S.JsonToArr(str)
	ai := []int64{}
	for _, v := range ax {
		ai = append(ai, X.ToI(v))
	}
	return ai
}

// query 1 colums of string
func (tx *Tx) QStrArr(query string, params ...interface{}) []string {
	res := []string{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	val := ``
	for rows.Scan(&val) {
		res = append(res, val)
	}
	return res
}

// do query all, also calls DoPrepare, don't forget to close the rows
func (tx *Tx) QAll(query string, params ...interface{}) (rows Records) {
	var start time.Time
	if DEBUG {
		start = time.Now()
	}
	ite := tx.Trans.Query(query, params...).Iter()
	rows = Records{ite, query, params}
	if DEBUG {
		L.TimeTrack(start, query)
	}
	return
}

// execute a select single value query, convert to string
func (tx *Tx) QStr(query string, params ...interface{}) (dest string) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	L.CheckIf(rows.Scan(&dest), `failed to Tx.QStr: %s %# v`, query, params)
	return
}

// execute a select single value query, convert to bool
func (tx *Tx) QBool(query string, params ...interface{}) (dest bool) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	L.CheckIf(rows.Scan(&dest), `failed to Tx.QBool: %s %# v`, query, params)
	return
}

//// check non-zero if unique exists
//func (tx *Tx) QId(table, id string) (id int64) {
//	// not implemented: no unique constraint on scylladb
//}

// check if id exists
func (tx *Tx) QExists(table string, id string) bool {
	query := `SELECT COUNT(*) FROM ` + table + ` WHERE id = ` + Z(id)
	count := tx.QInt(query)
	return count == 1
}

//// return non-empty string if id exists
//func (tx *Tx) QUniq(table string, id string) (uniq string) {
//	// not implemented: useless id is unique_id, and there are no additional unique constraint on scylladb
//}

// execute a select pair value query, convert to int64 and string
func (tx *Tx) QIntStr(query string, params ...interface{}) (i int64, s string) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	rows.Scan(&i, &s)
	return
}

// execute a select pair value query, convert to string and int
func (tx *Tx) QStrInt(query string, params ...interface{}) (s string, i int64) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	rows.Scan(&s, &i)
	return
}

// execute a select pair value query, convert to string and string
func (tx *Tx) QStrStr(query string, params ...interface{}) (s string, ss string) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	rows.Scan(&s, &ss)
	return
}

// execute a select single value query, convert to int64
func (tx *Tx) QInt(query string, params ...interface{}) (dest int64) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	L.CheckIf(rows.Scan(&dest), `failed to Tx.QInt: %s %# v`, query, params)
	return
}

func (tx *Tx) QFloat(query string, params ...interface{}) (dest float64) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	L.CheckIf(rows.Scan(&dest), `failed to Tx.QFloat: %s %# v`, query, params)
	return
}

// fetch a row to as Base struct
func (tx *Tx) QBase(table string, id string) (base Base) {
	query := `SELECT * FROM ` + S.ZZ(table) + ` WHERE id = ` + Z(id)
	rows := tx.QAll(query)
	defer rows.Close()
	res := rows.ScanMap()
	if len(res) > 0 {
		base.FromMap(res)
	}
	base.Table = table
	return
}

//// fetch a row to Base struct by unique_id
//func (tx *Tx) QBaseUniq(table, uid string) (base Base) {
//	// not implemented: no unique constraint on scylladb
//}

// execute anything that doesn't need LastInsertId or RowsAffected
func (tx *Tx) DoExec(query string, params ...interface{}) error {
	err := tx.Trans.Query(query, params...).Exec()
	L.PanicIf(err, query)
	return err
}

// generate insert command and execute it
func (tx *Tx) DoInsert(actor string, table string, kvparams M.SX) (ok bool) {
	return tx.DoForcedInsert(actor, table, kvparams)
}

// do insert without checking
func (tx *Tx) DoForcedInsert(actor string, table string, kvparams M.SX) bool {
	kvparams[`created_by`] = actor
	kvparams[`created_at`] = T.UnixNano()
	kvparams[`id`] = NextId()
	query, params := GenInsert(table, kvparams)
	return nil != tx.DoExec(query, params...)
}

// generate insert or update command and execute it
func (tx *Tx) DoUpsert(actor string, table string, kvparams M.SX) (ok bool) {
	exists := tx.QExists(table, kvparams.GetStr(`id`))
	if exists {
		ok = tx.DoUpdate(actor, table, kvparams)
	} else {
		ok = tx.DoForcedInsert(actor, table, kvparams)
		if !ok {
			kvparams[`id`] = ``
		}
	}
	return
}

// generate update command and execute it
func (tx *Tx) DoUpdate(actor string, table string, kvparams M.SX) bool {
	id := kvparams.GetStr(`id`)
	if id != `` {
		kvparams[`updated_by`] = actor
		kvparams[`updated_at`] = T.UnixNano()
		query, params := GenUpdateId(table, id, kvparams)
		return nil != tx.DoExec(query, params...)
	}
	panic(`tx.DoUpdate missing id`)
	return false
}

// execute delete (is_deleted = true)
func (tx *Tx) DoDelete(actor string, table string, id string) bool {
	kvparams := M.SX{}
	kvparams[`is_deleted`] = true
	kvparams[`deleted_by`] = actor
	query, params := GenUpdateId(table, id, kvparams)
	return tx.DoExec(query, params...) == nil
}

// execute delete (is_deleted = false)
func (tx *Tx) DoRestore(actor string, table string, id string) bool {
	kvparams := M.SX{}
	kvparams[`is_deleted`] = false
	kvparams[`restored_by`] = actor
	query, params := GenUpdateId(table, id, kvparams)
	return tx.DoExec(query, params...) == nil
}

// delete or restore
func (tx *Tx) DoWipeUnwipe(a string, actor string, table string, id string) bool {
	if a == `save` || a == `` {
		return false
	}
	switch a {
	case `restore`:
		return tx.DoRestore(actor, table, id)
	case `delete`:
		return tx.DoDelete(actor, table, id)
	default:
		if S.StartsWith(a, `restore_`) {
			return tx.DoRestore(actor, table, id)
		} else if S.StartsWith(a, `delete_`) {
			return tx.DoDelete(actor, table, id)
		}
	}
	return false
}

// execute delete (is_deleted = true)
func (tx *Tx) DoDelete2(actor string, table string, id string, lambda func(base *Base) string, ajax W.Ajax) bool {
	base := tx.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Delete(actor)
	}
	ajax.Error(err_msg)
	return false
}

// execute delete (is_deleted = false)
func (tx *Tx) DoRestore2(actor string, table string, id string, lambda func(base *Base) string, ajax W.Ajax) bool {
	base := tx.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Restore(actor)
	}
	ajax.Error(err_msg)
	return false
}

// fetch json data as map
func (tx *Tx) DataJsonMap(table string, id string) M.SX {
	res := M.SX{}
	if id == `` {
		return res
	}
	query := `SELECT data FROM ` + table + ` WHERE id = ` + Z(id)
	rows := tx.QAll(query)
	defer rows.Close()
	return rows.ScanMap()
}

// fecth json data and id (if exists)
func (tx *Tx) DataJsonMapUniq(table, unique_id string) (res M.SX, id string) {
	res = tx.QFirstMap(`SELECT * FROM ` + table + ` WHERE id = ` + Z(unique_id))
	id = X.ToS(res[`id`])
	return
}

// query any number of columns, returns first line of line
func (db *Tx) QFirstMap(query string, params ...interface{}) M.SX {
	rows := db.QAll(query)
	defer rows.Close()
	return rows.ScanMap()
}
