package saZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/Ch/clickhouse_orm_generator.go

import (
	`database/sql`
	`fmt`
	`log`
	`os`
	`testing`
	`time`

	_ `github.com/ClickHouse/clickhouse-go/v2`
	`github.com/ory/dockertest/v3`
	`github.com/stretchr/testify/assert`
	`github.com/kokizzu/gotro/D/Ch`
	`github.com/kokizzu/gotro/L`
)

var globalPool *dockertest.Pool
var dbConn *sql.DB
var reconnect func() *sql.DB

func prepareDb(onReady func(db *sql.DB) int) {
	const dockerRepo = "yandex/clickhouse-server"
	const dockerVer = "latest"
	const chPort = "9000/tcp"
	const dbDriver = "clickhouse"
	const dbConnStr = "tcp://127.0.0.1:%s?debug=true"
	var err error
	if globalPool == nil {
		globalPool, err = dockertest.NewPool("")
		if err != nil {
			log.Printf("Could not connect to docker: %s\n", err)
			os.Exit(onReady(nil))
		}
	}
	resource, err := globalPool.Run(dockerRepo, dockerVer, []string{})
	if err != nil {
		log.Printf("Could not start resource: %s\n", err)
		os.Exit(onReady(nil))
	}
	var db *sql.DB
	if err := globalPool.Retry(func() error {
		var err error
		connStr := fmt.Sprintf(dbConnStr, resource.GetPort(chPort))
		reconnect = func() *sql.DB {
			db, err = sql.Open(dbDriver, connStr)
			L.IsError(err, "sql.Open: "+connStr)
			return db
		}
		reconnect()
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Printf("Could not connect to docker: %s\n", err)
		os.Exit(onReady(nil))
	}
	code := onReady(db)
	if err := globalPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestMain(m *testing.M) {
	prepareDb(func(db *sql.DB) int {
		dbConn = db
		return m.Run()
	})
}


func TestGeneratedZzzHelpers(t *testing.T) {
	obj := NewZzz(nil)
	assert.NotNil(t, obj)
	assert.NotEmpty(t, obj.TableName())
	assert.NotEmpty(t, obj.SqlTableName())
	obj.Id = uint64(1)
	obj.CreatedAt = time.Now().UTC().Truncate(time.Second)
	obj.Name = "sample"
	obj.HeightMeter = 1.5
	assert.NotEmpty(t, obj.SqlInsert())
	assert.NotEmpty(t, obj.SqlCount())
	assert.NotEmpty(t, obj.SqlSelectAllFields())
	assert.NotEmpty(t, obj.SqlAllFields())
	arr := obj.ToArray()
	assert.Len(t, arr, 4)
	params := obj.SqlInsertParam()
	assert.Len(t, params, 4)
	_, ok := ZzzFieldTypeMap[`id`]
	assert.True(t, ok)
	assert.Equal(t, 0, obj.IdxId())
	assert.Equal(t, `id`, obj.SqlId())
	assert.Equal(t, 1, obj.IdxCreatedAt())
	assert.Equal(t, `createdAt`, obj.SqlCreatedAt())
	assert.Equal(t, 2, obj.IdxName())
	assert.Equal(t, `name`, obj.SqlName())
	assert.Equal(t, 3, obj.IdxHeightMeter())
	assert.Equal(t, `heightMeter`, obj.SqlHeightMeter())
	prep, ok := Preparators[obj.TableName()]
	assert.True(t, ok)
	assert.NotNil(t, prep)
}

func TestGeneratedZzzCRUD(t *testing.T) {
	if dbConn == nil {
		t.Skip("docker unavailable")
	}
	a := &Ch.Adapter{DB: dbConn, Reconnect: reconnect}
	obj := NewZzz(a)
	ok := a.UpsertTable(obj.TableName(), &Ch.TableProp{
Fields: []Ch.Field{
{`id`, Ch.UInt64},
{`createdAt`, Ch.DateTime},
{`name`, Ch.String},
{`heightMeter`, Ch.Float64},
},
Engine: `ReplacingMergeTree`,
Orders: []string{`id`},
})
	assert.True(t, ok)
	_, _ = a.Exec("TRUNCATE TABLE " + string(obj.TableName()))
	row1 := NewZzz(a)
	row1.Id = uint64(1)
	row1.CreatedAt = time.Now().UTC().Truncate(time.Second)
	row1.Name = "sample"
	row1.HeightMeter = 1.5
	_, err := a.Exec(row1.SqlInsert(), row1.SqlInsertParam()...)
	assert.NoError(t, err)
	rows, err := a.Query("SELECT " + row1.SqlSelectAllFields() + " FROM " + row1.SqlTableName() + " LIMIT 1")
	assert.NoError(t, err)
	assert.True(t, rows.Next())
	got := NewZzz(a)
	assert.NoError(t, got.ScanRowAllCols(rows))
	assert.NoError(t, rows.Close())
	row2 := NewZzz(a)
	row2.Id = uint64(2)
	row2.CreatedAt = time.Now().UTC().Add(time.Second).Truncate(time.Second)
	row2.Name = "sample2"
	row2.HeightMeter = 2.5
	_, err = a.Exec(row2.SqlInsert(), row2.SqlInsertParam()...)
	assert.NoError(t, err)
	rows, err = a.Query("SELECT " + row1.SqlSelectAllFields() + " FROM " + row1.SqlTableName() + " LIMIT 10")
	assert.NoError(t, err)
	parsed, err := row1.ScanRowsAllCols(rows, 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(parsed), 2)
	var cnt int64
	assert.NoError(t, a.QueryRow(row1.SqlCount()).Scan(&cnt))
	assert.GreaterOrEqual(t, cnt, int64(1))
}

