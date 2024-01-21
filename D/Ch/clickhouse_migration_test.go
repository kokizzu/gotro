package Ch

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/L"
)

var globalPool *dockertest.Pool

func prepareDb(onReady func(db *sql.DB) int) {
	const dockerRepo = `yandex/clickhouse-server`
	const dockerVer = `latest`
	const chPort = `9000/tcp`
	const dbDriver = "clickhouse"
	const dbConnStr = "tcp://127.0.0.1:%s?debug=true"
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
	var db *sql.DB
	if err := globalPool.Retry(func() error {
		var err error
		connStr := fmt.Sprintf(dbConnStr, resource.GetPort(chPort))
		reconnect = func() *sql.DB {
			db, err = sql.Open(dbDriver, connStr)
			L.IsError(err, `sql.Open: `+connStr)
			return db
		}
		reconnect()
		if err != nil {
			return err
		}
		return db.Ping()
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

var dbConn *sql.DB
var reconnect func() *sql.DB

func TestMain(m *testing.M) {
	prepareDb(func(db *sql.DB) int {
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
	t.Run(`create test must ok`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, &TableProp{
			Fields: []Field{
				{`id`, UInt64},
				{`nom`, String},
			},
			Orders: []string{`id`},
			Engine: ReplacingMergeTree,
		})
		assert.True(t, ok)
	})
	t.Run(`rename column must ok`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, &TableProp{
			Fields: []Field{
				{`id`, UInt64},
				{`name`, String}, // from nom
			},
			Orders: []string{`id`},
			Engine: ReplacingMergeTree,
		})
		assert.True(t, ok)
	})
	t.Run(`add column must ok`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, &TableProp{
			Fields: []Field{
				{`id`, UInt64},
				{`name`, String},
				{`age`, Int32}, // 1 more col
			},
			Orders: []string{`id`},
			Engine: ReplacingMergeTree,
		})
		assert.True(t, ok)
	})
	t.Run(`change non PK column data type must ok`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, &TableProp{
			Fields: []Field{
				{`id`, UInt64},
				{`name`, String},
				{`age`, Float32}, // to float32
			},
			Orders: []string{`id`},
			Engine: ReplacingMergeTree,
		})
		assert.True(t, ok)
	})
	t.Run(`add 2 columns must ok`, func(t *testing.T) {
		ok := a.UpsertTable(tableName, &TableProp{
			Fields: []Field{
				{`id`, UInt64},
				{`name`, String},
				{`age`, Float32},
				{`a`, Int32}, // 2 more cols
				{`b`, Int64},
			},
			Orders: []string{`id`},
			Engine: ReplacingMergeTree,
		})
		assert.True(t, ok)
	})
	t.Run(`create materialized view must ok`, func(t *testing.T) {
		const mvName = `test_mv`
		const mvSourceTable = `mv_source_table`
		ok := a.UpsertTable(mvSourceTable, &TableProp{
			Fields: []Field{
				{`id`, UInt32},
				{`created_at`, DateTime},
				{`name`, String},
			},
			Engine:     ReplacingMergeTree,
			Partitions: []string{`id`},
			Orders:     []string{`id`, `created_at`},
		})
		assert.True(t, ok)
		ok = a.CreateMaterializedViews(mvName, &MVProp{
			SourceTable: mvSourceTable,
			SourceColumns: []string{
				`toDate(toHour(created_at)) AS created_at_hour`,
				`*`,
			},
			Engine:     ReplacingMergeTree,
			Partitions: []string{`id`},
			Orders:     []string{`id`, `created_at_hour`},
		})
		assert.True(t, ok)
	})
}
