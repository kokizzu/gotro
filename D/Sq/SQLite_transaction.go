package Sq

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
	"time"
)

type Tx struct {
	Trans *sqlx.Tx
}

// query 2 colums of integer-integer as map
func (tx *Tx) QIntIntMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		key := int64(0)
		val := key
		rows.Scan(&key, &val)
		res[key] = val
	}
	return res
}

// query 2 colums of integer-string as map
func (tx *Tx) QIntStrMap(query string, params ...interface{}) M.IS {
	res := M.IS{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		key := int64(0)
		val := ``
		rows.Scan(&key, &val)
		res[key] = val
	}
	return res
}

// query 1 colums of string as map
// SELECT data->>'puis_dosenid'
func (tx *Tx) QStrBoolMap(query string, params ...interface{}) M.SB {
	res := M.SB{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	key := ``
	for rows.Next() {
		rows.Scan(&key)
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
	for rows.Next() {
		rows.Scan(&key)
		res[key] = true
	}
	return res
}

// query 2 colums of string-string as map
func (tx *Tx) QStrStrMap(query string, params ...interface{}) M.SS {
	res := M.SS{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		key := ``
		val := ``
		rows.Scan(&key, &val)
		res[key] = val
	}
	return res
}

// query 1+N colums of string-[]any as map
func (tx *Tx) QStrArrMap(query string, params ...interface{}) M.SAX {
	res := M.SAX{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		row := rows.ScanSlice()
		res[X.ToS(row[0])] = row[1:]
	}
	return res
}

// query 2 colums of string-integer as map
func (tx *Tx) QStrIntMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		key := ``
		val := int64(0)
		rows.Scan(&key, &val)
		res[key] = val
	}
	return res
}

// query 1 colums of integer
func (tx *Tx) QIntCountMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := int64(0)
		rows.Scan(&val)
		res[val] += 1
	}
	return res
}

// query 1 colums of string
func (tx *Tx) QStrCountMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := ``
		rows.Scan(&val)
		res[val] += 1
	}
	return res
}

// query 1 colums of integer
func (tx *Tx) QIntArr(query string, params ...interface{}) []int64 {
	res := []int64{}
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := int64(0)
		rows.Scan(&val)
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
	for rows.Next() {
		val := ``
		rows.Scan(&val)
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
	var err error
	rs, err := tx.Trans.Queryx(query, params...)
	L.PanicIf(err, `failed to QAll: %s %# v`, query, params)
	rows = Records{rs, query, params}
	if DEBUG {
		L.TimeTrack(start, query)
	}
	return
}

// execute a select single value query, convert to string
func (tx *Tx) QStr(query string, params ...interface{}) (dest string) {
	err := tx.Trans.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QStr: %s %# v`, query, params)
	return
}

// check non-zero if unique exists
func (tx *Tx) QId(table, uniq string) (id int64) {
	uniq = Z(uniq)
	query := `SELECT COUNT(*) FROM ` + table + ` WHERE unique_id = ` + uniq
	id = tx.QInt(query)
	if id == 0 {
		return
	}
	query = `SELECT id FROM ` + table + ` WHERE unique_id = ` + uniq
	id = tx.QInt(query)
	return
}

// check if id exists
func (tx *Tx) QExists(table string, id int64) bool {
	id_str := I.ToS(id)
	query := `SELECT COUNT(*) FROM ` + table + ` WHERE id = ` + id_str
	count := tx.QInt(query)
	if count == 0 {
		return false
	}
	return true
}

// return non-empty string if id exists
func (tx *Tx) QUniq(table string, id int64) (uniq string) {
	id_str := I.ToS(id)
	query := `SELECT COUNT(*) FROM ` + table + ` WHERE id = ` + id_str
	count := tx.QInt(query)
	if count == 0 {
		return
	}
	query = `SELECT unique_id FROM ` + table + ` WHERE id = ` + id_str
	uniq = tx.QStr(query)
	return
}

// execute a select pair value query, convert to int64 and string
func (tx *Tx) QIntStr(query string, params ...interface{}) (i int64, s string) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&i, &s)
		return
	}
	return
}

// execute a select pair value query, convert to string and int
func (tx *Tx) QStrInt(query string, params ...interface{}) (s string, i int64) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&s, &i)
		return
	}
	return
}

// execute a select pair value query, convert to string and string
func (tx *Tx) QStrStr(query string, params ...interface{}) (s string, ss string) {
	rows := tx.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&s, &ss)
		return
	}
	return
}

// execute a select single value query, convert to int64
func (tx *Tx) QInt(query string, params ...interface{}) (dest int64) {
	err := tx.Trans.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QInt: %s %# v`, query, params)
	return
}

func (tx *Tx) QFloat(query string, params ...interface{}) (dest float64) {
	err := tx.Trans.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QFloat: %s %# v`, query, params)
	return
}

// fetch a row to as Base struct
func (tx *Tx) QBase(table string, id int64) (base Base) {
	query := `SELECT * FROM ` + S.ZZ(table) + ` WHERE id = ` + I.ToS(id)
	rows := tx.QAll(query)
	defer rows.Close()
	if rows.Next() {
		rows.ScanStruct(&base)
		base.DataToMap()
	}
	base.Table = table
	return
}

// fetch a row to Base struct by unique_id
func (tx *Tx) QBaseUniq(table, uid string) (base Base) {
	query := `SELECT * FROM ` + ZZ(table) + ` WHERE unique_id = ` + Z(uid)
	rows := tx.QAll(query)
	defer rows.Close()
	if rows.Next() {
		rows.ScanStruct(&base)
		base.DataToMap()
	}
	return
}

// execute anything that doesn't need LastInsertId or RowsAffected
func (tx *Tx) DoExec(query string, params ...interface{}) sql.Result {
	res, err := tx.Trans.Exec(query, params...)
	L.PanicIf(err, query)
	return res
}

// generate insert command and execute it
func (tx *Tx) DoInsert(actor int64, table string, kvparams M.SX) (id int64) {
	uniq, ok := kvparams[`unique_id`].(string)
	id = tx.QId(table, uniq)
	if ok && id > 0 {
		// already exists, cancel insert
		id = 0
		return
	}
	id = tx.DoForcedInsert(actor, table, kvparams)
	return
}

// do insert without checking
func (tx *Tx) DoForcedInsert(actor int64, table string, kvparams M.SX) (id int64) {
	kvparams[`created_by`] = actor
	query, params := GenInsert(table, kvparams)
	//L.Describe(query, params)
	rows := tx.QAll(query, params...)
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&id)
	} else {
		L.PanicIf(rows.Err(), `failed to insert on %s: %s; %# v`, table, query, params)
	}
	return
}

// generate insert or update command and execute it
func (tx *Tx) DoUpsert(actor int64, table string, kvparams M.SX) (id int64) {
	id, ok := kvparams[`id`].(int64)
	exists := false
	if ok && id > 0 {
		exists = tx.QExists(table, id)
	} else {
		uniq, ok := kvparams[`unique_id`].(string)
		id = tx.QId(table, uniq)
		exists = (ok && id > 0)
	}
	if exists {
		// update
		tx.DoUpdate(actor, table, id, kvparams)
	} else {
		// insert
		if id == 0 {
			delete(kvparams, `id`)
		}
		id = tx.DoForcedInsert(actor, table, kvparams)
	}
	return
}

// generate update command and execute it
func (tx *Tx) DoUpdate(actor int64, table string, id int64, kvparams M.SX) (ra int64) {
	if id > 0 {
		kvparams[`updated_by`] = actor
		query, params := GenUpdateId(table, id, kvparams)
		rs := tx.DoExec(query, params...)
		ra, _ = rs.RowsAffected()
		return
	}
	// when id not given, find the id first
	uniq, ok := kvparams[`unique_id`].(string)
	if !ok {
		// doesn't exists, cancel update
		ra = 0
		return
	}
	id = tx.QId(table, uniq)
	if id == 0 {
		// doesn't exists, cancel update
		ra = 0
		return
	}
	return tx.DoUpdate(actor, table, id, kvparams) // try again with given ID
}

// execute delete (is_deleted = true)
func (tx *Tx) DoDelete(actor int64, table string, id int64) bool {
	kvparams := M.SX{}
	kvparams[`is_deleted`] = true
	kvparams[`deleted_by`] = actor
	query, params := GenUpdateId(table, id, kvparams)
	rs := tx.DoExec(query, params...)
	ra, _ := rs.RowsAffected()
	return ra > 0
}

// execute delete (is_deleted = false)
func (tx *Tx) DoRestore(actor int64, table string, id int64) bool {
	kvparams := M.SX{}
	kvparams[`is_deleted`] = false
	kvparams[`restored_by`] = actor
	query, params := GenUpdateId(table, id, kvparams)
	rs := tx.DoExec(query, params...)
	ra, _ := rs.RowsAffected()
	return ra > 0
}

// delete or restore
func (tx *Tx) DoWipeUnwipe(a string, actor int64, table string, id int64) bool {
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
func (tx *Tx) DoDelete2(actor int64, table string, id int64, lambda func(base *Base) string, ajax W.Ajax) bool {
	base := tx.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Delete(actor)
	}
	ajax.Error(err_msg)
	return false
}

// execute delete (is_deleted = false)
func (tx *Tx) DoRestore2(actor int64, table string, id int64, lambda func(base *Base) string, ajax W.Ajax) bool {
	base := tx.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Restore(actor)
	}
	ajax.Error(err_msg)
	return false
}

// fetch json data as map
func (tx *Tx) DataJsonMap(table string, id int64) M.SX {
	res := M.SX{}
	if id <= 0 {
		return res
	}
	query := `SELECT data FROM ` + table + ` WHERE id = ` + I.ToS(id)
	rows := tx.QAll(query)
	str := ``
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&str)
		res = S.JsonToMap(str)
	}
	return res
}

// fecth json data as map by unique id
func (tx *Tx) DataJsonMapUniq(table, unique_id string) (res M.SX, id int64) {
	res = M.SX{}
	id = 0
	if unique_id == `` {
		return
	}
	query := `SELECT data, id FROM ` + table + ` WHERE unique_id = ` + Z(unique_id)
	rows := tx.QAll(query)
	str := ``
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&str, &id)
		res = S.JsonToMap(str)
	}
	return
}

// query any number of columns, returns first line of line
func (db *Tx) QFirstMap(query string, params ...interface{}) M.SX {
	rows := db.QAll(query)
	defer rows.Close()
	for rows.Next() {
		return rows.ScanMap()
	}
	return M.SX{}
}
