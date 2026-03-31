package Tt

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool/v2"
)

type fakeResponse struct{}

func (f fakeResponse) Header() tarantool.Header          { return tarantool.Header{} }
func (f fakeResponse) Decode() ([]interface{}, error)    { return nil, nil }
func (f fakeResponse) DecodeTyped(res interface{}) error { return nil }

func TestNewAdapterUsesConnectFunc(t *testing.T) {
	calls := 0
	a := NewAdapter(func() *tarantool.Connection {
		calls++
		return nil
	})

	assert.NotNil(t, a)
	assert.Equal(t, 1, calls)
	assert.Nil(t, a.Connection)

	_ = a.Reconnect()
	assert.Equal(t, 2, calls)
}

func TestQueryMetaFromVariants(t *testing.T) {
	t.Run(`nil response nil error`, func(t *testing.T) {
		meta := QueryMetaFrom(nil, nil)
		assert.Equal(t, QueryMeta{}, meta)
	})

	t.Run(`nil response with error`, func(t *testing.T) {
		meta := QueryMetaFrom(nil, errors.New(`boom`))
		assert.Equal(t, `boom`, meta.Err)
	})

	t.Run(`non execute response without error`, func(t *testing.T) {
		meta := QueryMetaFrom(fakeResponse{}, nil)
		assert.Equal(t, `not ExecuteResponse`, meta.Err)
	})

	t.Run(`non execute response with error`, func(t *testing.T) {
		meta := QueryMetaFrom(fakeResponse{}, errors.New(`decode failed`))
		assert.Equal(t, `decode failed`, meta.Err)
	})

	t.Run(`execute response`, func(t *testing.T) {
		meta := QueryMetaFrom(&tarantool.ExecuteResponse{}, nil)
		assert.Empty(t, meta.Err)
		assert.Empty(t, meta.Columns)
		assert.Equal(t, uint64(0), meta.SqlInfo.AffectedCount)
	})

	t.Run(`execute response keeps upstream error`, func(t *testing.T) {
		meta := QueryMetaFrom(&tarantool.ExecuteResponse{}, errors.New(`upstream`))
		assert.Equal(t, `upstream`, meta.Err)
	})
}

func TestFieldKeyRenderer(t *testing.T) {
	unsigned := Field{Type: Unsigned}.KeyRenderer()
	assert.Equal(t, `tarantool.UintKey{I:uint(id)}`, unsigned(`id`))

	integer := Field{Type: Integer}.KeyRenderer()
	assert.Equal(t, `tarantool.IntKey{I:int(age)}`, integer(`age`))

	str := Field{Type: String}.KeyRenderer()
	assert.Equal(t, `tarantool.StringKey{S:name}`, str(`name`))

	fallback := Field{Type: Double}.KeyRenderer()
	assert.Equal(t, `DataTypeNotPossibleToBeAKey:double:height`, fallback(`height`))
}

func TestTypeGraphqlAndQuotes(t *testing.T) {
	assert.Equal(t, `ID`, TypeGraphql(Field{Name: `ownerId`, Type: Unsigned}))
	assert.Equal(t, `ID`, TypeGraphql(Field{Name: `createdBy`, Type: Integer}))
	assert.Equal(t, `Int`, TypeGraphql(Field{Name: `count`, Type: Integer}))
	assert.Equal(t, `Float`, TypeGraphql(Field{Name: `ratio`, Type: Double}))
	assert.Equal(t, `String`, TypeGraphql(Field{Name: `name`, Type: String}))
	assert.Equal(t, `Boolean`, TypeGraphql(Field{Name: `isActive`, Type: Boolean}))

	assert.Equal(t, "\n\t`github.com/example/mod`", qi(`github.com/example/mod`))
	assert.Equal(t, `"users"`, dq(`users`))
}

func TestMigrateTablesMisconfigurationPanics(t *testing.T) {
	a := &Adapter{}
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Contains(t, fmt.Sprint(r), `spatial index is not supported in vinyl engine`)
	}()

	a.MigrateTables(map[TableName]*TableProp{
		`bad_table`: {
			Fields:  []Field{{Name: `id`, Type: Unsigned}},
			Engine:  Vinyl,
			Spatial: `coords`,
		},
	})
}
