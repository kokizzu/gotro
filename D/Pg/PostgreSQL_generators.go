package Pg

import (
	"bytes"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/M"
)

////////////////
// GENERATORS

// generate insert, requires table name and field-value params
func GenInsert(table string, kvparams M.SX) (string, []interface{}) {
	query := bytes.Buffer{}
	params := []interface{}{}
	query.WriteString(`INSERT INTO ` + table + `( `)
	len := 0
	for key, val := range kvparams {
		if len > 0 {
			query.WriteString(`, `)
		}
		query.WriteString(key)
		params = append(params, val)
		len++
	}
	query.WriteString(` ) VALUES ( `)
	for z := 1; z <= len; z++ {
		if z > 1 {
			query.WriteString(`, `)
		}
		query.WriteString(`$` + I.ToStr(z))
	}
	query.WriteString(` ) RETURNING id`)
	return query.String(), params
}

// generate update, requires table name, id and field-value params
func GenUpdateId(table string, id int64, kvparams M.SX) (string, []interface{}) {
	return GenUpdateWhere(table, `id = `+I.ToS(id), kvparams)
}

// generate update requires table name, unique id and field-value params
func GenUpdateUniq(table, uniq string, kvparams M.SX) (string, []interface{}) {
	return GenUpdateWhere(table, `unique_id = `+Z(uniq), kvparams)
}

// generate update, requires table, id and
func GenUpdateWhere(table, where string, kvparams M.SX) (string, []interface{}) {
	query := bytes.Buffer{}
	params := []interface{}{}
	query.WriteString(`UPDATE ` + table + ` SET `)
	len := 1
	for key, val := range kvparams {
		if key == `unique_id` && val == `` {
			continue
		}
		if len > 1 {
			query.WriteString(`, `)
		}
		query.WriteString(key + ` = $` + I.ToStr(len))
		params = append(params, val)
		len++
	}
	str := ` WHERE ` + where
	query.WriteString(str)
	return query.String(), params
}

// generate update base, requires table name, id and field-value, and data json string
// http://michael.otacoo.com/postgresql-2/manipulating-jsonb-data-with-key-unique/
// http://schinckel.net/2014/09/29/adding-json%28b%29-operators-to-postgresql/
// data json string provided by: data := M.ToJson(a K.M)
func GenUpdateBase(name string, id int64, kvparams M.SX, data string) (string, []interface{}) {
	query := bytes.Buffer{}
	query.WriteString(`
UPDATE ` + name + `
SET data = (
	SELECT json_object_agg(key, value)::jsonb
	FROM (
		SELECT * FROM jsonb_each( SELECT data FROM ` + name + ` WHERE id = ` + I.ToS(id) + ` )
		UNION ALL
	   SELECT * FROM jsonb_each( $1 )
	) x1
)`)
	params := []interface{}{data}
	len := 2
	for key, val := range kvparams {
		query.WriteString(`, ` + key + ` = $` + I.ToStr(len))
		params = append(params, val)
		len++
	}
	str := ` WHERE id = ` + I.ToS(id)
	query.WriteString(str)
	return query.String(), params
}

// generate update command to unset json data by key
func GenUpdateBaseUnset(name string, id int64, unsetkeys string) (string, []interface{}) {
	query := `
UPDATE ` + name + `
SET data = (
	SELECT json_object_agg(key, value)::jsonb FROM
		( SELECT * FROM jsonb_each( SELECT data FROM ` + name + ` WHERE id = ` + I.ToS(id) + ` )
			WHERE key NOT IN ( $1 )
		) x1
) WHERE id = ` + I.ToS(id)
	params := []interface{}{unsetkeys}
	return query, params
}
