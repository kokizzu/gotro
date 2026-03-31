package wcZzz

import (
	"testing"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/D/Tt/mZzz/rqZzz"
	"github.com/kokizzu/gotro/M"
	"github.com/stretchr/testify/assert"
)

func TestNewMutatorAndMutationBookkeeping(t *testing.T) {
	z := NewZzzMutator(nil)
	assert.NotNil(t, z)
	assert.Empty(t, z.Logs())
	assert.False(t, z.HaveMutation())
	assert.NotNil(t, z.Coords)
	assert.Equal(t, 0, len(z.Coords))

	assert.False(t, z.SetName(``))
	assert.True(t, z.SetName(`alpha`))
	assert.True(t, z.SetCreatedAt(123))
	assert.True(t, z.SetCoords([]any{1, 2}))
	assert.True(t, z.SetHeightMeter(2.4))
	assert.True(t, z.SetId(9))
	assert.True(t, z.HaveMutation())
	assert.Len(t, z.Logs(), 5)

	z.ClearMutations()
	assert.False(t, z.HaveMutation())
	assert.Empty(t, z.Logs())
	assert.True(t, z.DoUpdateById())
	assert.True(t, z.DoUpdateByName())
	assert.True(t, z.DoUpsertById())
}

func TestSetAll(t *testing.T) {
	z := NewZzzMutator(nil)

	changed := z.SetAll(rqZzz.Zzz{
		Id:          7,
		CreatedAt:   88,
		Coords:      []any{12.3, 45.6},
		Name:        `  Everest  `,
		HeightMeter: 8848.86,
	}, nil, nil)
	assert.True(t, changed)
	assert.Equal(t, uint64(7), z.Id)
	assert.Equal(t, int64(88), z.CreatedAt)
	assert.Equal(t, []any{12.3, 45.6}, z.Coords)
	assert.Equal(t, `Everest`, z.Name)
	assert.Equal(t, 8848.86, z.HeightMeter)

	changed = z.SetAll(rqZzz.Zzz{
		Id:          1,
		CreatedAt:   2,
		Coords:      []any{1},
		Name:        `K2`,
		HeightMeter: 8611,
	}, M.SB{`id`: true, `name`: true}, nil)
	assert.True(t, changed)
	assert.Equal(t, uint64(7), z.Id)
	assert.Equal(t, `Everest`, z.Name)
	assert.Equal(t, int64(2), z.CreatedAt)
	assert.Equal(t, []any{1}, z.Coords)
	assert.Equal(t, 8611.0, z.HeightMeter)

	changed = z.SetAll(rqZzz.Zzz{}, nil, nil)
	assert.False(t, changed)

	changed = z.SetAll(rqZzz.Zzz{}, nil, M.SB{
		`id`:          true,
		`createdAt`:   true,
		`coords`:      true,
		`name`:        true,
		`heightMeter`: true,
	})
	assert.True(t, changed)
	assert.Equal(t, uint64(0), z.Id)
	assert.Equal(t, int64(0), z.CreatedAt)
	assert.Nil(t, z.Coords)
	assert.Equal(t, ``, z.Name)
	assert.Equal(t, 0.0, z.HeightMeter)
}

func TestDbMethodsPanicWithUninitializedAdapter(t *testing.T) {
	z := NewZzzMutator(&Tt.Adapter{})
	z.Id = 1
	z.Name = `abc`
	z.SetName(`def`)

	assert.Panics(t, func() { _ = z.DoOverwriteById() })
	assert.Panics(t, func() { _ = z.DoDeletePermanentById() })
	assert.Panics(t, func() { _ = z.DoOverwriteByName() })
	assert.Panics(t, func() { _ = z.DoUpdateById() })
	assert.Panics(t, func() { _ = z.DoUpdateByName() })
	assert.Panics(t, func() { _ = z.DoDeletePermanentByName() })
	assert.Panics(t, func() { _ = z.DoInsert() })
	assert.Panics(t, func() { _ = z.DoUpsertById() })
}
