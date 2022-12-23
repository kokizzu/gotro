package Pg

import (
	"github.com/jmoiron/sqlx"

	"github.com/kokizzu/gotro/L"
)

type Records struct {
	ResultSet   *sqlx.Rows
	Query       string
	QueryParams []any
}

func (r *Records) ErrorCheck(err error, msg string) {
	if len(r.QueryParams) == 0 {
		L.IsError(err, `failed `+msg, r.Query)
	} else {
		L.IsError(err, `failed `+msg, r.Query, r.QueryParams)
	}
}
func (r *Records) Err() error {
	return r.ResultSet.Err()
}
func (r *Records) Next() bool {
	return r.ResultSet.Next()
}
func (r *Records) Close() {
	r.ResultSet.Close()
}
func (r *Records) ScanSlice() []any {
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
func (r *Records) ScanStruct(dest any) bool {
	err := r.ResultSet.StructScan(dest)
	r.ErrorCheck(err, `StructScan`)
	return err == nil
}
func (r *Records) Scan(dest ...any) bool {
	err := r.ResultSet.Scan(dest...)
	r.ErrorCheck(err, `Scan`)
	return err == nil
}
func (r *Records) ScanMap() map[string]any {
	res := map[string]any{}
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
