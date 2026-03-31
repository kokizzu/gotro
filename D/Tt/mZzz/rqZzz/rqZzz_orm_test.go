package rqZzz

import (
	"testing"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/stretchr/testify/assert"
)

func TestZzzSchemaHelpers(t *testing.T) {
	z := NewZzz(nil)
	assert.Equal(t, `zzz`, z.SpaceName())
	assert.Equal(t, `"zzz"`, z.SqlTableName())
	assert.Equal(t, `id`, z.UniqueIndexId())
	assert.Equal(t, `coords`, z.SpatialIndexCoords())
	assert.Equal(t, `name`, z.UniqueIndexName())
	assert.Equal(t, 0, z.IdxId())
	assert.Equal(t, `"id"`, z.SqlId())
	assert.Equal(t, 1, z.IdxCreatedAt())
	assert.Equal(t, `"createdAt"`, z.SqlCreatedAt())
	assert.Equal(t, 2, z.IdxCoords())
	assert.Equal(t, `"coords"`, z.SqlCoords())
	assert.Equal(t, 3, z.IdxName())
	assert.Equal(t, `"name"`, z.SqlName())
	assert.Equal(t, 4, z.IdxHeightMeter())
	assert.Equal(t, `"heightMeter"`, z.SqlHeightMeter())
	assert.NotNil(t, z.ToUpdateArray())
	assert.Contains(t, z.SqlSelectAllFields(), `"heightMeter"`)
	assert.Contains(t, z.SqlSelectAllUncensoredFields(), `"coords"`)
}

func TestZzzToAndFromArray(t *testing.T) {
	z := &Zzz{
		Id:          5,
		CreatedAt:   123,
		Coords:      []any{1.2, 3.4},
		Name:        `mountain`,
		HeightMeter: 3456.7,
	}
	arr := z.ToArray()
	assert.Equal(t, uint64(5), arr[0])
	assert.Equal(t, int64(123), arr[1])
	assert.Equal(t, []any{1.2, 3.4}, arr[2])
	assert.Equal(t, `mountain`, arr[3])
	assert.Equal(t, 3456.7, arr[4])

	z.Id = 0
	arr = z.ToArray()
	assert.Nil(t, arr[0])

	decoded := (&Zzz{}).FromArray([]any{uint64(9), int64(8), []any{7, 6}, `peak`, 4.2})
	assert.Equal(t, uint64(9), decoded.Id)
	assert.Equal(t, int64(8), decoded.CreatedAt)
	assert.Equal(t, []any{7, 6}, decoded.Coords)
	assert.Equal(t, `peak`, decoded.Name)
	assert.Equal(t, 4.2, decoded.HeightMeter)

	decoded2 := (&Zzz{}).FromUncensoredArray([]any{uint64(2), int64(3), []any{1, 2}, `hill`, 9.1})
	assert.Equal(t, uint64(2), decoded2.Id)
	assert.Equal(t, int64(3), decoded2.CreatedAt)
	assert.Equal(t, []any{1, 2}, decoded2.Coords)
	assert.Equal(t, `hill`, decoded2.Name)
	assert.Equal(t, 9.1, decoded2.HeightMeter)
}

func TestZzzFieldTypeMap(t *testing.T) {
	assert.Equal(t, Tt.Unsigned, ZzzFieldTypeMap[`id`])
	assert.Equal(t, Tt.Integer, ZzzFieldTypeMap[`createdAt`])
	assert.Equal(t, Tt.Array, ZzzFieldTypeMap[`coords`])
	assert.Equal(t, Tt.String, ZzzFieldTypeMap[`name`])
	assert.Equal(t, Tt.Double, ZzzFieldTypeMap[`heightMeter`])
}

func TestDbMethodsPanicWithUninitializedAdapter(t *testing.T) {
	z := NewZzz(&Tt.Adapter{})
	z.Id = 1
	z.Name = `abc`

	assert.Panics(t, func() { _ = z.FindById() })
	assert.Panics(t, func() { _ = z.FindByName() })
	assert.Panics(t, func() { _ = z.FindOffsetLimit(0, 10, z.UniqueIndexId()) })
	assert.Panics(t, func() { _, _ = z.FindArrOffsetLimit(0, 10, z.UniqueIndexId()) })
	assert.Panics(t, func() { _ = z.Total() })
}
