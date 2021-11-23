package wcZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`github.com/kokizzu/gotro/D/Tt/mZzz/rqZzz`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcZzz__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type wcZzz__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type wcZzz__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type wcZzz__ORM.GEN.go
// go:generate msgp -tests=false -file wcZzz__ORM.GEN.go -o wcZzz__MSG.GEN.go

type ZzzMutator struct {
	rqZzz.Zzz
	mutations []A.X
}

func NewZzzMutator(adapter *Tt.Adapter) *ZzzMutator {
	return &ZzzMutator{Zzz: rqZzz.Zzz{Adapter: adapter}}
}

func (z *ZzzMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(z.mutations) > 0
}

// Overwrite all columns, error if not exists
func (z *ZzzMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := z.Adapter.Update(z.SpaceName(), z.UniqueIndexId(), A.X{z.Id}, z.ToUpdateArray())
	return !L.IsError(err, `Zzz.DoOverwriteById failed: `+z.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (z *ZzzMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !z.HaveMutation() {
		return true
	}
	_, err := z.Adapter.Update(z.SpaceName(), z.UniqueIndexId(), A.X{z.Id}, z.mutations)
	return !L.IsError(err, `Zzz.DoUpdateById failed: `+z.SpaceName())
}

func (z *ZzzMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := z.Adapter.Delete(z.SpaceName(), z.UniqueIndexId(), A.X{z.Id})
	return !L.IsError(err, `Zzz.DoDeletePermanentById failed: `+z.SpaceName())
}

// func (z *ZzzMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := z.Adapter.Upsert(z.SpaceName(), z.ToArray(), A.X{
//		A.X{`=`, 0, z.Id},
//		A.X{`=`, 1, z.CreatedAt},
//	})
//	return !L.IsError(err, `Zzz.DoUpsert failed: `+z.SpaceName())
// }

// insert, error if exists
func (z *ZzzMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := z.Adapter.Insert(z.SpaceName(), z.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			z.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `Zzz.DoInsert failed: `+z.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (z *ZzzMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := z.Adapter.Replace(z.SpaceName(), z.ToArray())
	return !L.IsError(err, `Zzz.DoReplace failed: `+z.SpaceName())
}

func (z *ZzzMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != z.Id {
		z.mutations = append(z.mutations, A.X{`=`, 0, val})
		z.Id = val
		return true
	}
	return false
}

func (z *ZzzMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != z.CreatedAt {
		z.mutations = append(z.mutations, A.X{`=`, 1, val})
		z.CreatedAt = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

