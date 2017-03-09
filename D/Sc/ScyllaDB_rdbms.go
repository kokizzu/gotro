package Sc

import (
	"bytes"
	"errors"
	"github.com/go-lang-plugin-org/go-lang-idea-plugin/testData/mockSdk-1.1.2/src/pkg/fmt"
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	_ "github.com/gocql/gocql"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/I"
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
	Adapter gocassa.KeySpace
	Cluster *gocql.ClusterConfig
}

// create new postgresql connection to localhost
func NewConn(user, pass, ip, db string) *RDBMS {
	conn, err := gocassa.ConnectToKeySpace(db, []string{ip}, user, pass)
	L.PanicIf(err, `Failed connect to keyspace`)
	clust := gocql.NewCluster(ip)
	clust.Authenticator = gocql.PasswordAuthenticator{
		Username: user,
		Password: pass,
	}
	clust.Keyspace = db
	return &RDBMS{
		Name:    `sc://` + user + `:` + pass + `@` + ip + `/` + db,
		Adapter: conn,
		Cluster: clust,
	}
}

// rename a base table
func (db *RDBMS) RenameBaseTable(oldname, newname string) {
	db.DoTransaction(func(tx *Tx) string {
		query := `ALTER TABLE ` + ZZ(oldname) + ` RENAME TO ` + ZZ(newname)
		tx.DoExec(query)
		oldtrig := oldname + `__trigger`
		newtrig := newname + `__trigger`
		query = `ALTER TRIGGER ` + ZZ(oldtrig) + ` ON ` + ZZ(newname) + ` RENAME TO ` + ZZ(newtrig)
		tx.DoExec(query)
		oldseq := oldname + `_id_seq`
		newseq := newname + `_id_seq`
		query = `ALTER SEQUENCE ` + ZZ(oldseq) + ` RENAME TO ` + ZZ(newseq)
		tx.DoExec(query)
		oldlog := `_log_` + oldname
		newlog := `_log_` + newname
		query = `ALTER TABLE ` + ZZ(oldlog) + ` RENAME TO ` + ZZ(newlog)
		tx.DoExec(query)
		oldseq = oldlog + `_id_seq`
		newseq = newlog + `_id_seq`
		query = `ALTER SEQUENCE ` + ZZ(oldseq) + ` RENAME TO ` + ZZ(newseq)
		tx.DoExec(query)
		return ``
	})
}

// init trigger
func (db *RDBMS) InitTrigger() {
	// auto update timestamp trigger
	query := `
CREATE OR REPLACE FUNCTION timestamp_changer() RETURNS trigger AS $$
DECLARE
	changed BOOLEAN  := FALSE;
	log_table TEXT := quote_ident('_log_' || TG_TABLE_NAME);
	info TEXT := '';
	mod_time TIMESTAMP := CURRENT_TIMESTAMP;
	actor BIGINT;
	query TEXT := '';
BEGIN
	IF (OLD.unique_id <> NEW.unique_id) THEN
		NEW.updated_at := mod_time;
		actor := NEW.updated_by;
		changed := TRUE;
		IF info <> '' THEN info := info || chr(10); END IF;
		info := info || 'unique' || E'\t' || OLD.unique_id || E'\t' || NEW.unique_id;
	END IF;
	IF (OLD.is_deleted = TRUE) AND (NEW.is_deleted = FALSE) THEN
		NEW.restored_at := mod_time;
		actor := NEW.restored_by;
		IF info <> '' THEN info := info || chr(10); END IF;
		info := info || 'restore';
		changed := TRUE;
	END IF;
	IF (OLD.is_deleted = FALSE) AND (NEW.is_deleted = TRUE) THEN
		NEW.deleted_at := mod_time;
		actor := NEW.deleted_by;
		IF info <> '' THEN info := info || chr(10); END IF;
		info := info || 'delete';
		changed := TRUE;
	END IF;
	IF (OLD.data <> NEW.data) THEN
		NEW.updated_at := mod_time;
		IF info <> '' THEN info := info || chr(10); END IF;
		info := info || 'update';
		query := 'INSERT INTO ' || log_table || '( record_id, user_id, date, info, data_before, data_after )' || ' VALUES(' || OLD.id || ',' || NEW.updated_by || ',' || quote_literal(mod_time) || ',' || quote_literal(info) || ',' || quote_literal(NEW.data) || ',' || quote_literal(OLD.data) || ')';
		EXECUTE query;
		changed := TRUE;
	ELSEIF changed THEN
		query := 'INSERT INTO ' || log_table || '( record_id, user_id, date, info )' || ' VALUES(' || OLD.id || ',' || actor || ',' || quote_literal(mod_time) || ',' || quote_literal(info) || ')';
		EXECUTE query;
	END IF;
	IF changed THEN NEW.modified_at := mod_time; END IF;
	RETURN NEW;
END; $$ LANGUAGE plpgsql;`
	db.DoTransaction(func(tx *Tx) string {
		tx.DoExec(query)
		return ``
	})
}

// create a base table
func (db *RDBMS) CreateBaseTable(name string) {
	db.DoTransaction(func(tx *Tx) string {
		query := `
CREATE TABLE IF NOT EXISTS ` + name + ` (
	id PRIMARY KEY NOT NULL AUTO_INCREMENT,
	unique_id VARCHAR(120) UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP,
	deleted_at TIMESTAMP,
	restored_at TIMESTAMP,
	modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_by BIGINT,
	updated_by BIGINT,
	deleted_by BIGINT,
	restored_by BIGINT,
	is_deleted BOOLEAN DEFAULT FALSE,
	data JSON
);`
		is_deleted__index := name + `__is_deleted__index`
		modified_at__index := name + `__modified_at__index`
		unique_patern__index := name + `__unique__patern__index`
		query_count_index := `SELECT COUNT(*) FROM pg_indexes WHERE indexname = `
		err := tx.DoExec(query)
		if err == nil {
			query = query_count_index + Z(is_deleted__index)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX ` + name + `__is_deleted__index ON ` + name + `(is_deleted);`
				tx.DoExec(query)
			}
			query = query_count_index + Z(modified_at__index)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX ` + name + `__modified_at__index ON ` + name + `(modified_at);`
				tx.DoExec(query)
			}
			query = query_count_index + Z(unique_patern__index)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX ` + name + `__unique__pattern ON ` + name + ` (unique_id varchar_pattern_ops);`
				tx.DoExec(query)
			}
			query = query_count_index + Z(name)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX ON ` + name + ` USING GIN(data)`
				tx.DoExec(query)
			}
		}
		trig_name := name + `__timestamp_changer`
		query = `SELECT COUNT(*) FROM pg_trigger WHERE tgname = ` + Z(trig_name)
		if tx.QInt(query) == 0 {
			query = `DROP TRIGGER IF EXISTS ` + trig_name + ` ON ` + name
			tx.DoExec(query)
			query = `CREATE TRIGGER ` + trig_name + ` BEFORE UPDATE ON ` + name + ` FOR EACH ROW EXECUTE PROCEDURE timestamp_changer();`
		}
		tx.DoExec(query)
		// logs
		if name == `users` {
			// TODO: tambahkan record ke access_log ketika login, renew session, logout
			query = `
CREATE TABLE IF NOT EXISTS access_logs (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	log_count BIGINT,
	session VARCHAR(256),
	ip_address VARCHAR(256),
	logs TEXT,
	CONSTRAINT unique__user_id__logs UNIQUE(user_id,session)
);`
			tx.DoExec(query)
		}
		query = `
CREATE TABLE IF NOT EXISTS _log_` + name + ` (
	id BIGSERIAL PRIMARY KEY,
	record_id BIGINT REFERENCES ` + name + `(id) ON UPDATE CASCADE,
	user_id BIGINT REFERENCES users(id),
	date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	info TEXT,
	data_before JSONB NULL,
	data_after JSONB NULL
);`
		tx.DoExec(query)
		idx_name := `_log_` + name + `__record_id__idx`
		query = query_count_index + Z(idx_name)
		if tx.QInt(query) == 0 {
			query = `CREATE INDEX	` + idx_name + ` ON 	_log_` + name + `	(record_id);`
			tx.DoExec(query)
		}

		idx_name = `_log_` + name + `__date__idx`
		query = query_count_index + Z(idx_name)
		if tx.QInt(query) == 0 {
			query = `CREATE INDEX	` + idx_name + ` ON 	_log_` + name + `	(date);`
			tx.DoExec(query)
		}
		idx_name = `_log_` + name + `__user_id__idx`
		query = query_count_index + Z(idx_name)
		if tx.QInt(query) == 0 {
			query = `CREATE INDEX	` + idx_name + ` ON 	_log_` + name + `	(user_id);`
			tx.DoExec(query)
		}
		return ``
	})

}

// reset sequence, to be called after TRUNCATE TABLE tablename
func (db *RDBMS) FixSerialSequence(table string) {
	db.DoTransaction(func(tx *Tx) string {
		next := tx.QInt(`SELECT COALESCE(MAX(id)+1,1) FROM ` + table)
		tx.DoExec(`ALTER SEQUENCE ` + table + `_id_seq RESTART WITH ` + I.ToS(next))
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
// SELECT id, puis_dosenid
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
// SELECT data->>'puis_dosenid'
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
