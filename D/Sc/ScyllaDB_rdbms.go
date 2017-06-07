package Sc

import (
	"bytes"
	"errors"
	"github.com/go-lang-plugin-org/go-lang-idea-plugin/testData/mockSdk-1.1.2/src/pkg/fmt"
	"github.com/gocql/gocql"
	_ "github.com/gocql/gocql"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
	"time"
)

// TODO: change postgresql specific tables (pg_indexes) and trigger (InitTrigger) syntax to scylladb's

// wrapper for GO's sql.DB
type RDBMS struct {
	Name    string
	Cluster *gocql.ClusterConfig
	Session *gocql.Session
}

// create new scylla connection to localhost
// CREATE KEYSPACE "replace_this" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
func NewConn(user, pass, ip, db string) *RDBMS {
	clust := gocql.NewCluster(ip)
	clust.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	clust.Timeout = 8 * time.Second
	clust.ConnectTimeout = 15 * time.Second
	clust.Keyspace = db
	if user != `` {
		clust.Authenticator = gocql.PasswordAuthenticator{
			Username: user,
			Password: pass,
		}
	}
	sess, err := clust.CreateSession()
	L.PanicIf(err, `Failed create session Sc`)
	return &RDBMS{
		Name:    `sc://` + user + `:` + pass + `@` + ip + `/` + db,
		Cluster: clust,
		Session: sess,
	}
}

// rename a base table
func (db *RDBMS) RenameBaseTable(oldname, newname string) {
	db.DoTransaction(func(tx *Tx) string {
		query := `ALTER TABLE ` + ZZ(oldname) + ` RENAME TO ` + ZZ(newname)
		tx.DoExec(query)
		return ``
	})
}

// init trigger: not supported
//func (db *RDBMS) InitTrigger() {}

// create a base table
func (db *RDBMS) CreateBaseTable(name string) {
	db.DoTransaction(func(tx *Tx) string {
		query := `
CREATE TABLE IF NOT EXISTS ` + name + ` (
	clust TEXT, 
	id TEXT,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,
	restored_at TIMESTAMP,
	modified_at TIMESTAMP,
	created_by TEXT,
	updated_by TEXT,
	deleted_by TEXT,
	restored_by TEXT,
	is_deleted BOOLEAN,
	data MAP<TEXT,TEXT>,
	PRIMARY KEY (clust, id)
);`
		tx.DoExec(query)
		return ``
	})

}

//// query any number of columns, returns array of slice (to be exported directly to json, not for processing)
//func (db *RDBMS) QArray(query string, params ...interface{}) A.X {
//	// not implemented, far too inefficient with current gocql
//}

// query any number of columns, returns first line of line
func (db *RDBMS) QFirstMap(query string, params ...interface{}) M.SX {
	rows := db.QAll(query)
	defer rows.Close()
	return rows.ScanMap()
}

// query any number of columns, returns first line of line
func (db *RDBMS) QFirstArray(query string, params ...interface{}) A.X {
	rows := db.QAll(query)
	defer rows.Close()
	return rows.CurrentSlice()
}

// query any number of columns, returns array of slice (to be exported directly to json, not for processing)
func (db *RDBMS) QMapArray(query string, params ...interface{}) A.MSX {
	rows := db.QAll(query)
	res := A.MSX{}
	defer rows.Close()
	for {
		m := rows.ScanMap()
		if len(m) == 0 {
			break
		}
		res = append(res, m)
	}
	return res
}

// query to tsv file
func (db *RDBMS) QTsv(header, query string, params ...interface{}) bytes.Buffer {
	res := bytes.Buffer{}
	res.WriteString(header)
	rows := db.QAll(query)
	defer rows.Close()
	all := rows.GetAMSX()
	col := rows.ResultSet.Columns()
	for _, row := range all {
		for _, k := range col {
			res.WriteString(X.ToS(row[k.Name]))
			res.WriteRune('\t')
		}
		res.WriteRune('\n')
	}
	return res
}

//// query any number of columns, returns map of string, array (to be exported directly to json, not for processing)
//// the key_idx will be converted to string and taken as key
//func (db *RDBMS) QStrIdxArrMap(key_idx int64, query string, params ...interface{}) M.SAX {
//	// not implemented, far too inefficient with current gocql
//}

//// query any number of columns, returns map of string, array (to be exported directly to json, not for processing)
//// the first index will be converted to string and taken as key
//func (db *RDBMS) QStrShiftArrMap(query string, params ...interface{}) M.SAX {
//	// not implemented, far too inefficient with current gocql
//}

// query any number of columns, returns map of string, map (to be exported directly to json, not for processing)
// the key_idx will be converted to string and taken as key
func (db *RDBMS) QStrMapMap(key_idx string, query string, params ...interface{}) M.SX {
	rows := db.QAll(query)
	res := M.SX{}
	defer rows.Close()
	all := rows.GetAMSX()
	for _, row := range all {
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
	key := int64(0)
	val := key
	for rows.Scan(&key, &val) {
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
	key := int64(0)
	val := ``
	for rows.Scan(&key, &val) {
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
	for rows.Scan(&key) {
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
	for rows.Scan(&key) {
		res[key] = true
	}
	return res
}

// query 2 colums of string-string as map
func (db *RDBMS) QStrStrMap(query string, params ...interface{}) M.SS {
	res := M.SS{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	key := ``
	val := ``
	for rows.Scan(&key, &val) {
		res[key] = val
	}
	return res
}

//// query 1+N columns of string-[]any as map
//func (db *RDBMS) QStrArrMap(query string, params ...interface{}) M.SAX {
//	// not implemented, far too inefficient with current gocql
//}

// query 2 colums of string-integer as map
func (db *RDBMS) QStrIntMap(query string, params ...interface{}) M.SI {
	res := M.SI{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	key := ``
	val := int64(0)
	for rows.Scan(&key, &val) {
		res[key] = val
	}
	return res
}

// query 1 colums of integer as map
func (db *RDBMS) QIntCountMap(query string, params ...interface{}) M.II {
	res := M.II{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	val := int64(0)
	for rows.Scan(&val) {
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
	val := ``
	for rows.Scan(&val) {
		res[val] += 1
	}
	return res
}

// query 1 colums of integer
func (db *RDBMS) QIntArr(query string, params ...interface{}) []int64 {
	res := []int64{}
	rows := db.QAll(query, params...)
	defer rows.Close()
	val := int64(0)
	for rows.Scan(&val) {
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
	val := ``
	for rows.Scan(&val) {
		res = append(res, val)
	}
	return res
}

// do query all, also calls DoPrepare, don't forget to close the rows
func (db *RDBMS) QAll(query string, params ...interface{}) (rows Records) {
	var start time.Time
	if DEBUG {
		start = time.Now()
	}
	db.DoTransaction(func(tx *Tx) string {
		rows = tx.QAll(query, params...)
		return ``
	})
	if DEBUG {
		L.LogTrack(start, query)
	}
	return
}

// execute query, since scylladb doesn't have transaction
func (db *RDBMS) DoExec(query string, params ...interface{}) error {
	err := db.Session.Query(query, params...).Exec()
	L.PanicIf(err, query)
	return err
}

// execute a select single value query, convert to string
func (db *RDBMS) QStr(query string, params ...interface{}) (dest string) {
	db.DoTransaction(func(tx *Tx) string {
		dest = tx.QStr(query, params...)
		return ``
	})
	return
}

//// check if unique exists
//func (db *RDBMS) QId(table, key string) (id int64) {
//	// not implemented: no unique constraint on scylladb
//}

// check if id exists
func (db *RDBMS) QExists(table string, id string) (ex bool) {
	db.DoTransaction(func(tx *Tx) string {
		ex = tx.QExists(table, id)
		return ``
	})
	return
}

//// get unique_id from id
//func (db *RDBMS) QUniq(table string, key int64) (uniq string) {
//	// not implemented: useless id is unique_id, and there are no additional unique constraint on scylladb
//}

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
	db.DoTransaction(func(tx *Tx) string {
		dest = tx.QBool(query, params...)
		return ``
	})
	return
}

// execute a select single value query, convert to int64
func (db *RDBMS) QInt(query string, params ...interface{}) (dest int64) {
	db.DoTransaction(func(tx *Tx) string {
		dest = tx.QInt(query, params...)
		return ``
	})
	return
}

// query count from a table
func (db *RDBMS) QCount(query string) (dest int64) {
	query = `
SELECT COUNT(*)
FROM (
	` + query + `
) count_sq0`
	db.DoTransaction(func(tx *Tx) string {
		dest = tx.QInt(query)
		return ``
	})
	return
}

func (db *RDBMS) QFloat(query string, params ...interface{}) (dest float64) {
	db.DoTransaction(func(tx *Tx) string {
		dest = tx.QFloat(query)
		return ``
	})
	return
}

// fetch a row to as Base struct
func (db *RDBMS) QBase(table string, id string) (base Base) {
	db.DoTransaction(func(tx *Tx) string {
		base = tx.QBase(table, id)
		base.Connection = db
		return ``
	})
	return
}

//// fetch a row to Base struct by unique_id
//func (db *RDBMS) QBaseUniq(table, uid string) (base Base) {
//	// not implemented: no unique constraint on scylladb
//}

// generate insert command and execute it
func (db *RDBMS) DoInsert(actor string, table string, kvparams M.SX) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoInsert(actor, table, kvparams)
		return ``
	})
	return
}

// generate insert or update command and execute it
func (db *RDBMS) DoUpsert(actor string, table string, kvparams M.SX) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoUpsert(actor, table, kvparams)
		return ``
	})
	return
}

// generate update command and execute it
func (db *RDBMS) DoUpdate(actor string, table string, kvparams M.SX) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoUpdate(actor, table, kvparams)
		return ``
	})
	return
}

// delete base table
func (db *RDBMS) DoDelete(actor string, table string, id string) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoDelete(actor, table, id)
		return ``
	})
	return
}

// delete or restore
func (db *RDBMS) DoWipeUnwipe(a string, actor string, table string, id string) bool {
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
func (db *RDBMS) DoRestore(actor string, table string, id string) (ok bool) {
	db.DoTransaction(func(tx *Tx) string {
		ok = tx.DoRestore(actor, table, id)
		return ``
	})
	return
}

// execute delete (is_deleted = true)
func (db *RDBMS) DoDelete2(actor string, table string, id string, lambda func(base D.Record) string, ajax W.Ajax) bool {
	base := db.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Delete(actor)
	}
	ajax.Error(err_msg)
	return false
}

// execute delete (is_deleted = false)
func (db *RDBMS) DoRestore2(actor string, table string, id string, lambda func(base D.Record) string, ajax W.Ajax) bool {
	base := db.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		return base.Restore(actor)
	}
	ajax.Error(err_msg)
	return false
}

// get connection from pool but does not rollback automatically when there are error
func (db *RDBMS) DoTransaction(lambda func(tx *Tx) string) (ok string) {
	tx, err := db.Cluster.CreateSession()
	L.PanicIf(err, `Failed to create session to keyspace`)
	defer func() {
		str := recover()
		tx.Close()
		if str != nil {
			L.PanicIf(errors.New(`transaction error`), `failed to end transaction %# v`, str)
			return
		}
		if ok != `` {
			L.PanicIf(fmt.Errorf("%v", str), `query executed but has error`)
			return
		}
		L.Describe(ok)
	}()
	ok = lambda(&Tx{tx})
	return
}

func (db *RDBMS) ViewExists(viewname string) bool {
	query := `SELECT COALESCE((SELECT COUNT(*) FROM information_schema.views WHERE table_name = ` + Z(viewname) + `),0)`
	return db.QInt(query) > 0
}

// 2015-12-04 Kiz: replacement for JsonLine
// lambda should return empty string if it's correct row
func (db *RDBMS) JsonRow(table string, id string, lambda func(rec D.Record) string) W.Ajax {
	base := db.QBase(table, id)
	err_msg := lambda(&base)
	if err_msg == `` {
		ajax := base.XData
		ajax.Set(`id`, base.Id)
		ajax.Set(`is_deleted`, base.IsDeleted)
		return W.Ajax{ajax}
	} else {
		ajax := W.NewAjax()
		ajax.Error(err_msg)
		return ajax
	}
}
