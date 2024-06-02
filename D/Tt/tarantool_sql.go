package Tt

import (
	"log"
	"time"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
	"github.com/tarantool/go-tarantool/v2"
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
	res, err := a.Connection.Do(tarantool.NewCallRequest("box.execute").Args(params)).Get()
	if L.IsError(err, `ExecSql box.execute failed: `+query) {
		log.Println(`ERROR ExecSql !!! ` + err.Error())
		//L.DescribeSql(query, parameters)
		L.Describe(parameters)
		//tracer.PanicOnDev(err)
		return map[any]any{`error`: err.Error()}
	}
	if len(res) > 0 {
		// go-tarantool/v2 use this:
		if tup, ok := res[0].(map[any]any); ok {
			// tup have metadata and rows for query
			return tup
		}
		if tup, ok := res[0].([][]any); ok {
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
	if len(res) > 1 {
		if tup, ok := res[1].([][]any); ok {
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
		// go-tarantool v2 use this kind of error
		if boxErr, ok := res[1].(*tarantool.BoxError); ok && boxErr != nil {
			L.Describe(query)
			L.Print(boxErr)
			return map[any]any{`error`: boxErr.Error()}
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

func QueryMetaFrom(res tarantool.Response, err error) QueryMeta {
	if res == nil {
		if err != nil {
			return QueryMeta{
				Err: err.Error(),
			}
		}
		return QueryMeta{}
	}
	exRes, ok := res.(*tarantool.ExecuteResponse)
	if !ok {
		if err != nil {
			return QueryMeta{
				Err: err.Error(),
			}
		}
		return QueryMeta{
			Err: `not ExecuteResponse`,
		}
	}
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	columns, _ := exRes.MetaData()
	sqlInfo, err := exRes.SQLInfo()
	if err != nil && errStr == `` {
		errStr = err.Error()
	}
	return QueryMeta{
		Columns: columns,
		SqlInfo: sqlInfo,
		Err:     errStr,
	}
}
