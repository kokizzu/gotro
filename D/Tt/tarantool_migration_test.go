package Tt

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/L"
)

var globalPool *dockertest.Pool

func prepareDb(onReady func(db *tarantool.Connection) int) {
	const dockerRepo = `tarantool/tarantool`
	const dockerVer = `2.7.2`
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
	resource, err := globalPool.Run(dockerRepo, dockerVer, []string{})
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
