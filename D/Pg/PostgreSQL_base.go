package Pg

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

// base structure for all model, including schema
type Base struct {
	Id         int64 `db:"id"`
	Table      string
	Connection *RDBMS
	UniqueId   sql.NullString `db:"unique_id"`
	CreatedAt  pq.NullTime    `db:"created_at"`
	UpdatedAt  pq.NullTime    `db:"updated_at"`
	DeletedAt  pq.NullTime    `db:"deleted_at"`
	RestoredAt pq.NullTime    `db:"restored_at"`
	ModifiedAt pq.NullTime    `db:"modified_at"`
	CreatedBy  sql.NullInt64  `db:"created_by"`
	UpdatedBy  sql.NullInt64  `db:"updated_by"`
	DeletedBy  sql.NullInt64  `db:"deleted_by"`
	RestoredBy sql.NullInt64  `db:"restored_by"`
	IsDeleted  bool           `db:"is_deleted"`
	DataStr    string         `db:"data"` // json object
	XData      M.SX
}

func (b *Base) DataToMap() M.SX {
	if b.DataStr != `` {
		b.XData = S.JsonToMap(b.DataStr)
	}
	return b.XData
}

func (b *Base) MapToData() string {
	b.DataStr = M.ToJson(b.XData)
	return b.DataStr
}

// 2015-08-27 also update unique_id
func (b *Base) Save(actor int64) bool {
	records := M.SX{
		`id`:   b.Id,
		`data`: b.MapToData(),
	}
	if b.UniqueId.Valid {
		records[`unique_id`] = b.UniqueId.String
	}
	id := b.Connection.DoUpsert(actor, b.Table, records)
	if id == 0 {
		return false
	}
	b.Id = id
	return true
}

// 2016-05-25 delete row, TODO: add DeletedAt
func (b *Base) Delete(actor int64) bool {
	if b.Connection.DoDelete(actor, b.Table, b.Id) {
		b.IsDeleted = true
		return true
	}
	return false
}

// 2016-05-25 restore row, TODO: add RestoredAt
func (b *Base) Restore(actor int64) bool {
	if b.Connection.DoRestore(actor, b.Table, b.Id) {
		b.IsDeleted = false
		return true
	}
	return false
}

// get string from Data
func (b *Base) GetMSX(key string) M.SX {
	return b.XData.GetMSX(key)
}

// get string from Data
func (b *Base) GetStr(key string) string {
	return b.XData.GetStr(key)
}

// get boolean from Data
func (b *Base) GetBool(key string) bool {
	return b.XData.GetBool(key)
}

// get int64 from Data
func (b *Base) GetInt(key string) int64 {
	return b.XData.GetInt(key)
}

// get []any from Data
func (b *Base) GetArr(key string) []any {
	return X.ToArr(b.XData[key])
}

// get float64 from Data
func (b *Base) GetFloat(key string) float64 {
	return b.XData.GetFloat(key)
}

// get id
func (b *Base) GetId() int64 {
	return b.Id
}

// get unique id
func (b *Base) GetUniqueId() string {
	return b.UniqueId.String
}

func (b *Base) SetVal(key string, val any) {
	b.XData[key] = val
}
