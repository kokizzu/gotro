package Tt

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alitto/pond"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

var globalPool *dockertest.Pool

func prepareDb(onReady func(db *tarantool.Connection) int) {
	const dockerRepo = `tarantool/tarantool`
	const imageVer = `2.11.1`
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
	resource, err := globalPool.Run(dockerRepo, imageVer, []string{
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
			db, err = tarantool.Connect(connStr, tarantool.Opts{
				User: dbUser,
				Pass: dbPass,
			})
			L.IsError(err, `tarantool.Connect: `+connStr)
			return db
		}
		reconnect()
		if err != nil {
			return err
		}
		_, err = db.Ping()
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
	a := Adapter{dbConn, reconnect}
	const tableName = `test1`
	t.Run(`create test table`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, &TableProp{
			Fields: []Field{
				{`id`, Unsigned},
				{`name`, String},
			},
			Unique1: `id`,
			Engine:  Vinyl,
		})
		assert.True(t, ok)
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
	})
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
	})
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
	})
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
					row, err := a.Insert(tableName, []any{nil, `name`, age, 1, 2, `c`, true, 1.2, []any{}})
					if !L.IsError(err, `insert failed`) {
						if age%printEvery == 0 {
							tup := row.Tuples()
							if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
								lastId := X.ToU(tup[0][0])
								fmt.Println(`lastId`, lastId)
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
