package Sc

import (
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/X"
	"time"
)

// base structure for all model, including schema
type Base struct {
	Table      string
	Connection *RDBMS
	Id         string    `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	DeletedAt  time.Time `db:"deleted_at"`
	RestoredAt time.Time `db:"restored_at"`
	ModifiedAt time.Time `db:"modified_at"`
	CreatedBy  string    `db:"created_by"`
	UpdatedBy  string    `db:"updated_by"`
	DeletedBy  string    `db:"deleted_by"`
	RestoredBy string    `db:"restored_by"`
	IsDeleted  bool      `db:"is_deleted"`
	XData      M.SX
}

func (b *Base) FromMap(m M.SX) {
	b.Id = X.ToS(m[`id`])
	b.CreatedAt = time.Unix(0, X.ToI(m[`created_at`]))
	b.UpdatedAt = time.Unix(0, X.ToI(m[`updated_at`]))
	b.DeletedAt = time.Unix(0, X.ToI(m[`deleted_at`]))
	b.RestoredAt = time.Unix(0, X.ToI(m[`restored_at`]))
	b.ModifiedAt = time.Unix(0, X.ToI(m[`modified_at`]))
	b.CreatedBy = X.ToS(m[`created_by`])
	b.UpdatedBy = X.ToS(m[`updated_by`])
	b.DeletedBy = X.ToS(m[`deleted_by`])
	b.RestoredBy = X.ToS(m[`restored_by`])
	b.IsDeleted = X.ToBool(m[`is_deleted`])
	b.XData = m
}

func (b *Base) ToMap() M.SX {
	b.XData[`id`] = b.Id
	b.XData[`created_at`] = b.CreatedAt.UnixNano()
	b.XData[`updated_at`] = b.UpdatedAt.UnixNano()
	b.XData[`deleted_at`] = b.DeletedAt.UnixNano()
	b.XData[`restored_at`] = b.RestoredAt.UnixNano()
	b.XData[`modified_at`] = b.ModifiedAt.UnixNano()
	b.XData[`created_by`] = b.CreatedBy
	b.XData[`updated_by`] = b.UpdatedBy
	b.XData[`deleted_by`] = b.DeletedBy
	b.XData[`restored_by`] = b.RestoredBy
	b.XData[`is_deleted`] = b.IsDeleted
	return b.XData
}

// 2015-08-27 also update unique_id
func (b *Base) Save(actor string) bool {
	records := b.ToMap()
	return b.Connection.DoUpsert(actor, b.Table, records)
}

// 2016-05-25 delete row, TODO: add DeletedAt
func (b *Base) Delete(actor string) bool {
	if b.Connection.DoDelete(actor, b.Table, b.Id) {
		b.IsDeleted = true
		return true
	}
	return false
}

// 2016-05-25 restore row, TODO: add RestoredAt
func (b *Base) Restore(actor string) bool {
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

// get []interface{} from Data
func (b *Base) GetArr(key string) []interface{} {
	return X.ToArr(b.XData[key])
}

// get float64 from Data
func (b *Base) GetFloat(key string) float64 {
	return b.XData.GetFloat(key)
}

// get id
func (b *Base) GetId() string {
	return b.Id
}

// get unique id
func (b *Base) GetUniqueId() string {
	return b.Id
}

func (b *Base) SetVal(key string, val interface{}) {
	b.XData[key] = val
}
