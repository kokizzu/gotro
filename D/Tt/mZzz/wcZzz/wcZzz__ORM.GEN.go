package wcZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`github.com/kokizzu/gotro/D/Tt/mZzz/rqZzz`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/M`
	`github.com/kokizzu/gotro/S`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcZzz__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type wcZzz__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type wcZzz__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type wcZzz__ORM.GEN.go
// ZzzMutator DAO writer/command struct
type ZzzMutator struct {
	rqZzz.Zzz
	mutations *tarantool.Operations
	logs	  []A.X
}

// NewZzzMutator create new ORM writer/command object
func NewZzzMutator(adapter *Tt.Adapter) (res *ZzzMutator) {
	res = &ZzzMutator{Zzz: rqZzz.Zzz{Adapter: adapter}}
	res.mutations = tarantool.NewOperations()
	res.Coords = []any{}
	return
}

// Logs get array of logs [field, old, new]
func (z *ZzzMutator) Logs() []A.X { //nolint:dupl false positive
	return z.logs
}

// HaveMutation check whether Set* methods ever called
func (z *ZzzMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(z.logs) > 0
}

// ClearMutations clear all previously called Set* methods
func (z *ZzzMutator) ClearMutations() { //nolint:dupl false positive
	z.mutations = tarantool.NewOperations()
	z.logs = []A.X{}
}

// DoOverwriteById update all columns, error if not exists, not using mutations/Set*
func (z *ZzzMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := z.Adapter.RetryDo(tarantool.NewUpdateRequest(z.SpaceName()).
		Index(z.UniqueIndexId()).
		Key(tarantool.UintKey{I:uint(z.Id)}).
		Operations(z.ToUpdateArray()),
	)
	return !L.IsError(err, `Zzz.DoOverwriteById failed: `+z.SpaceName())
}

// DoUpdateById update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (z *ZzzMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !z.HaveMutation() {
		return true
	}
	_, err := z.Adapter.RetryDo(
		tarantool.NewUpdateRequest(z.SpaceName()).
		Index(z.UniqueIndexId()).
		Key(tarantool.UintKey{I:uint(z.Id)}).
		Operations(z.mutations),
	)
	return !L.IsError(err, `Zzz.DoUpdateById failed: `+z.SpaceName())
}

// DoDeletePermanentById permanent delete
func (z *ZzzMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := z.Adapter.RetryDo(
		tarantool.NewDeleteRequest(z.SpaceName()).
		Index(z.UniqueIndexId()).
		Key(tarantool.UintKey{I:uint(z.Id)}),
	)
	return !L.IsError(err, `Zzz.DoDeletePermanentById failed: `+z.SpaceName())
}

// DoOverwriteByName update all columns, error if not exists, not using mutations/Set*
func (z *ZzzMutator) DoOverwriteByName() bool { //nolint:dupl false positive
	_, err := z.Adapter.RetryDo(tarantool.NewUpdateRequest(z.SpaceName()).
		Index(z.UniqueIndexName()).
		Key(tarantool.StringKey{S:z.Name}).
		Operations(z.ToUpdateArray()),
	)
	return !L.IsError(err, `Zzz.DoOverwriteByName failed: `+z.SpaceName())
}

// DoUpdateByName update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (z *ZzzMutator) DoUpdateByName() bool { //nolint:dupl false positive
	if !z.HaveMutation() {
		return true
	}
	_, err := z.Adapter.RetryDo(
		tarantool.NewUpdateRequest(z.SpaceName()).
		Index(z.UniqueIndexName()).
		Key(tarantool.StringKey{S:z.Name}).
		Operations(z.mutations),
	)
	return !L.IsError(err, `Zzz.DoUpdateByName failed: `+z.SpaceName())
}

// DoDeletePermanentByName permanent delete
func (z *ZzzMutator) DoDeletePermanentByName() bool { //nolint:dupl false positive
	_, err := z.Adapter.RetryDo(
		tarantool.NewDeleteRequest(z.SpaceName()).
		Index(z.UniqueIndexName()).
		Key(tarantool.StringKey{S:z.Name}),
	)
	return !L.IsError(err, `Zzz.DoDeletePermanentByName failed: `+z.SpaceName())
}

// DoInsert insert, error if already exists
func (z *ZzzMutator) DoInsert() bool { //nolint:dupl false positive
	arr := z.ToArray()
	row, err := z.Adapter.RetryDo(
		tarantool.NewInsertRequest(z.SpaceName()).
		Tuple(arr),
	)
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				z.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `Zzz.DoInsert failed: `+z.SpaceName() + `\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// tarantool's replace/upsert can only match by primary key
// previous name: DoReplace
func (z *ZzzMutator) DoUpsertById() bool { //nolint:dupl false positive
	if z.Id > 0 {
		return z.DoUpdateById()
	}
	return z.DoInsert()
}

// SetId create mutations, should not duplicate
func (z *ZzzMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != z.Id {
		z.mutations.Assign(0, val)
		z.logs = append(z.logs, A.X{`id`, z.Id, val})
		z.Id = val
		return true
	}
	return false
}

// SetCreatedAt create mutations, should not duplicate
func (z *ZzzMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != z.CreatedAt {
		z.mutations.Assign(1, val)
		z.logs = append(z.logs, A.X{`createdAt`, z.CreatedAt, val})
		z.CreatedAt = val
		return true
	}
	return false
}

// SetCoords create mutations, should not duplicate
func (z *ZzzMutator) SetCoords(val []any) bool { //nolint:dupl false positive
	z.mutations.Assign(2, val)
	z.logs = append(z.logs, A.X{`coords`, z.Coords, val})
	z.Coords = val
	return true
}

// SetName create mutations, should not duplicate
func (z *ZzzMutator) SetName(val string) bool { //nolint:dupl false positive
	if val != z.Name {
		z.mutations.Assign(3, val)
		z.logs = append(z.logs, A.X{`name`, z.Name, val})
		z.Name = val
		return true
	}
	return false
}

// SetHeightMeter create mutations, should not duplicate
func (z *ZzzMutator) SetHeightMeter(val float64) bool { //nolint:dupl false positive
	if val != z.HeightMeter {
		z.mutations.Assign(4, val)
		z.logs = append(z.logs, A.X{`heightMeter`, z.HeightMeter, val})
		z.HeightMeter = val
		return true
	}
	return false
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (z *ZzzMutator) SetAll(from rqZzz.Zzz, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`id`] && (forceMap[`id`] || from.Id != 0) {
		z.Id = from.Id
		changed = true
	}
	if !excludeMap[`createdAt`] && (forceMap[`createdAt`] || from.CreatedAt != 0) {
		z.CreatedAt = from.CreatedAt
		changed = true
	}
	if !excludeMap[`coords`] && (forceMap[`coords`] || from.Coords != nil) {
		z.Coords = from.Coords
		changed = true
	}
	if !excludeMap[`name`] && (forceMap[`name`] || from.Name != ``) {
		z.Name = S.Trim(from.Name)
		changed = true
	}
	if !excludeMap[`heightMeter`] && (forceMap[`heightMeter`] || from.HeightMeter != 0) {
		z.HeightMeter = from.HeightMeter
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

