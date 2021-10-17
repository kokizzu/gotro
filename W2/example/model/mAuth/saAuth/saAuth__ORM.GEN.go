package saAuth

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go

import (
	"database/sql"
	"github.com/kokizzu/gotro/W2/example/model/mAuth"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	chBuffer "github.com/kokizzu/ch-timed-buffer"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/L"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file saAuth__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type saAuth__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type saAuth__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type saAuth__ORM.GEN.go
// go:generate msgp -tests=false -file saAuth__ORM.GEN.go -o saAuth__MSG.GEN.go

var userLogsDummy = UserLogs{}
var Preparators = map[Ch.TableName]chBuffer.Preparator{
	mAuth.TableUserLogs: func(tx *sql.Tx) *sql.Stmt {
		query := userLogsDummy.sqlInsert()
		stmt, err := tx.Prepare(query)
		L.IsError(err, `failed to tx.Prepare: `+query)
		return stmt
	},
}

type UserLogs struct {
	Adapter   *Ch.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	CreatedAt time.Time
	RequestId uint64
	ActorId   uint64
	Error     string
	IpAddr4   string
	IpAddr6   string
	UserAgent string
}

func NewUserLogs(adapter *Ch.Adapter) *UserLogs {
	return &UserLogs{Adapter: adapter}
}

func (u UserLogs) TableName() Ch.TableName { //nolint:dupl false positive
	return mAuth.TableUserLogs
}

func (u *UserLogs) sqlTableName() string { //nolint:dupl false positive
	return `"userLogs"`
}

// insert, error if exists
func (u *UserLogs) sqlInsert() string { //nolint:dupl false positive
	return `INSERT INTO ` + u.sqlTableName() + `(` + u.sqlAllFields() + `) VALUES (?,?,?,?,?,?,?)`
}

func (u *UserLogs) sqlCount() string { //nolint:dupl false positive
	return `SELECT COUNT(*) FROM ` + u.sqlTableName()
}

func (u *UserLogs) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` createdAt
	, requestId
	, actorId
	, error
	, ipAddr4
	, ipAddr6
	, userAgent
	`
}

func (u *UserLogs) sqlAllFields() string { //nolint:dupl false positive
	return `createdAt, requestId, actorId, error, ipAddr4, ipAddr6, userAgent`
}

func (u UserLogs) SqlInsertParam() []interface{} { //nolint:dupl false positive
	return []interface{}{
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

func (u *UserLogs) sqlCreatedAt() string { //nolint:dupl false positive
	return `createdAt`
}

func (u *UserLogs) IdxRequestId() int { //nolint:dupl false positive
	return 1
}

func (u *UserLogs) sqlRequestId() string { //nolint:dupl false positive
	return `requestId`
}

func (u *UserLogs) IdxActorId() int { //nolint:dupl false positive
	return 2
}

func (u *UserLogs) sqlActorId() string { //nolint:dupl false positive
	return `actorId`
}

func (u *UserLogs) IdxError() int { //nolint:dupl false positive
	return 3
}

func (u *UserLogs) sqlError() string { //nolint:dupl false positive
	return `error`
}

func (u *UserLogs) IdxIpAddr4() int { //nolint:dupl false positive
	return 4
}

func (u *UserLogs) sqlIpAddr4() string { //nolint:dupl false positive
	return `ipAddr4`
}

func (u *UserLogs) IdxIpAddr6() int { //nolint:dupl false positive
	return 5
}

func (u *UserLogs) sqlIpAddr6() string { //nolint:dupl false positive
	return `ipAddr6`
}

func (u *UserLogs) IdxUserAgent() int { //nolint:dupl false positive
	return 6
}

func (u *UserLogs) sqlUserAgent() string { //nolint:dupl false positive
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
