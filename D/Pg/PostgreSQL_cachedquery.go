package Pg

import (
	"fmt"
	"github.com/OneOfOne/cmap"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/W"
	"time"
)

// cached query module, all queries will be cached for TTL seconds per ram_key, expired per ram bucket
// bucket = the bucket of keys, all keys will be removed when it's expired
// ram_key = the key of the query

var CACHE_INV cmap.CMap  // table invalidate date
var CACHE_TTL cmap.CMap  // cache invalidate time
var CACHE_BORN cmap.CMap // cache born time
var CACHE cmap.CMap      // cache real storage

const TTL = 4

func init() {
	CACHE_INV = cmap.New()
	CACHE_BORN = cmap.New()
	CACHE_TTL = cmap.New()
	CACHE = cmap.New()
}

func RamDel_ByQuery(ram_key, query string) {
	CACHE.Delete(ram_key)
}

func RamGlobalEvict_ByAjax_ByBucket(ajax W.Ajax, bucket string) {
	if !ajax.HasError() {
		CACHE_INV.Set(bucket, T.UnixNano())
	}
}

func RamExpired_ByBucket_ByKey(bucket, ram_key string) bool {
	global := CACHE_INV.Get(bucket)
	if global != nil {
		born := CACHE_BORN.Get(ram_key)
		if born != nil {
			// when new data exists, invalidate all cache
			if born.(int64) <= global.(int64) {
				CACHE.Delete(ram_key)
				if DEBUG {
					L.Print(`RamExpired_ByBucket_ByKey: ` + bucket + ` ` + ram_key)
				}
				return true
			}
		}
	}
	local := CACHE_TTL.Get(ram_key)
	if local != nil {
		now := T.UnixNano()
		if local.(int64) <= now {
			CACHE.Delete(ram_key)
			if DEBUG {
				L.Print(`RamExpired_ByBucket_ByKey: ` + bucket + ` ` + ram_key)
			}
			return true
		}
	}
	return false
}

func RamSet_ByBucket_ByRamKey_ByQuery(bucket, ram_key, query string, val interface{}, sec int64) {
	if DEBUG {
		L.Print(`RamSet_ByBucket_ByRamKey_ByQuery: ` + bucket + ` ` + ram_key + fmt.Sprintf("\n%# v", val))
	}
	CACHE_BORN.Set(ram_key, T.UnixNano())
	CACHE.Set(ram_key, val)
	dur := time.Second * time.Duration(sec)
	CACHE_TTL.Set(ram_key, T.UnixNanoAfter(dur))
	go (func() {
		time.Sleep(dur)
		RamExpired_ByBucket_ByKey(bucket, ram_key)
	})()
}

// returns false when cache expired
func RamGet_ByRamKey_ByQuery(ram_key, query string) interface{} {
	res := CACHE.Get(ram_key)
	if res == nil {
		if DEBUG {
			L.Print(`RamGet_ByRamKey_ByQuery: (miss) ` + ram_key)
		}
	}
	return res
}

// M.SX
// only set when not yet being set, preventing double write and double delete
func RamSetMSX(bucket, ram_key, query string, val M.SX, sec int64) M.SX {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, ram_key, query, val, sec)
	return val
}

func RamGetMSX(ram_key, query string) (M.SX, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return M.SX{}, false
	}
	v, ok := val.(M.SX)
	return v, ok
}

// A.MSX
// only set when not yet being set, preventing double write and double delete
func RamSetAMSX(bucket, key, query string, val A.MSX, sec int64) A.MSX {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetAMSX(ram_key, query string) (A.MSX, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return A.MSX{}, false
	}
	v, ok := val.(A.MSX)
	return v, ok
}

// []STRING
// only set when not yet being set, preventing double write and double delete
func RamSetAS(bucket, key, query string, val []string, sec int64) []string {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetAS(ram_key, query string) ([]string, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return []string{}, false
	}
	v, ok := val.([]string)
	return v, ok
}

// []STRING
// only set when not yet being set, preventing double write and double delete
func RamSetAI(bucket, key, query string, val []int64, sec int64) []int64 {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetAI(ram_key, query string) ([]int64, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return []int64{}, false
	}
	v, ok := val.([]int64)
	return v, ok
}

// MAP[INT64]STRING
// only set when not yet being set, preventing double write and double delete
func RamSetMIS(bucket, key, query string, val M.IS, sec int64) M.IS {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetMIS(ram_key, query string) (M.IS, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return M.IS{}, false
	}
	v, ok := val.(M.IS)
	return v, ok
}

// MAP[STRING]STRING
// only set when not yet being set, preventing double write and double delete
func RamSetMSS(bucket, key, query string, val M.SS, sec int64) M.SS {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetMSS(ram_key, query string) (M.SS, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return M.SS{}, false
	}
	v, ok := val.(M.SS)
	return v, ok
}

// FLOAT64
// only set when not yet being set, preventing double write and double delete
func RamSetFloat(bucket, key, query string, val float64, sec int64) float64 {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetFloat(ram_key, query string) (float64, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return 0, false
	}
	v, ok := val.(float64)
	return v, ok
}

// BOOL
// only set when not yet being set, preventing double write and double delete
func RamSetBool(bucket, key, query string, val bool, sec int64) bool {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetBool(ram_key, query string) (bool, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return false, false
	}
	v, ok := val.(bool)
	return v, ok
}

// INT64
// only set when not yet being set, preventing double write and double delete
func RamSetInt(bucket, key, query string, val int64, sec int64) int64 {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetInt(ram_key, query string) (int64, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return 0, false
	}
	v, ok := val.(int64)
	return v, ok
}

// STRING
// only set when not yet being set, preventing double write and double delete
func RamSetStr(bucket, key, query string, val string, sec int64) string {
	RamSet_ByBucket_ByRamKey_ByQuery(bucket, key, query, val, sec)
	return val
}

func RamGetStr(ram_key, query string) (string, bool) {
	val := RamGet_ByRamKey_ByQuery(ram_key, query)
	if val == nil {
		return ``, false
	}
	v, ok := val.(string)
	return v, ok
}

// query and cache AMSX for TTL seconds
func (conn *RDBMS) CQMapArray(bucket, ram_key, query string) A.MSX {
	if val, ok := RamGetAMSX(ram_key, query); ok {
		return val
	}
	val := conn.QMapArray(query)
	return RamSetAMSX(bucket, ram_key, query, val, TTL)
}

// query and cache M.SX for TTL seconds
func (conn *RDBMS) CQFirstMap(bucket, ram_key, query string) M.SX {
	if val, ok := RamGetMSX(ram_key, query); ok {
		return val
	}
	val := conn.QFirstMap(query)
	return RamSetMSX(bucket, ram_key, query, val, TTL)
}

// query and cache []string for TTL seconds
func (conn *RDBMS) CQStrArr(bucket, ram_key, query string) []string {
	if val, ok := RamGetAS(ram_key, query); ok {
		return val
	}
	val := conn.QStrArr(query)
	return RamSetAS(bucket, ram_key, query, val, TTL)
}

// query and cache bool for TTL seconds
func (conn *RDBMS) CQBool(bucket, ram_key, query string) bool {
	if val, ok := RamGetBool(ram_key, query); ok {
		return val
	}
	val := conn.QBool(query)
	return RamSetBool(bucket, ram_key, query, val, TTL)
}

// query and cache int64 for TTL seconds
func (conn *RDBMS) CQInt(bucket, ram_key, query string) int64 {
	if val, ok := RamGetInt(ram_key, query); ok {
		return val
	}
	val := conn.QInt(query)
	return RamSetInt(bucket, ram_key, query, val, TTL)
}

// Query and cache float64 for TTL seconds
func (conn *RDBMS) CQFloat(bucket, ram_key, query string) float64 {
	if val, ok := RamGetFloat(ram_key, query); ok {
		return val
	}
	val := conn.QFloat(query)
	return RamSetFloat(bucket, ram_key, query, val, TTL)
}

// query and cache string for TTL seconds
func (conn *RDBMS) CQStr(bucket, ram_key, query string) string {
	if val, ok := RamGetStr(ram_key, query); ok {
		return val
	}
	val := conn.QStr(query)
	return RamSetStr(bucket, ram_key, query, val, TTL)
}

// query and cache M.SX for TTL seconds
func (conn *RDBMS) CQStrMapMap(bucket, ram_key, index, query string) M.SX {
	if val, ok := RamGetMSX(ram_key, query); ok {
		return val
	}
	val := conn.QStrMapMap(index, query)
	return RamSetMSX(bucket, ram_key, query, val, TTL)
}

// query and cache []int64 for TTL seconds
func (conn *RDBMS) CQIntArr(bucket, ram_key, query string) []int64 {
	if val, ok := RamGetAI(ram_key, query); ok {
		return val
	}
	val := conn.QIntArr(query)
	return RamSetAI(bucket, ram_key, query, val, TTL)
}

// query and cache map[int64]string for TTL seconds
func (conn *RDBMS) CQIntStrMap(bucket, ram_key, query string) map[int64]string {
	if val, ok := RamGetMIS(ram_key, query); ok {
		return val
	}
	val := conn.QIntStrMap(query)
	return RamSetMIS(bucket, ram_key, query, val, TTL)
}

// query and cache map[string]string for TTL seconds
func (conn *RDBMS) CQStrStrMap(bucket, ram_key, query string) map[string]string {
	if val, ok := RamGetMSS(ram_key, query); ok {
		return val
	}
	val := conn.QStrStrMap(query)
	return RamSetMSS(bucket, ram_key, query, val, TTL)
}
