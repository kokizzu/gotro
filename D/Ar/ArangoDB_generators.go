package Ar

import (
	"bytes"
	"github.com/kokizzu/gotro/M"
)

// generate insert, requires table name and field-value params
func GenInsert(table string, kvparams M.SX) (string, []interface{}) {
	query := bytes.Buffer{}
	params := []interface{}{}
	query.WriteString(`INSERT ` + M.ToJson(kvparams) + ` IN ` + table + ` RETURN NEW._key`)
	for _, val := range kvparams {
		params = append(params, val)
	}
	return query.String(), params
}

/*
// generate update, requires table name, id and field-value params
func GenUpdateId(table string, id int64, kvparams M.SX) (string, []interface{}) {
	return GenUpdateWhere(table, `id = `+I.ToS(id), kvparams)
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
*/