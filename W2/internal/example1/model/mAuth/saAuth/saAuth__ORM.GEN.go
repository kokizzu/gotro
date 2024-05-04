package saAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go

import (
	"database/sql"
	"net"
	"time"

	"example1/model/mAuth"

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

var userLogsDummy = UserLogs{}
var Preparators = map[Ch.TableName]chBuffer.Preparator{
	mAuth.TableUserLogs: func(tx *sql.Tx) *sql.Stmt {
		query := userLogsDummy.SqlInsert()
		stmt, err := tx.Prepare(query)
		L.IsError(err, `failed to tx.Prepare: `+query)
		return stmt
	},
}

type UserLogs struct {
	Adapter   *Ch.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	CreatedAt time.Time   `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	RequestId uint64      `json:"requestId,string" form:"requestId" query:"requestId" long:"requestId" msg:"requestId"`
	ActorId   uint64      `json:"actorId,string" form:"actorId" query:"actorId" long:"actorId" msg:"actorId"`
	Error     string      `json:"error" form:"error" query:"error" long:"error" msg:"error"`
	IpAddr4   net.IP      `json:"ipAddr4" form:"ipAddr4" query:"ipAddr4" long:"ipAddr4" msg:"ipAddr4"`
	IpAddr6   net.IP      `json:"ipAddr6" form:"ipAddr6" query:"ipAddr6" long:"ipAddr6" msg:"ipAddr6"`
	UserAgent string      `json:"userAgent" form:"userAgent" query:"userAgent" long:"userAgent" msg:"userAgent"`
}

func NewUserLogs(adapter *Ch.Adapter) *UserLogs {
	return &UserLogs{Adapter: adapter}
}

// UserLogsFieldTypeMap returns key value of field name and key
var UserLogsFieldTypeMap = map[string]Ch.DataType{ //nolint:dupl false positive
	`createdAt`: Ch.DateTime,
	`requestId`: Ch.UInt64,
	`actorId`:   Ch.UInt64,
	`error`:     Ch.String,
	`ipAddr4`:   Ch.IPv4,
	`ipAddr6`:   Ch.IPv6,
	`userAgent`: Ch.String,
}

func (u *UserLogs) TableName() Ch.TableName { //nolint:dupl false positive
	return mAuth.TableUserLogs
}

func (u *UserLogs) SqlTableName() string { //nolint:dupl false positive
	return `"userLogs"`
}

func (u *UserLogs) ScanRowAllCols(rows *sql.Rows) (err error) { //nolint:dupl false positive
	return rows.Scan(
		&u.CreatedAt,
		&u.RequestId,
		&u.ActorId,
		&u.Error,
		&u.IpAddr4,
		&u.IpAddr6,
		&u.UserAgent,
	)
}

func (u *UserLogs) ScanRowsAllCols(rows *sql.Rows, estimateRows int) (res []UserLogs, err error) { //nolint:dupl false positive
	res = make([]UserLogs, 0, estimateRows)
	defer rows.Close()
	for rows.Next() {
		var row UserLogs
		err = row.ScanRowAllCols(rows)
		if err != nil {
			return
		}
		res = append(res, row)
	}
	return
}

// insert, error if exists
func (u *UserLogs) SqlInsert() string { //nolint:dupl false positive
	return `INSERT INTO ` + u.SqlTableName() + `(` + u.SqlAllFields() + `) VALUES (?,?,?,?,?,?,?)`
}

func (u *UserLogs) SqlCount() string { //nolint:dupl false positive
	return `SELECT COUNT(*) FROM ` + u.SqlTableName()
}

func (u *UserLogs) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` createdAt
	, requestId
	, actorId
	, error
	, ipAddr4
	, ipAddr6
	, userAgent
	`
}

func (u *UserLogs) SqlAllFields() string { //nolint:dupl false positive
	return `createdAt, requestId, actorId, error, ipAddr4, ipAddr6, userAgent`
}

func (u UserLogs) SqlInsertParam() []any { //nolint:dupl false positive
	return []any{
		u.CreatedAt, // 0
		u.RequestId, // 1
		u.ActorId,   // 2
		u.Error,     // 3
		u.IpAddr4,   // 4
		u.IpAddr6,   // 5
		u.UserAgent, // 6
	}
}

func (u *UserLogs) IdxCreatedAt() int { //nolint:dupl false positive
	return 0
}

func (u *UserLogs) SqlCreatedAt() string { //nolint:dupl false positive
	return `createdAt`
}

func (u *UserLogs) IdxRequestId() int { //nolint:dupl false positive
	return 1
}

func (u *UserLogs) SqlRequestId() string { //nolint:dupl false positive
	return `requestId`
}

func (u *UserLogs) IdxActorId() int { //nolint:dupl false positive
	return 2
}

func (u *UserLogs) SqlActorId() string { //nolint:dupl false positive
	return `actorId`
}

func (u *UserLogs) IdxError() int { //nolint:dupl false positive
	return 3
}

func (u *UserLogs) SqlError() string { //nolint:dupl false positive
	return `error`
}

func (u *UserLogs) IdxIpAddr4() int { //nolint:dupl false positive
	return 4
}

func (u *UserLogs) SqlIpAddr4() string { //nolint:dupl false positive
	return `ipAddr4`
}

func (u *UserLogs) IdxIpAddr6() int { //nolint:dupl false positive
	return 5
}

func (u *UserLogs) SqlIpAddr6() string { //nolint:dupl false positive
	return `ipAddr6`
}

func (u *UserLogs) IdxUserAgent() int { //nolint:dupl false positive
	return 6
}

func (u *UserLogs) SqlUserAgent() string { //nolint:dupl false positive
	return `userAgent`
}

func (u *UserLogs) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		u.CreatedAt, // 0
		u.RequestId, // 1
		u.ActorId,   // 2
		u.Error,     // 3
		u.IpAddr4,   // 4
		u.IpAddr6,   // 5
		u.UserAgent, // 6
	}
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go
