package Sc

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kokizzu/gotro/L"
)

type Records struct {
	ResultSet   *gocql.Iter
	Query       string
	QueryParams []interface{}
}

func (r *Records) ErrorCheck(err error, msg string) {
	if len(r.QueryParams) == 0 {
		L.IsError(err, `failed `+msg, r.Query)
	} else {
		L.IsError(err, `failed `+msg, r.Query, r.QueryParams)
	}
}
func (r *Records) Err() error {
	warn := r.ResultSet.Warnings()
	if len(warn) > 0 {
		return fmt.Errorf(`%v`, warn)
	}
	return nil
}

func (r *Records) Close() {
	r.ResultSet.Close()
}

func (r *Records) CurrentSlice() []interface{} {
	rd, err := r.ResultSet.RowData()
	r.ErrorCheck(err, `failed Records.CurrentSlice`)
	return rd.Values
}

//func (r *Records) GetAX() []interface{} {
//	// not implemented: gocql does not have this kind of function
//}

func (r *Records) GetAMSX() []map[string]interface{} {
	arr, err := r.ResultSet.SliceMap()
	r.ErrorCheck(err, `failed Records.GetAMSX`)
	return arr
}

func (r *Records) ScanStruct(dest interface{}) bool {
	return r.ResultSet.Scan(dest)
}

func (r *Records) Scan(dest ...interface{}) bool {
	return r.ResultSet.Scan(dest...)
}

func (r *Records) ScanMap() map[string]interface{} {
	res := map[string]interface{}{}
	r.ResultSet.MapScan(res) // ignore error
	for k, v := range res {
		bs, ok := v.([]uint8)
		if ok {
			res[k] = string(bs)
		}
	}
	return res
}
