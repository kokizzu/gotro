package wcZzz

import (
	"fmt"
	"testing"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool/v2"
)

func TestSanity(t *testing.T) {
	conn := dbConn

	// migrate
	a := &Tt.Adapter{Connection: conn, Reconnect: reconnect}
	a.MigrateTables(map[Tt.TableName]*Tt.TableProp{
		`bands`: {
			Fields: []Tt.Field{
				{`id`, Tt.Unsigned},
				{`band_name`, Tt.String},
				{`year`, Tt.Unsigned},
			},
			AutoIncrementId: true,
			Unique1:         `band_name`,
			Indexes:         []string{`year`},
		},
	},
	)

	// Insert data
	tuples := [][]interface{}{
		{1, "Roxette", 1986},
		{2, "Scorpions", 1965},
		{3, "Ace of Base", 1987},
		{4, "The Beatles", 1960},
	}
	var futures []*tarantool.Future
	for _, tuple := range tuples {
		request := tarantool.NewInsertRequest("bands").Tuple(tuple)
		futures = append(futures, conn.Do(request))
	}
	fmt.Println("Inserted tuples:")
	for _, future := range futures {
		result, err := future.Get()
		assert.NoError(t, err)
		fmt.Println(result)
	}

	// Select by primary key
	data, err := conn.Do(
		tarantool.NewSelectRequest("bands").
			Limit(10).
			Iterator(tarantool.IterEq).
			Key([]interface{}{uint(1)}),
	).Get()
	assert.NoError(t, err)
	fmt.Println("Tuple selected by the primary key value:", data)

	// Select by secondary key
	data, err = conn.Do(
		tarantool.NewSelectRequest("bands").
			Index("band_name").
			Limit(10).
			Iterator(tarantool.IterEq).
			Key([]interface{}{"The Beatles"}),
	).Get()
	assert.NoError(t, err)
	fmt.Println("Tuple selected by the secondary key value:", data)

	// Update
	data, err = conn.Do(
		tarantool.NewUpdateRequest("bands").
			Key(tarantool.IntKey{2}).
			Operations(tarantool.NewOperations().Assign(1, "Pink Floyd")),
	).Get()
	assert.NoError(t, err)
	fmt.Println("Updated tuple:", data)

	// Upsert
	data, err = conn.Do(
		tarantool.NewUpsertRequest("bands").
			Tuple([]interface{}{uint(5), "The Rolling Stones", 1962}).
			Operations(tarantool.NewOperations().Assign(1, "The Doors")),
	).Get()
	assert.NoError(t, err)
	fmt.Println(data)

	// Replace
	data, err = conn.Do(
		tarantool.NewReplaceRequest("bands").
			Tuple([]interface{}{1, "Queen", 1970}),
	).Get()
	assert.NoError(t, err)
	fmt.Println("Replaced tuple:", data)

	// Delete
	data, err = conn.Do(
		tarantool.NewDeleteRequest("bands").
			Key([]interface{}{uint(5)}),
	).Get()
	assert.NoError(t, err)
	fmt.Println("Deleted tuple:", data)

	// Create storProc
	data, err = conn.Do(
		tarantool.NewEvalRequest(`
box.schema.func.create('get_bands_older_than', {
    body = [[
    function(year)
        return box.space.bands.index.year:select({ year }, { iterator = 'LT', limit = 10 })
    end
    ]]
})`)).Get()
	assert.NoError(t, err)
	fmt.Println(data)

	// Call
	data, err = conn.Do(
		tarantool.NewCallRequest("get_bands_older_than").Args([]interface{}{1966}),
	).Get()
	assert.NoError(t, err)
	fmt.Println("Stored procedure result:", data)

	// Close connection
	//conn.CloseGraceful()
	//fmt.Println("Connection is closed")
}
