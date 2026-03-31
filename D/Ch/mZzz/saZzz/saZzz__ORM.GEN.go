package saZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go

import (
	`database/sql`
	`time`

	`github.com/kokizzu/gotro/D/Ch/mZzz`

	_ `github.com/ClickHouse/clickhouse-go/v2`
	chBuffer `github.com/kokizzu/ch-timed-buffer`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Ch`
	`github.com/kokizzu/gotro/L`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file saZzz__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type saZzz__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type saZzz__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type saZzz__ORM.GEN.go
// go:generate msgp -tests=false -file saZzz__ORM.GEN.go -o saZzz__MSG.GEN.go

var zzzDummy = Zzz{}
var Preparators = map[Ch.TableName]chBuffer.Preparator{
	mZzz.TableZzz: func(tx *sql.Tx) *sql.Stmt {
		query := zzzDummy.SqlInsert()
		stmt, err := tx.Prepare(query)
		L.IsError(err, `failed to tx.Prepare: `+query)
		return stmt
	},
}
type Zzz struct {
	Adapter *Ch.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id          uint64
	CreatedAt   time.Time
	Name        string
	HeightMeter float64
}

func NewZzz(adapter *Ch.Adapter) *Zzz {
	return &Zzz{Adapter: adapter}
}

// ZzzFieldTypeMap returns key value of field name and key
var ZzzFieldTypeMap = map[string]Ch.DataType { //nolint:dupl false positive
	`id`:          Ch.UInt64,
	`createdAt`:   Ch.DateTime,
	`name`:        Ch.String,
	`heightMeter`: Ch.Float64,
}

func (z *Zzz) TableName() Ch.TableName { //nolint:dupl false positive
	return mZzz.TableZzz
}

func (z *Zzz) SqlTableName() string { //nolint:dupl false positive
	return `"zzz"`
}

func (z *Zzz) ScanRowAllCols(rows *sql.Rows) (err error) { //nolint:dupl false positive
	return rows.Scan(
		&z.Id,
		&z.CreatedAt,
		&z.Name,
		&z.HeightMeter,
	)
}

func (z *Zzz) ScanRowsAllCols(rows *sql.Rows, estimateRows int) (res []Zzz, err error) { //nolint:dupl false positive
	res = make([]Zzz, 0, estimateRows)
	defer rows.Close()
	for rows.Next() {
		var row Zzz
		err = row.ScanRowAllCols(rows)
		if err != nil {
			return
		}
		res = append(res, row)
	}
	return
}

// insert, error if exists
func (z *Zzz) SqlInsert() string { //nolint:dupl false positive
	return `INSERT INTO ` + z.SqlTableName() + ` (` + z.SqlAllFields() + `) VALUES (?,?,?,?)`
}

func (z *Zzz) SqlCount() string { //nolint:dupl false positive
	return `SELECT COUNT(*) FROM ` + z.SqlTableName()
}

func (z *Zzz) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` id
	, createdAt
	, name
	, heightMeter
	`
}

func (z *Zzz) SqlAllFields() string { //nolint:dupl false positive
	return `id, createdAt, name, heightMeter`
}

func (z Zzz) SqlInsertParam() []any { //nolint:dupl false positive
	return []any{
		z.Id, // 0 
		z.CreatedAt, // 1 
		z.Name, // 2 
		z.HeightMeter, // 3 
	}
}

func (z *Zzz) IdxId() int { //nolint:dupl false positive
	return 0
}

func (z *Zzz) SqlId() string { //nolint:dupl false positive
	return `id`
}

func (z *Zzz) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

func (z *Zzz) SqlCreatedAt() string { //nolint:dupl false positive
	return `createdAt`
}

func (z *Zzz) IdxName() int { //nolint:dupl false positive
	return 2
}

func (z *Zzz) SqlName() string { //nolint:dupl false positive
	return `name`
}

func (z *Zzz) IdxHeightMeter() int { //nolint:dupl false positive
	return 3
}

func (z *Zzz) SqlHeightMeter() string { //nolint:dupl false positive
	return `heightMeter`
}

func (z *Zzz) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		z.Id,          // 0
		z.CreatedAt,   // 1
		z.Name,        // 2
		z.HeightMeter, // 3
	}
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go

