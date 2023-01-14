package Tt

import (
	"log"
	"time"

	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

/*
function T(sql_statement)
	start_time = os.clock()
    res = box.execute(sql_statement)
	end_time = os.clock()
	return box.tuple.new(res, 'Query done in ' .. string.format("%.2f",(end_time - start_time)*1000) .. ' ms')
end
*/

func (a *Adapter) ExecSql(query string, parameters ...MSX) map[any]any {
	// https://www.tarantool.io/en/doc/latest/reference/reference_lua/box_sql/#box-sql-box-execute
	params := A.X{query}
	for _, v := range parameters {
		params = append(params, v)
	}
	//L.Describe(params)
	res, err := a.Call(`box.execute`, params)
	if L.IsError(err, `ExecSql box.execute failed: `+query) {
		log.Println(`ERROR ExecSql !!! ` + err.Error())
		//L.DescribeSql(query, parameters)
		L.Describe(parameters)
		//tracer.PanicOnDev(err)
		return map[any]any{`error`: err.Error()}
	}
	tup := res.Tuples()
	if len(tup) > 0 {
		if len(tup[0]) > 0 {
			if tup[0][0] != nil {
				kv, ok := tup[0][0].(map[any]any)
				// row_count for UPDATE
				// metadata, rows for SELECT
				if ok {
					return kv
				}
			}
		}
	}
	// possible error
	if len(tup) > 1 {
		if len(tup[1]) > 0 {
			if tup[1][0] != nil {
				errStr := X.ToS(tup[1][0])
				log.Println(`ERROR ExecSql syntax: ` + errStr)
				//L.DescribeSql(query, parameters)
				L.Describe(query)
				L.Describe(parameters)
				//tracer.PanicOnDev(errors.New(errStr))
				return map[any]any{`error`: tup[1][0]}
			}
		}
	}
	return map[any]any{}
}

func (a *Adapter) QuerySql(query string, callback func(row []any), parameters ...MSX) []any {
	if DebugPerf {
		defer L.TimeTrack(time.Now(), query)
	}
	kv := a.ExecSql(query, parameters...)
	rows, ok := kv[`rows`].([]any)
	if ok {
		for _, v := range rows {
			callback(v.([]any))
		}
		return rows
	}
	return nil
}

var DebugPerf = false

type QueryMeta struct {
	Columns []tarantool.ColumnMetaData
	SqlInfo tarantool.SQLInfo
	Err     string
	Code    uint32
}

func QueryMetaFrom(res *tarantool.Response, err error) QueryMeta {
	if res == nil {
		return QueryMeta{
			Err: err.Error(),
		}
	}
	return QueryMeta{
		Columns: res.MetaData,
		SqlInfo: res.SQLInfo,
		Err:     res.Error,
		Code:    res.Code,
	}
}
