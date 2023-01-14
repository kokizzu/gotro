package rqZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`github.com/kokizzu/gotro/D/Tt/mZzz`

	`github.com/tarantool/go-tarantool`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqZzz__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type rqZzz__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type rqZzz__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type rqZzz__ORM.GEN.go
// go:generate msgp -tests=false -file rqZzz__ORM.GEN.go -o rqZzz__MSG.GEN.go

// Zzz DAO reader/query struct
type Zzz struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id        uint64
	CreatedAt int64
}

// NewZzz create new ORM reader/query object
func NewZzz(adapter *Tt.Adapter) *Zzz {
	return &Zzz{Adapter: adapter}
}

// sqlTableName returns full package and table name
func (z *Zzz) SpaceName() string { //nolint:dupl false positive
	return string(mZzz.TableZzz)
}

// sqlTableName returns quoted table name
func (z *Zzz) sqlTableName() string { //nolint:dupl false positive
	return `"zzz"`
}

func (z *Zzz) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (z *Zzz) FindById() bool { //nolint:dupl false positive
	res, err := z.Adapter.Select(z.SpaceName(), z.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{z.Id})
	if L.IsError(err, `Zzz.FindById failed: `+z.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		z.FromArray(rows[0])
		return true
	}
	return false
}

// sqlSelectAllFields generate sql select fields
func (z *Zzz) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "created_at"
	`
}

// ToUpdateArray generate slice of update command
func (z *Zzz) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, z.Id},
		A.X{`=`, 1, z.CreatedAt},
	}
}

// IdxId return name of the index
func (z *Zzz) IdxId() int { //nolint:dupl false positive
	return 0
}

// sqlIdxId return name of the column being indexed
func (z *Zzz) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxCreatedAt return name of the index
func (z *Zzz) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

// sqlIdxCreatedAt return name of the column being indexed
func (z *Zzz) sqlCreatedAt() string { //nolint:dupl false positive
	return `"created_at"`
}

// ToArray receiver fields to slice
func (z *Zzz) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if z.Id != 0 {
		id = z.Id
	}
	return A.X{
		id,
		z.CreatedAt, // 1
	}
}

// FromArray convert slice to receiver fields
func (z *Zzz) FromArray(a A.X) *Zzz { //nolint:dupl false positive
	z.Id = X.ToU(a[0])
	z.CreatedAt = X.ToI(a[1])
	return z
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (z *Zzz) FindOffsetLimit(offset, limit uint32, idx string) []Zzz { //nolint:dupl false positive
	var rows []Zzz
	res, err := z.Adapter.Select(z.SpaceName(), idx, offset, limit, 2, A.X{})
	if L.IsError(err, `Zzz.FindOffsetLimit failed: `+z.SpaceName()) {
		return rows
	}
	for _, row := range res.Tuples() {
		item := Zzz{}
		rows = append(rows, *item.FromArray(row))
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (z *Zzz) FindArrOffsetLimit(offset, limit uint32, idx string) []A.X { //nolint:dupl false positive
	var rows []A.X
	res, err := z.Adapter.Select(z.SpaceName(), idx, offset, limit, 2, A.X{})
	if L.IsError(err, `Zzz.FindOffsetLimit failed: `+z.SpaceName()) {
		return rows
	}
	tuples := res.Tuples()
	rows = make([]A.X, len(tuples))
	for z, row := range tuples {
		rows[z] = row
	}
	return rows
}

// Total count number of rows
func (z *Zzz) Total() int64 { //nolint:dupl false positive
	rows := z.Adapter.CallBoxSpace(z.SpaceName() + `:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

