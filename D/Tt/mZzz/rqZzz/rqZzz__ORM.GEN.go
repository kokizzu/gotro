package rqZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`github.com/kokizzu/gotro/D/Tt/mZzz`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqZzz__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type rqZzz__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type rqZzz__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type rqZzz__ORM.GEN.go
// Zzz DAO reader/query struct
type Zzz struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id          uint64
	CreatedAt   int64
	Coords      []any
	Name        string
	HeightMeter float64
}

// NewZzz create new ORM reader/query object
func NewZzz(adapter *Tt.Adapter) *Zzz {
	return &Zzz{Adapter: adapter}
}

// SpaceName returns full package and table name
func (z *Zzz) SpaceName() string { //nolint:dupl false positive
	return string(mZzz.TableZzz) // casting required to string from Tt.TableName
}

// SqlTableName returns quoted table name
func (z *Zzz) SqlTableName() string { //nolint:dupl false positive
	return `"zzz"`
}

func (z *Zzz) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (z *Zzz) FindById() bool { //nolint:dupl false positive
	res, err := z.Adapter.RetryDo(
		tarantool.NewSelectRequest(z.SpaceName()).
		Index(z.UniqueIndexId()).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(tarantool.UintKey{I:uint(z.Id)}),
	)
	if L.IsError(err, `Zzz.FindById failed: `+z.SpaceName()) {
		return false
	}
	if len(res) == 1 {
		if row, ok := res[0].([]any); ok {
			z.FromArray(row)
			return true
		}
	}
	return false
}

// SpatialIndexCoords return spatial index name
func (z *Zzz) SpatialIndexCoords() string { //nolint:dupl false positive
	return `coords`
}

// UniqueIndexName return unique index name
func (z *Zzz) UniqueIndexName() string { //nolint:dupl false positive
	return `name`
}

// FindByName Find one by Name
func (z *Zzz) FindByName() bool { //nolint:dupl false positive
	res, err := z.Adapter.RetryDo(
		tarantool.NewSelectRequest(z.SpaceName()).
		Index(z.UniqueIndexName()).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(tarantool.StringKey{S:z.Name}),
	)
	if L.IsError(err, `Zzz.FindByName failed: `+z.SpaceName()) {
		return false
	}
	if len(res) == 1 {
		if row, ok := res[0].([]any); ok {
			z.FromArray(row)
			return true
		}
	}
	return false
}

// SqlSelectAllFields generate Sql select fields
func (z *Zzz) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "createdAt"
	, "coords"
	, "name"
	, "heightMeter"
	`
}

// SqlSelectAllUncensoredFields generate Sql select fields
func (z *Zzz) SqlSelectAllUncensoredFields() string { //nolint:dupl false positive
	return ` "id"
	, "createdAt"
	, "coords"
	, "name"
	, "heightMeter"
	`
}

// ToUpdateArray generate slice of update command
func (z *Zzz) ToUpdateArray() *tarantool.Operations { //nolint:dupl false positive
	return tarantool.NewOperations().
		Assign(0, z.Id).
		Assign(1, z.CreatedAt).
		Assign(2, z.Coords).
		Assign(3, z.Name).
		Assign(4, z.HeightMeter)
}

// IdxId return name of the index
func (z *Zzz) IdxId() int { //nolint:dupl false positive
	return 0
}

// SqlId return name of the column being indexed
func (z *Zzz) SqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxCreatedAt return name of the index
func (z *Zzz) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

// SqlCreatedAt return name of the column being indexed
func (z *Zzz) SqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

// IdxCoords return name of the index
func (z *Zzz) IdxCoords() int { //nolint:dupl false positive
	return 2
}

// SqlCoords return name of the column being indexed
func (z *Zzz) SqlCoords() string { //nolint:dupl false positive
	return `"coords"`
}

// IdxName return name of the index
func (z *Zzz) IdxName() int { //nolint:dupl false positive
	return 3
}

// SqlName return name of the column being indexed
func (z *Zzz) SqlName() string { //nolint:dupl false positive
	return `"name"`
}

// IdxHeightMeter return name of the index
func (z *Zzz) IdxHeightMeter() int { //nolint:dupl false positive
	return 4
}

// SqlHeightMeter return name of the column being indexed
func (z *Zzz) SqlHeightMeter() string { //nolint:dupl false positive
	return `"heightMeter"`
}

// ToArray receiver fields to slice
func (z *Zzz) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if z.Id != 0 {
		id = z.Id
	}
	return A.X{
		id,
		z.CreatedAt,   // 1
		z.Coords,      // 2
		z.Name,        // 3
		z.HeightMeter, // 4
	}
}

// FromArray convert slice to receiver fields
func (z *Zzz) FromArray(a A.X) *Zzz { //nolint:dupl false positive
	z.Id = X.ToU(a[0])
	z.CreatedAt = X.ToI(a[1])
	z.Coords = X.ToArr(a[2])
	z.Name = X.ToS(a[3])
	z.HeightMeter = X.ToF(a[4])
	return z
}

// FromUncensoredArray convert slice to receiver fields
func (z *Zzz) FromUncensoredArray(a A.X) *Zzz { //nolint:dupl false positive
	z.Id = X.ToU(a[0])
	z.CreatedAt = X.ToI(a[1])
	z.Coords = X.ToArr(a[2])
	z.Name = X.ToS(a[3])
	z.HeightMeter = X.ToF(a[4])
	return z
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (z *Zzz) FindOffsetLimit(offset, limit uint32, idx string) []Zzz { //nolint:dupl false positive
	var rows []Zzz
	res, err := z.Adapter.RetryDo(
		tarantool.NewSelectRequest(z.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll),
	)
	if L.IsError(err, `Zzz.FindOffsetLimit failed: `+z.SpaceName()) {
		return rows
	}
	for _, row := range res {
		item := Zzz{}
		row, ok := row.([]any)
		if ok {
			rows = append(rows, *item.FromArray(row))
		}
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (z *Zzz) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	resp, err := z.Adapter.RetryDoResp(
		tarantool.NewSelectRequest(z.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll),
	)
	if L.IsError(err, `Zzz.FindOffsetLimit failed: `+z.SpaceName()) {
		return rows, Tt.QueryMetaFrom(resp, err)
	}
	res, err := resp.Decode()
	if L.IsError(err, `Zzz.FindOffsetLimit failed: `+z.SpaceName()) {
		return rows, Tt.QueryMetaFrom(resp, err)
	}
	rows = make([]A.X, len(res))
	for _, row := range res {
		row, ok := row.([]any)
		if ok {
			rows = append(rows, row)
		}
	}
	return rows, Tt.QueryMetaFrom(resp, nil)
}

// Total count number of rows
func (z *Zzz) Total() int64 { //nolint:dupl false positive
	rows := z.Adapter.CallBoxSpace(z.SpaceName() + `:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// ZzzFieldTypeMap returns key value of field name and key
var ZzzFieldTypeMap = map[string]Tt.DataType { //nolint:dupl false positive
	`id`:          Tt.Unsigned,
	`createdAt`:   Tt.Integer,
	`coords`:      Tt.Array,
	`name`:        Tt.String,
	`heightMeter`: Tt.Double,
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

