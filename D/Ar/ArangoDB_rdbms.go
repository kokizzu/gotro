package Ar

import (
	ara "github.com/diegogub/aranGO"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/X"
	"time"
)

///////////
// RDBMS

type RDBMS struct {
	Name    string
	Adapter *ara.Database
}

type Doc struct {
	ara.Document
	UniqueId   float64 `db:"unique_id"`
	CreatedAt  float64 `db:"created_at"`
	UpdatedAt  float64 `db:"updated_at"`
	DeletedAt  float64 `db:"deleted_at"`
	RestoredAt float64 `db:"restored_at"`
	ModifiedAt float64 `db:"modified_at"`
	CreatedBy  int64   `db:"created_by"`
	UpdatedBy  int64   `db:"updated_by"`
	DeletedBy  int64   `db:"deleted_by"`
	RestoredBy int64   `db:"restored_by"`
	IsDeleted  bool    `db:"is_deleted"`
	DataStr    string  `db:"data"` // json object
	XData      M.SX
}

// create new arangodb connection to localhost
func NewConn(host, user, pass, db string) *RDBMS {
	sess, err := ara.Connect(host, user, pass, false)
	if err != nil {
		panic(err)
	}
	conn := sess.DB(db)
	name := `arangodb::` + user + `@` + host + `/` + db
	return &RDBMS{
		Name:    name,
		Adapter: conn,
	}
}

// create a base table
func (db *RDBMS) CreateBaseTable(name string) {
	if !db.Adapter.ColExist(name) {
		collectionOptions := ara.NewCollectionOptions(name, true)
		db.Adapter.CreateCollection(collectionOptions)
	}
	if !db.Adapter.ColExist(`_log_` + name) {
		collectionOptions := ara.NewCollectionOptions(`_log_`+name, true)
		db.Adapter.CreateCollection(collectionOptions)
	}
}

// query any number of columns, returns array of slice (to be exported directly to json, not for processing)
func (db *RDBMS) QArray(query string, params ...interface{}) A.X {
	rows := db.QAll(query)
	res := A.X{}
	for _, row := range rows.Cursor.Result {
		res = append(res, row)
	}
	return res
}

// query any number of columns, returns first line of line
func (db *RDBMS) QFirstMap(query string, params ...interface{}) M.SX {
	rows := db.QAll(query)
	if len(rows.Cursor.Result) > 0 {
		return X.ToMSX(rows.Cursor.Result[0])
	}
	return M.SX{}
}

// query any number of columns, returns first line of line
func (db *RDBMS) QFirstArray(query string, params ...interface{}) A.X {
	rows := db.QAll(query)
	if len(rows.Cursor.Result) > 0 {
		return X.ToAX(rows.Cursor.Result[0])
	}
	return A.X{}
}

// query any number of columns, returns array of slice (to be exported directly to json, not for processing)
func (db *RDBMS) QMapArray(query string, params ...interface{}) A.MSX {
	rows := db.QAll(query)
	res := A.MSX{}
	for _, row := range rows.Cursor.Result {
		res = append(res, X.ToMSX(row))
	}
	return res
}

/*
// query to tsv file
func (db *RDBMS) QTsv(header, query string, params ...interface{}) bytes.Buffer {
	res := bytes.Buffer{}
	res.WriteString(header)
	rows := db.QAll(query)
	defer rows.Close()
	for rows.Next() {
		row := rows.ScanSlice()
		for k := range row {
			res.WriteString(X.ToS(row[k]))
			res.WriteRune('\t')
		}
		res.WriteRune('\n')
	}
	return res
}

// query any number of columns, returns map of string, array (to be exported directly to json, not for processing)
// the key_idx will be converted to string and taken as key
func (db *RDBMS) QStrIdxArrMap(key_idx int64, query string, params ...interface{}) M.SAX {
	rows := db.QAll(query)
	res := M.SAX{}
	defer rows.Close()
	for rows.Next() {
		row := rows.ScanSlice()
		key_str := X.ToS(row[key_idx])
		res[key_str] = row
	}
	return res
}

// query any number of columns, returns map of string, array (to be exported directly to json, not for processing)
// the first index will be converted to string and taken as key
func (db *RDBMS) QStrShiftArrMap(query string, params ...interface{}) M.SAX {
	rows := db.QAll(query)
	res := M.SAX{}
	defer rows.Close()
	for rows.Next() {
		row := rows.ScanSlice()
		key_str := X.ToS(row[0])
		res[key_str] = row[1:]
	}
	return res
}

// query any number of columns, returns map of string, map (to be exported directly to json, not for processing)
// the key_idx will be converted to string and taken as key
func (db *RDBMS) QStrMapMap(key_idx string, query string, params ...interface{}) M.SX {
	rows := db.QAll(query)
	res := M.SX{}
	defer rows.Close()
	for rows.Next() {
		row := rows.ScanMap()
		key_str := X.ToS(row[key_idx])
		res[key_str] = row
	}
	return res
}

// query 2 colums of integer-integer as map
// SELECT id, COUNT(*)
func (db *RDBMS) QIntIntMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := db.QAll(query, params...)
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
// SELECT id, unique_id
func (db *RDBMS) QIntStrMap(query string, params ...interface{}) M.IS {
	res := M.IS{}
	rows := db.QAll(query, params...)
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
// SELECT unique_id
func (db *RDBMS) QStrBoolMap(query string, params ...interface{}) M.SB {
	res := M.SB{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	key := ``
	for rows.Next() {
		rows.Scan(&key)
		res[key] = true
	}
	return res
}

// query single column int64, return with true value
func (db *RDBMS) QIntBoolMap(query string, params ...interface{}) M.IB {
	res := M.IB{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	key := int64(0)
	for rows.Next() {
		rows.Scan(&key)
		res[key] = true
	}
	return res
}

// query 2 colums of string-string as map
func (db *RDBMS) QStrStrMap(query string, params ...interface{}) M.SS {
	res := M.SS{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		key := ``
		val := ``
		rows.Scan(&key, &val)
		res[key] = val
	}
	return res
}

// query 1+N columns of string-[]any as map
func (db *RDBMS) QStrArrMap(query string, params ...interface{}) M.SAX {
	res := M.SAX{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		row := rows.ScanSlice()
		res[X.ToS(row[0])] = row[1:]
	}
	return res
}

// query 2 colums of string-integer as map
func (db *RDBMS) QStrIntMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		key := ``
		val := int64(0)
		rows.Scan(&key, &val)
		res[key] = val
	}
	return res
}

// query 1 colums of integer as map
func (db *RDBMS) QIntCountMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := int64(0)
		rows.Scan(&val)
		res[val] += 1
	}
	return res
}

// query 1 colums of string as map
// result equal to: SELECT name, COUNT(*) FROM tabel1 GROUP BY 1
// map[string]int
func (db *RDBMS) QStrCountMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := ``
		rows.Scan(&val)
		res[val] += 1
	}
	return res
}

// query 1 colums of integer
func (db *RDBMS) QIntArr(query string, params ...interface{}) []int64 {
	res := []int64{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := int64(0)
		rows.Scan(&val)
		res = append(res, val)
	}
	return res
}

// query one cell [1,2,3,...] and return array of int64
func (db *RDBMS) QJsonIntArr(query string, params ...interface{}) []int64 {
	str := db.QStr(query, params...)
	ax := S.JsonToArr(str)
	ai := []int64{}
	for _, v := range ax {
		ai = append(ai, X.ToI(v))
	}
	return ai
}

// query 1 colums of string
func (db *RDBMS) QStrArr(query string, params ...interface{}) []string {
	res := []string{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	for rows.Next() {
		val := ``
		rows.Scan(&val)
		res = append(res, val)
	}
	return res
}
*/

// do query all, also calls DoPrepare, don't forget to close the rows
func (db *RDBMS) QAll(query string, params ...interface{}) (rows Records) {
	var start time.Time
	if DEBUG {
		start = time.Now()
	}
	var err error
	q := ara.NewQuery(query)
	cursor, err := db.Adapter.Execute(q)
	L.PanicIf(err, `failed to QAll: %s %# v`, query, params)
	rows = Records{cursor, query, params}
	if DEBUG {
		L.LogTrack(start, query)
	}
	return
}

/*
// execute a select single value query, convert to string
func (db *RDBMS) QStr(query string, params ...interface{}) (dest string) {
	err := db.Adapter.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QStr: %s %# v`, query, params)
	return
}

// check if unique exists
func (db *RDBMS) QId(table, key string) (id int64) {
	db.DoTransaction(func(tx *Tx) string {
		id = tx.QId(table, key)
		return ``
	})
	return
}

// check if id exists
func (db *RDBMS) QExists(table string, id int64) (ex bool) {
	db.DoTransaction(func(tx *Tx) string {
		ex = tx.QExists(table, id)
		return ``
	})
	return
}

// get unique_id from id
func (db *RDBMS) QUniq(table string, key int64) (uniq string) {
	db.DoTransaction(func(tx *Tx) string {
		uniq = tx.QUniq(table, key)
		return ``
	})
	return
}

// execute a select pair query, convert to int64 and string
func (db *RDBMS) QIntStr(query string, params ...interface{}) (i int64, s string) {
	db.DoTransaction(func(tx *Tx) string {
		i, s = tx.QIntStr(query, params...)
		return ``
	})
	return
}

// execute a select pair query, convert to string and int64
func (db *RDBMS) QStrInt(query string, params ...interface{}) (s string, i int64) {
	db.DoTransaction(func(tx *Tx) string {
		s, i = tx.QStrInt(query, params...)
		return ``
	})
	return
}

// execute a select pair query, convert to string and string
func (db *RDBMS) QStrStr(query string, params ...interface{}) (s string, ss string) {
	db.DoTransaction(func(tx *Tx) string {
		s, ss = tx.QStrStr(query, params...)
		return ``
	})
	return
}

// query single column string, return with true value

// execute a select single value query, convert to int64
func (db *RDBMS) QBool(query string, params ...interface{}) (dest bool) {
	err := db.Adapter.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QBool: %s %# v`, query, params)
	return
}

// execute a select single value query, convert to int64
func (db *RDBMS) QInt(query string, params ...interface{}) (dest int64) {
	err := db.Adapter.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QInt: %s %# v`, query, params)
	return
}

// query count from a table
func (db *RDBMS) QCount(query string) (dest int64) {
	query = `
SELECT COUNT(*)
FROM (
	` + query + `
) count_sq0`
	err := db.Adapter.Get(&dest, query)
	L.PanicIf(err, `failed to QCount: %s`, query)
	return
}

func (db *RDBMS) QFloat(query string, params ...interface{}) (dest float64) {
	err := db.Adapter.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QFloat: %s %# v`, query, params)
	return
}

// fetch a row to as Base struct
func (db *RDBMS) QBase(table string, id int64) (base Base) {
	db.DoTransaction(func(tx *Tx) string {
		base = tx.QBase(table, id)
		base.Connection = db
		return ``
	})
	return
}

// fetch a row to Base struct by unique_id
func (db *RDBMS) QBaseUniq(table, uid string) (base Base) {
	db.DoTransaction(func(tx *Tx) string {
		base = tx.QBaseUniq(table, uid)
		base.Connection = db
		return ``
	})
	return
}
*/

// generate insert command and execute it
func (db *RDBMS) DoInsert(actor int64, table string, kvparams M.SX) (key string) {
	kvparams[`created_by`] = actor
	query, _ := GenInsert(table, kvparams)
	q := ara.NewQuery(query)
	L.Print(query)
	cursor, err := db.Adapter.Execute(q)
	L.PanicIf(err, `failed to DoInsert`)
	key = X.ToS(cursor.Result[0])
	return
}

/*
// generate insert or update command and execute it
func (db *RDBMS) DoUpsert(actor int64, table string, kvparams M.SX) (id int64) {
	id = 0
	db.DoTransaction(func(tx *Tx) string {
		id = tx.DoUpsert(actor, table, kvparams)
		return ``
	})
	return
}

// generate update command and execute it
func (db *RDBMS) DoUpdate(actor int64, table string, id int64, kvparams M.SX) (ra int64) {
	db.DoTransaction(func(tx *Tx) string {
		ra = tx.DoUpdate(actor, table, id, kvparams)
		return ``
	})
	return
}

// delete base table
func (db *RDBMS) DoDelete(actor int64, table string, id int64) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoDelete(actor, table, id)
		return ``
	})
	return
}

// delete or restore
func (db *RDBMS) DoWipeUnwipe(a string, actor int64, table string, id int64) bool {
	if a == `save` || a == `` {
		return false
	}
	switch a {
	case `restore`:
		return db.DoRestore(actor, table, id)
	case `delete`:
		return db.DoDelete(actor, table, id)
	default:
		if S.StartsWith(a, `restore_`) {
			return db.DoRestore(actor, table, id)
		} else if S.StartsWith(a, `delete_`) {
			return db.DoDelete(actor, table, id)
		}
	}
	return false
}

// restore base table
func (db *RDBMS) DoRestore(actor int64, table string, id int64) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoRestore(actor, table, id)
		return ``
	})
	return
}

// execute delete (is_deleted = true)
func (db *RDBMS) DoDelete2(actor int64, table string, id int64, lambda func(base D.Record) string, ajax W.Ajax) bool {
	base := db.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Delete(actor)
	}
	ajax.Error(err_msg)
	return false
}

// execute delete (is_deleted = false)
func (db *RDBMS) DoRestore2(actor int64, table string, id int64, lambda func(base D.Record) string, ajax W.Ajax) bool {
	base := db.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Restore(actor)
	}
	ajax.Error(err_msg)
	return false
}

// begin, commit transaction and rollback automatically when there are error
func (db *RDBMS) DoTransaction(lambda func(tx *Tx) string) {
	tx := db.Adapter.MustBegin()
	ok := ``
	defer func() {
		str := recover()
		if str != nil {
			// rollback when error
			err := tx.Rollback()
			L.PanicIf(errors.New(`transaction error`), `failed to end transaction %# v / %# v`, str, err)
			return
		}
		if ok == `` {
			// commit when empty string or nil
			err := tx.Commit()
			L.PanicIf(err, `failed to commit transaction`)
			return
		}
		L.Describe(ok)
		tx.Rollback() // rollback when there is a string
	}()
	ok = lambda(&Tx{tx})
}

func (db *RDBMS) ViewExists(viewname string) bool {
	query := `SELECT COALESCE((SELECT COUNT(*) FROM information_schema.views WHERE table_name = ` + Z(viewname) + `),0)`
	return db.QInt(query) > 0
}

// 2015-12-04 Kiz: replacement for JsonLine
// lambda should return empty string if it's correct row
func (db *RDBMS) JsonRow(table string, id int64, lambda func(rec D.Record) string) W.Ajax {
	base := db.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		ajax := base.XData
		ajax.Set(`unique_id`, base.UniqueId.String)
		ajax.Set(`is_deleted`, base.IsDeleted)
		return W.Ajax{ajax}
	} else {
		ajax := W.NewAjax()
		ajax.Error(err_msg)
		return ajax
	}
}
*/
