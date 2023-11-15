package saAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go

import (
	"database/sql"
	"net"
	"time"

	"example2/model/mAuth"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	chBuffer "github.com/kokizzu/ch-timed-buffer"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/L"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file saAuth__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type saAuth__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type saAuth__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type saAuth__ORM.GEN.go
// go:generate msgp -tests=false -file saAuth__ORM.GEN.go -o saAuth__MSG.GEN.go

var actionLogsDummy = ActionLogs{}
var Preparators = map[Ch.TableName]chBuffer.Preparator{
	mAuth.TableActionLogs: func(tx *sql.Tx) *sql.Stmt {
		query := actionLogsDummy.SqlInsert()
		stmt, err := tx.Prepare(query)
		L.IsError(err, `failed to tx.Prepare: `+query)
		return stmt
	},
}

type ActionLogs struct {
	Adapter    *Ch.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	CreatedAt  time.Time   `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	RequestId  string      `json:"requestId,string" form:"requestId" query:"requestId" long:"requestId" msg:"requestId"`
	ActorId    uint64      `json:"actorId,string" form:"actorId" query:"actorId" long:"actorId" msg:"actorId"`
	Action     string      `json:"action" form:"action" query:"action" long:"action" msg:"action"`
	StatusCode int16       `json:"statusCode" form:"statusCode" query:"statusCode" long:"statusCode" msg:"statusCode"`
	Traces     string      `json:"traces" form:"traces" query:"traces" long:"traces" msg:"traces"`
	Error      string      `json:"error" form:"error" query:"error" long:"error" msg:"error"`
	IpAddr4    net.IP      `json:"ipAddr4" form:"ipAddr4" query:"ipAddr4" long:"ipAddr4" msg:"ipAddr4"`
	IpAddr6    net.IP      `json:"ipAddr6" form:"ipAddr6" query:"ipAddr6" long:"ipAddr6" msg:"ipAddr6"`
	UserAgent  string      `json:"userAgent" form:"userAgent" query:"userAgent" long:"userAgent" msg:"userAgent"`
	Latency    float64     `json:"latency" form:"latency" query:"latency" long:"latency" msg:"latency"`
	TenantCode string      `json:"tenantCode" form:"tenantCode" query:"tenantCode" long:"tenantCode" msg:"tenantCode"`
	RefId      uint64      `json:"refId,string" form:"refId" query:"refId" long:"refId" msg:"refId"`
}

func NewActionLogs(adapter *Ch.Adapter) *ActionLogs {
	return &ActionLogs{Adapter: adapter}
}

// ActionLogsFieldTypeMap returns key value of field name and key
var ActionLogsFieldTypeMap = map[string]Ch.DataType{ //nolint:dupl false positive
	`createdAt`:  Ch.DateTime,
	`requestId`:  Ch.String,
	`actorId`:    Ch.UInt64,
	`action`:     Ch.String,
	`statusCode`: Ch.Int16,
	`traces`:     Ch.String,
	`error`:      Ch.String,
	`ipAddr4`:    Ch.IPv4,
	`ipAddr6`:    Ch.IPv6,
	`userAgent`:  Ch.String,
	`latency`:    Ch.Float64,
	`tenantCode`: Ch.String,
	`refId`:      Ch.UInt64,
}

func (a *ActionLogs) TableName() Ch.TableName { //nolint:dupl false positive
	return mAuth.TableActionLogs
}

func (a *ActionLogs) SqlTableName() string { //nolint:dupl false positive
	return `"actionLogs"`
}

func (a *ActionLogs) ScanRowAllCols(rows *sql.Rows) (err error) { //nolint:dupl false positive
	return rows.Scan(
		&a.CreatedAt,
		&a.RequestId,
		&a.ActorId,
		&a.Action,
		&a.StatusCode,
		&a.Traces,
		&a.Error,
		&a.IpAddr4,
		&a.IpAddr6,
		&a.UserAgent,
		&a.Latency,
		&a.TenantCode,
		&a.RefId,
	)
}

func (a *ActionLogs) ScanRowsAllCols(rows *sql.Rows, estimateRows int) (res []ActionLogs, err error) { //nolint:dupl false positive
	res = make([]ActionLogs, 0, estimateRows)
	defer rows.Close()
	for rows.Next() {
		var row ActionLogs
		err = row.ScanRowAllCols(rows)
		if err != nil {
			return
		}
		res = append(res, row)
	}
	return
}

// insert, error if exists
func (a *ActionLogs) SqlInsert() string { //nolint:dupl false positive
	return `INSERT INTO ` + a.SqlTableName() + `(` + a.SqlAllFields() + `) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
}

func (a *ActionLogs) SqlCount() string { //nolint:dupl false positive
	return `SELECT COUNT(*) FROM ` + a.SqlTableName()
}

func (a *ActionLogs) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` createdAt
	, requestId
	, actorId
	, action
	, statusCode
	, traces
	, error
	, ipAddr4
	, ipAddr6
	, userAgent
	, latency
	, tenantCode
	, refId
	`
}

func (a *ActionLogs) SqlAllFields() string { //nolint:dupl false positive
	return `createdAt, requestId, actorId, action, statusCode, traces, error, ipAddr4, ipAddr6, userAgent, latency, tenantCode, refId`
}

func (a ActionLogs) SqlInsertParam() []any { //nolint:dupl false positive
	return []any{
		a.CreatedAt,  // 0
		a.RequestId,  // 1
		a.ActorId,    // 2
		a.Action,     // 3
		a.StatusCode, // 4
		a.Traces,     // 5
		a.Error,      // 6
		a.IpAddr4,    // 7
		a.IpAddr6,    // 8
		a.UserAgent,  // 9
		a.Latency,    // 10
		a.TenantCode, // 11
		a.RefId,      // 12
	}
}

func (a *ActionLogs) IdxCreatedAt() int { //nolint:dupl false positive
	return 0
}

func (a *ActionLogs) SqlCreatedAt() string { //nolint:dupl false positive
	return `createdAt`
}

func (a *ActionLogs) IdxRequestId() int { //nolint:dupl false positive
	return 1
}

func (a *ActionLogs) SqlRequestId() string { //nolint:dupl false positive
	return `requestId`
}

func (a *ActionLogs) IdxActorId() int { //nolint:dupl false positive
	return 2
}

func (a *ActionLogs) SqlActorId() string { //nolint:dupl false positive
	return `actorId`
}

func (a *ActionLogs) IdxAction() int { //nolint:dupl false positive
	return 3
}

func (a *ActionLogs) SqlAction() string { //nolint:dupl false positive
	return `action`
}

func (a *ActionLogs) IdxStatusCode() int { //nolint:dupl false positive
	return 4
}

func (a *ActionLogs) SqlStatusCode() string { //nolint:dupl false positive
	return `statusCode`
}

func (a *ActionLogs) IdxTraces() int { //nolint:dupl false positive
	return 5
}

func (a *ActionLogs) SqlTraces() string { //nolint:dupl false positive
	return `traces`
}

func (a *ActionLogs) IdxError() int { //nolint:dupl false positive
	return 6
}

func (a *ActionLogs) SqlError() string { //nolint:dupl false positive
	return `error`
}

func (a *ActionLogs) IdxIpAddr4() int { //nolint:dupl false positive
	return 7
}

func (a *ActionLogs) SqlIpAddr4() string { //nolint:dupl false positive
	return `ipAddr4`
}

func (a *ActionLogs) IdxIpAddr6() int { //nolint:dupl false positive
	return 8
}

func (a *ActionLogs) SqlIpAddr6() string { //nolint:dupl false positive
	return `ipAddr6`
}

func (a *ActionLogs) IdxUserAgent() int { //nolint:dupl false positive
	return 9
}

func (a *ActionLogs) SqlUserAgent() string { //nolint:dupl false positive
	return `userAgent`
}

func (a *ActionLogs) IdxLatency() int { //nolint:dupl false positive
	return 10
}

func (a *ActionLogs) SqlLatency() string { //nolint:dupl false positive
	return `latency`
}

func (a *ActionLogs) IdxTenantCode() int { //nolint:dupl false positive
	return 11
}

func (a *ActionLogs) SqlTenantCode() string { //nolint:dupl false positive
	return `tenantCode`
}

func (a *ActionLogs) IdxRefId() int { //nolint:dupl false positive
	return 12
}

func (a *ActionLogs) SqlRefId() string { //nolint:dupl false positive
	return `refId`
}

func (a *ActionLogs) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		a.CreatedAt,  // 0
		a.RequestId,  // 1
		a.ActorId,    // 2
		a.Action,     // 3
		a.StatusCode, // 4
		a.Traces,     // 5
		a.Error,      // 6
		a.IpAddr4,    // 7
		a.IpAddr6,    // 8
		a.UserAgent,  // 9
		a.Latency,    // 10
		a.TenantCode, // 11
		a.RefId,      // 12
	}
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go
