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
//go:generate replacer -afterprefix 'Id" form' 'Id,string" form' type rqZzz__ORM.GEN.go
//go:generate replacer -afterprefix 'json:"id"' 'json:"id,string"' type rqZzz__ORM.GEN.go
//go:generate replacer -afterprefix 'By" form' 'By,string" form' type rqZzz__ORM.GEN.go
// go:generate msgp -tests=false -file rqZzz__ORM.GEN.go -o rqZzz__MSG.GEN.go

// Zzz DAO reader/query struct
type Zzz struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id        uint64
	CreatedAt int64
	Coords    []any
}

// NewZzz create new ORM reader/query object
func NewZzz(adapter *Tt.Adapter) *Zzz {
	return &Zzz{Adapter: adapter}
}

// SpaceName returns full package and table name
func (z *Zzz) SpaceName() string { //nolint:dupl false positive
	return string(mZzz.TableZzz) // casting required to string from Tt.TableName
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

// SpatialIndexCoords return spatial index name
func (z *Zzz) SpatialIndexCoords() string { //nolint:dupl false positive
	return `coords`
}

// sqlSelectAllFields generate sql select fields
func (z *Zzz) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "created_at"
	, "coords"
	`
}

// ToUpdateArray generate slice of update command
func (z *Zzz) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, z.Id},
		A.X{`=`, 1, z.CreatedAt},
		A.X{`=`, 2, z.Coords},
	}
}

// IdxId return name of the index
func (z *Zzz) IdxId() int { //nolint:dupl false positive
	return 0
}

// sqlId return name of the column being indexed
func (z *Zzz) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxCreatedAt return name of the index
func (z *Zzz) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

// sqlCreatedAt return name of the column being indexed
func (z *Zzz) sqlCreatedAt() string { //nolint:dupl false positive
	return `"created_at"`
}

// IdxCoords return name of the index
func (z *Zzz) IdxCoords() int { //nolint:dupl false positive
	return 2
}

// sqlCoords return name of the column being indexed
func (z *Zzz) sqlCoords() string { //nolint:dupl false positive
	return `"coords"`
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
		z.Coords,    // 2
	}
}

// FromArray convert slice to receiver fields
func (z *Zzz) FromArray(a A.X) *Zzz { //nolint:dupl false positive
	z.Id = X.ToU(a[0])
	z.CreatedAt = X.ToI(a[1])
	z.Coords = X.ToArr(a[2])
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
func (z *Zzz) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	res, err := z.Adapter.Select(z.SpaceName(), idx, offset, limit, 2, A.X{})
	if L.IsError(err, `Zzz.FindOffsetLimit failed: `+z.SpaceName()) {
		return rows, Tt.QueryMetaFrom(res, err)
	}
	tuples := res.Tuples()
	rows = make([]A.X, len(tuples))
	for z, row := range tuples {
		rows[z] = row
	}
	return rows, Tt.QueryMetaFrom(res, nil)
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

