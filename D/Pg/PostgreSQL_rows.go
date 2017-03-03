package Pg

import (
	"github.com/jmoiron/sqlx"
	"github.com/kokizzu/gotro/L"
)

////////////////
// RECORD SET

type Rows struct {
	ResultSet   *sqlx.Rows
	Query       string
	QueryParams []interface{}
}

func (r *Rows) ErrorCheck(err error, msg string) {
	if len(r.QueryParams) == 0 {
		L.IsError(err, `failed `+msg, r.Query)
	} else {
		L.IsError(err, `failed `+msg, r.Query, r.QueryParams)
	}
}
func (r *Rows) Err() error {
	return r.ResultSet.Err()
}
func (r *Rows) Next() bool {
	return r.ResultSet.Next()
}
func (r *Rows) Close() {
	r.ResultSet.Close()
}
func (r *Rows) ScanSlice() []interface{} {
	arr, err := r.ResultSet.SliceScan()
	r.ErrorCheck(err, `ScanSlice`)
	for k, v := range arr {
		bs, ok := v.([]uint8)
		if ok {
			arr[k] = string(bs)
		}
	}
	return arr
}
func (r *Rows) ScanStruct(dest interface{}) bool {
	err := r.ResultSet.StructScan(dest)
	r.ErrorCheck(err, `StructScan`)
	return err == nil
}
func (r *Rows) Scan(dest ...interface{}) bool {
	err := r.ResultSet.Scan(dest...)
	r.ErrorCheck(err, `Scan`)
	return err == nil
}
func (r *Rows) ScanMap() map[string]interface{} {
	res := map[string]interface{}{}
	err := r.ResultSet.MapScan(res)
	r.ErrorCheck(err, `MapScan`)
	for k, v := range res {
		bs, ok := v.([]uint8)
		if ok {
			res[k] = string(bs)
		}
	}
	return res
}
