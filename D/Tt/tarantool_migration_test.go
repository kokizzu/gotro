package Tt

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alitto/pond"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool/v2"
)

var globalPool *dockertest.Pool

func prepareDb(onReady func(db *tarantool.Connection) int) {
	const dockerRepo = `tarantool/tarantool`
	const imageVer = `3.1`
	const ttPort = `3301/tcp`
	const dbConnStr = `127.0.0.1:%s`
	const dbUser = `guest`
	const dbPass = ``
	var err error
	if globalPool == nil {
		globalPool, err = dockertest.NewPool("")
		if err != nil {
			log.Printf("Could not connect to docker: %s\n", err)
			return
		}
	}
	resource, err := globalPool.Run(dockerRepo, imageVer, []string{ // TODO: update to 3.1
		`TT_READAHEAD=1632000`,       // 10x
		`TT_VINYL_MEMORY=1717986918`, // 16x
		`TT_VINYL_CACHE=536870912`,   // 2x
		`TT_NET_MSG_MAX=76800`,       // 100x
		`TT_MEMTX_MEMORY=900000000`,  // ~3x
		`TT_VINYL_PAGE_SIZE=8192`,    // 1x
	})
	if err != nil {
		log.Printf("Could not start resource: %s\n", err)
		return
	}
	var db *tarantool.Connection
	if err := globalPool.Retry(func() error {
		var err error
		connStr := fmt.Sprintf(dbConnStr, resource.GetPort(ttPort))
		reconnect = func() *tarantool.Connection {
			db, err = tarantool.Connect(context.Background(), tarantool.NetDialer{
				Address:  connStr,
				User:     dbUser,
				Password: dbPass,
			}, tarantool.Opts{
				Timeout: 8 * time.Second,
			})
			if err != nil && !S.Contains(err.Error(), `failed to read greeting: EOF`) {
				L.IsError(err, `tarantool.Connect`)
			}
			return db
		}
		reconnect()
		if err != nil {
			return err
		}
		_, err = db.Do(tarantool.NewPingRequest()).Get()
		return err
	}); err != nil {
		log.Printf("Could not connect to docker: %s\n", err)
		return
	}
	code := onReady(db)
	if err := globalPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

var reconnect func() *tarantool.Connection
var dbConn *tarantool.Connection

func TestMain(m *testing.M) {
	prepareDb(func(db *tarantool.Connection) int {
		dbConn = db
		if db != nil {
			return m.Run()
		}
		return 0
	})
}

func TestMigration(t *testing.T) {
	a := &Adapter{dbConn, reconnect}
	const tableName = `test1`
	dummyTable := &TableProp{
		Fields: []Field{
			{`id`, Unsigned},
			{`name`, String},
		},
		Unique1: `id`,
		Engine:  Vinyl,
	}
	t.Run(`check tables must panic because doesnt exists yet`, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Contains(t, fmt.Sprint(r), "Procedure 'box.space.test1:format' is not defined")
			}
		}()
		CheckTarantoolTables(a, map[TableName]*TableProp{
			tableName: dummyTable,
		})
	})
	t.Run(`create test table`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, dummyTable)
		assert.True(t, ok)

		t.Run(`check tables must not panic because it matches`, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fail()
				}
			}()
			CheckTarantoolTables(a, map[TableName]*TableProp{
				tableName: dummyTable,
			})
		})

		t.Run(`check tables must panic if columns changed`, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Contains(t, fmt.Sprint(r), `please run FIELD_TYPE name migration on test1 for tarantool: string <> integer`)
				}
			}()
			CheckTarantoolTables(a, map[TableName]*TableProp{
				tableName: {
					Fields: []Field{
						{`id`, Unsigned},
						{`name`, Integer},
					},
					Unique1: `id`,
					Engine:  Vinyl,
				},
			})
		})

		t.Run(`add column test table`, func(t *testing.T) {
			ok := a.UpsertTable(tableName, &TableProp{
				Fields: []Field{
					{`id`, Unsigned},
					{`name`, String},
					{`age`, Integer},
				},
				Unique1: `id`,
				Engine:  Vinyl,
			})
			assert.True(t, ok)

			t.Run(`add 2 columns test table`, func(t *testing.T) {
				ok := a.UpsertTable(tableName, &TableProp{
					Fields: []Field{
						{`id`, Unsigned},
						{`name`, String},
						{`age`, Integer},
						{`a`, Unsigned},
						{`b`, Integer},
					},
					Unique1: `id`,
					Engine:  Vinyl,
				})
				assert.True(t, ok)

				t.Run(`auto increment`, func(t *testing.T) {
					ok := a.UpsertTable(tableName, &TableProp{
						Fields: []Field{
							{`id`, Unsigned},
							{`name`, String},
							{`age`, Integer},
							{`a`, Unsigned},
							{`b`, Integer},
						},
						Unique1:         `id`,
						Engine:          Vinyl,
						AutoIncrementId: true,
					})
					assert.True(t, ok)

					t.Run(`auto increment again`, func(t *testing.T) {
						ok := a.UpsertTable(tableName, &TableProp{
							Fields: []Field{
								{`id`, Unsigned},
								{`name`, String},
								{`age`, Integer},
								{`a`, Unsigned},
								{`b`, Integer},
							},
							Unique1:         `id`,
							Engine:          Vinyl,
							Uniques:         []string{`a`, `b`},
							AutoIncrementId: true,
						})
						assert.True(t, ok)
					})
				})
			})
		})
	})
}

func TestMigrationBig(t *testing.T) {
	a := Adapter{dbConn, reconnect}
	const tableName = `test2`
	const insertCount = 4_000_000 / 10 // increase TT_MEMTX_MEMORY if needed
	const threadCount = 32
	oldTable := &TableProp{
		Fields: []Field{
			{`id`, Unsigned},
			{`name`, String},
			{`age`, Integer},
			{`a`, Unsigned},
			{`b`, Integer},
			{`c`, String},
			{`d`, Boolean},
			{`e`, Double},
			{`f`, Array},
		},
		AutoIncrementId: true,
		Unique1:         `age`, // just for example
		Engine:          Memtx,
	}
	t.Run(`create test table`, func(t *testing.T) {
		a.MigrateTables(map[TableName]*TableProp{
			tableName: oldTable,
		})

		t.Run(`insert rows`, func(t *testing.T) {
			DEBUG = false
			pool := pond.New(threadCount, insertCount)
			inserted := uint64(0)
			const printEvery = 100_000
			for i := 0; i < insertCount; i++ {
				pool.Submit(func() {
					age := atomic.AddUint64(&inserted, 1)
					row, err := a.Connection.Do(tarantool.NewInsertRequest(tableName).Tuple([]any{nil, `name`, age, 1, 2, `c`, true, 1.2, []any{}})).Get()
					if !L.IsError(err, `insert failed`) {
						if age%printEvery == 0 {
							if len(row) > 0 {
								if tup, ok := row[0].([]any); ok {
									lastId := X.ToU(tup[0])
									fmt.Println(`lastId`, lastId)
								}
							}
						}
					}
				})
			}
			pool.StopAndWait()
			DEBUG = true

			t.Run(`alter table add 6 more columns`, func(t *testing.T) {
				defer L.TimeTrack(time.Now(), `alter table completed`)
				newTable := oldTable // mutate oldTable
				newTable.Fields = append(newTable.Fields, []Field{
					{`g`, Unsigned},
					{`h`, Integer},
					{`i`, String},
					{`j`, Boolean},
					{`k`, Double},
					{`l`, Array},
				}...)
				a.MigrateTables(map[TableName]*TableProp{
					tableName: newTable,
				})
			})
		})
	})
}
