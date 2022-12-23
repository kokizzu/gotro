package Pg

import (
	"bytes"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kokizzu/gotro/W"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

///////////
// RDBMS

// wrapper for GO's sql.DB
type RDBMS struct {
	Name    string
	Adapter *sqlx.DB
}

// create new postgresql connection to localhost
func NewConn(user, db string) *RDBMS {
	opt := `user=` + user + ` dbname=` + db + ` sslmode=disable`
	conn := sqlx.MustConnect(`postgres`, opt)
	//conn.DB.SetMaxIdleConns(1)  // according to http://jmoiron.github.io/sqlx/#connectionPool
	conn.DB.SetMaxOpenConns(61) // TODO: change this according to postgresql.conf -3
	name := `pg::` + user + `@/` + db
	return &RDBMS{
		Name:    name,
		Adapter: conn,
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
		query := 'INSERT INTO ' || log_table || '( record_id, user_id, date, info, data_after, data_before )' || ' VALUES(' || OLD.id || ',' || NEW.updated_by || ',' || quote_literal(mod_time) || ',' || quote_literal(info) || ',' || quote_literal(NEW.data) || ',' || quote_literal(OLD.data) || ')';
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
func (db *RDBMS) CreateBaseTable(name, users_table string) {
	users_table = S.IfEmpty(users_table, `users`)
	db.DoTransaction(func(tx *Tx) string {
		query := `
CREATE TABLE IF NOT EXISTS ` + name + ` (
	id BIGSERIAL PRIMARY KEY,
	unique_id VARCHAR(4096) UNIQUE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE,
	restored_at TIMESTAMP WITH TIME ZONE,
	modified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	created_by BIGINT,
	updated_by BIGINT,
	deleted_by BIGINT,
	restored_by BIGINT,
	is_deleted BOOLEAN DEFAULT FALSE,
	data JSONB
);`
		is_deleted__index := name + `__is_deleted__index`
		modified_at__index := name + `__modified_at__index`
		unique_patern__index := name + `__unique__patern__index`
		query_count_index := `SELECT COUNT(*) FROM pg_indexes WHERE indexname = `
		ra, _ := tx.DoExec(query).RowsAffected()
		if ra > 0 {
			query = query_count_index + Z(is_deleted__index)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX IF NOT EXISTS ` + name + `__is_deleted__index ON ` + name + `(is_deleted);`
				tx.DoExec(query)
			}
			query = query_count_index + Z(modified_at__index)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX IF NOT EXISTS ` + name + `__modified_at__index ON ` + name + `(modified_at);`
				tx.DoExec(query)
			}
			query = query_count_index + Z(unique_patern__index)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX IF NOT EXISTS ` + name + `__unique__pattern ON ` + name + ` (unique_id varchar_pattern_ops);`
				tx.DoExec(query)
			}
			query = query_count_index + Z(name)
			if tx.QInt(query) == 0 {
				query = `CREATE INDEX IF NOT EXISTS ON ` + name + ` USING GIN(data)`
				tx.DoExec(query)
			}
		}
		trig_name := name + `__timestamp_changer`
		query = `SELECT COUNT(*) FROM pg_trigger WHERE tgname = ` + Z(trig_name)
		if tx.QInt(query) == 0 {
			query = `DROP TRIGGER IF EXISTS ` + trig_name + ` ON ` + name
			tx.DoExec(query)
			query = `CREATE TRIGGER ` + trig_name + ` BEFORE UPDATE ON ` + name + ` FOR EACH ROW EXECUTE PROCEDURE timestamp_changer();`
			tx.DoExec(query)
		}
		// logs
		query = `
CREATE TABLE IF NOT EXISTS _log_` + name + ` (
	id BIGSERIAL PRIMARY KEY,
	record_id BIGINT REFERENCES ` + name + `(id) ON UPDATE CASCADE,
	user_id BIGINT REFERENCES ` + users_table + `(id),
	date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	info TEXT,
	data_before JSONB NULL,
	data_after JSONB NULL
);`
		tx.DoExec(query)
		idx_name := `_log_` + name + `__record_id__idx`
		query = query_count_index + Z(idx_name)
		if tx.QInt(query) == 0 {
			query = `CREATE INDEX	IF NOT EXISTS ` + idx_name + ` ON 	_log_` + name + `	(record_id);`
			tx.DoExec(query)
		}

		idx_name = `_log_` + name + `__date__idx`
		query = query_count_index + Z(idx_name)
		if tx.QInt(query) == 0 {
			query = `CREATE INDEX	IF NOT EXISTS ` + idx_name + ` ON 	_log_` + name + `	(date);`
			tx.DoExec(query)
		}
		idx_name = `_log_` + name + `__user_id__idx`
		query = query_count_index + Z(idx_name)
		if tx.QInt(query) == 0 {
			query = `CREATE INDEX	IF NOT EXISTS ` + idx_name + ` ON 	_log_` + name + `	(user_id);`
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

// query any number of columns, returns array of slice (to be exported directly to json, not for processing)
func (db *RDBMS) QArray(query string, params ...any) A.X {
	rows := db.QAll(query)
	res := A.X{}
	defer rows.Close()
	for rows.Next() {
		res = append(res, rows.ScanSlice())
	}
	return res
}

// query any number of columns, returns first line of line
func (db *RDBMS) QFirstMap(query string, params ...any) M.SX {
	rows := db.QAll(query)
	defer rows.Close()
	for rows.Next() {
		return rows.ScanMap()
	}
	return M.SX{}
}

// query any number of columns, returns first line of line
func (db *RDBMS) QFirstArray(query string, params ...any) A.X {
	rows := db.QAll(query)
	defer rows.Close()
	for rows.Next() {
		return rows.ScanSlice()
	}
	return A.X{}
}

// query any number of columns, returns array of slice (to be exported directly to json, not for processing)
func (db *RDBMS) QMapArray(query string, params ...any) A.MSX {
	rows := db.QAll(query)
	res := A.MSX{}
	defer rows.Close()
	for rows.Next() {
		res = append(res, rows.ScanMap())
	}
	return res
}

// query to tsv file
func (db *RDBMS) QTsv(header, query string, params ...any) bytes.Buffer {
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
func (db *RDBMS) QStrIdxArrMap(key_idx int64, query string, params ...any) M.SAX {
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
func (db *RDBMS) QStrShiftArrMap(query string, params ...any) M.SAX {
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
func (db *RDBMS) QStrMapMap(key_idx string, query string, params ...any) M.SX {
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
func (db *RDBMS) QIntIntMap(query string, params ...any) M.II {
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
func (db *RDBMS) QIntStrMap(query string, params ...any) M.IS {
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
func (db *RDBMS) QStrBoolMap(query string, params ...any) M.SB {
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
func (db *RDBMS) QIntBoolMap(query string, params ...any) M.IB {
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
func (db *RDBMS) QStrStrMap(query string, params ...any) M.SS {
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
func (db *RDBMS) QStrArrMap(query string, params ...any) M.SAX {
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
func (db *RDBMS) QStrIntMap(query string, params ...any) M.SI {
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
func (db *RDBMS) QIntCountMap(query string, params ...any) M.II {
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
func (db *RDBMS) QStrCountMap(query string, params ...any) M.SI {
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
func (db *RDBMS) QIntArr(query string, params ...any) []int64 {
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
func (db *RDBMS) QJsonIntArr(query string, params ...any) []int64 {
	str := db.QStr(query, params...)
	ax := S.JsonToArr(str)
	ai := []int64{}
	for _, v := range ax {
		ai = append(ai, X.ToI(v))
	}
	return ai
}

// query 1 colums of string
func (db *RDBMS) QStrArr(query string, params ...any) []string {
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

// do query all, also calls DoPrepare, don't forget to close the rows
func (db *RDBMS) QAll(query string, params ...any) (rows Records) {
	var start time.Time
	if DEBUG {
		start = time.Now()
	}
	var err error

	rs, err := db.Adapter.Queryx(query, params...)
	L.PanicIf(err, `failed to QAll: %s %# v`, query, params)
	rows = Records{rs, query, params}
	if DEBUG {
		L.LogTrack(start, query)
	}
	return
}

// execute a select single value query, convert to string
func (db *RDBMS) QStr(query string, params ...any) (dest string) {
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
func (db *RDBMS) QIntStr(query string, params ...any) (i int64, s string) {
	db.DoTransaction(func(tx *Tx) string {
		i, s = tx.QIntStr(query, params...)
		return ``
	})
	return
}

// execute a select pair query, convert to string and int64
func (db *RDBMS) QStrInt(query string, params ...any) (s string, i int64) {
	db.DoTransaction(func(tx *Tx) string {
		s, i = tx.QStrInt(query, params...)
		return ``
	})
	return
}

// execute a select pair query, convert to string and string
func (db *RDBMS) QStrStr(query string, params ...any) (s string, ss string) {
	db.DoTransaction(func(tx *Tx) string {
		s, ss = tx.QStrStr(query, params...)
		return ``
	})
	return
}

// query single column string, return with true value

// execute a select single value query, convert to int64
func (db *RDBMS) QBool(query string, params ...any) (dest bool) {
	err := db.Adapter.Get(&dest, query, params...)
	L.PanicIf(err, `failed to QBool: %s %# v`, query, params)
	return
}

// execute a select single value query, convert to int64
func (db *RDBMS) QInt(query string, params ...any) (dest int64) {
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

func (db *RDBMS) QFloat(query string, params ...any) (dest float64) {
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

// generate insert command and execute it
func (db *RDBMS) DoInsert(actor int64, table string, kvparams M.SX) (id int64) {
	db.DoTransaction(func(tx *Tx) string {
		id = tx.DoInsert(actor, table, kvparams)
		return ``
	})
	return
}

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
		L.IsError(tx.Rollback(), `RDBMS.DoTransaction.Rollback`) // rollback when there is a string
	}()
	ok = lambda(&Tx{tx})
}

func (db *RDBMS) ViewExists(viewname string) bool {
	query := `SELECT COALESCE((SELECT COUNT(*) FROM information_schema.views WHERE table_name = ` + Z(viewname) + `),0)`
	return db.QInt(query) > 0
}

func (db *RDBMS) TableExists(tableName string) bool {
	query := `SELECT COALESCE((SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ` + Z(tableName) + `),0)`
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
