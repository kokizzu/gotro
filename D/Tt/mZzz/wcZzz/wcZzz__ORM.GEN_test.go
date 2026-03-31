package wcZzz

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`context`
	`fmt`
	`log`
	`os`
	`testing`
	`time`
	`github.com/kokizzu/gotro/D/Tt/mZzz`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/S`
	`github.com/ory/dockertest/v3`
	`github.com/stretchr/testify/assert`
	`github.com/tarantool/go-tarantool/v2`
)

var globalPool *dockertest.Pool
var reconnect func() *tarantool.Connection
var dbConn *tarantool.Connection

func prepareDb(onReady func(db *tarantool.Connection) int) {
	const dockerRepo = "tarantool/tarantool"
	const dockerVer = "3.1"
	const ttPort = "3301/tcp"
	const dbConnStr = "127.0.0.1:%s"
	const dbUser = "guest"
	const dbPass = ""
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
			if err != nil && !S.Contains(err.Error(), "failed to read greeting: EOF") {
				L.IsError(err, "tarantool.Connect")
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
		os.Exit(onReady(nil))
	}
	code := onReady(db)
	if err := globalPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestMain(m *testing.M) {
	prepareDb(func(db *tarantool.Connection) int {
		dbConn = db
		return m.Run()
	})
}

func TestGeneratedSanity(t *testing.T) {
	if dbConn == nil {
		t.Skip("docker unavailable")
	}
	conn := dbConn
	a := &Tt.Adapter{Connection: conn, Reconnect: reconnect}
	a.MigrateTables(map[Tt.TableName]*Tt.TableProp{
		"bands": {
			Fields: []Tt.Field{
				{"id", Tt.Unsigned},
				{"band_name", Tt.String},
				{"year", Tt.Unsigned},
			},
			AutoIncrementId: true,
			Unique1:         "band_name",
			Indexes:         []string{"year"},
		},
	})
	tuples := [][]any{
		{1, "Roxette", 1986},
		{2, "Scorpions", 1965},
		{3, "Ace of Base", 1987},
		{4, "The Beatles", 1960},
	}
	for _, tuple := range tuples {
		_, err := conn.Do(tarantool.NewInsertRequest("bands").Tuple(tuple)).Get()
		assert.NoError(t, err)
	}
	_, err := conn.Do(tarantool.NewSelectRequest("bands").Limit(10).Iterator(tarantool.IterEq).Key([]any{uint(1)})).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewSelectRequest("bands").Index("band_name").Limit(10).Iterator(tarantool.IterEq).Key([]any{"The Beatles"})).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewUpdateRequest("bands").Key(tarantool.IntKey{2}).Operations(tarantool.NewOperations().Assign(1, "Pink Floyd"))).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewUpsertRequest("bands").Tuple([]any{uint(5), "The Rolling Stones", 1962}).Operations(tarantool.NewOperations().Assign(1, "The Doors"))).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewReplaceRequest("bands").Tuple([]any{1, "Queen", 1970})).Get()
	assert.NoError(t, err)
	_, err = conn.Do(tarantool.NewDeleteRequest("bands").Key([]any{uint(5)})).Get()
	assert.NoError(t, err)
}

func TestGeneratedZzzUnit(t *testing.T) {
	m := NewZzzMutator(nil)
	assert.NotNil(t, m)
	assert.False(t, m.HaveMutation())
	assert.Empty(t, m.Logs())
	assert.True(t, m.SetId(uint64(1)))
	assert.False(t, m.SetId(uint64(1)))
	assert.True(t, m.SetCreatedAt(int64(1)))
	assert.False(t, m.SetCreatedAt(int64(1)))
	assert.True(t, m.SetCoords([]any{1.1, 2.2}))
	assert.True(t, m.SetName("sample"))
	assert.False(t, m.SetName("sample"))
	assert.True(t, m.SetHeightMeter(1.5))
	assert.False(t, m.SetHeightMeter(1.5))
	assert.True(t, m.HaveMutation())
	from := m.Zzz
	assert.True(t, m.SetAll(from, nil, nil))
	m2 := NewZzzMutator(nil)
	fromZero := m2.Zzz
	fromZero.Coords = nil
	assert.False(t, m.SetAll(fromZero, nil, nil))
	m.ClearMutations()
	assert.False(t, m.HaveMutation())
	assert.True(t, m.DoUpdateById())
	assert.True(t, m.DoUpdateByName())
}

func TestGeneratedZzzCRUD(t *testing.T) {
	if dbConn == nil {
		t.Skip("docker unavailable")
	}
	a := &Tt.Adapter{Connection: dbConn, Reconnect: reconnect}
	ok := a.UpsertTable(mZzz.TableZzz, mZzz.TarantoolTables[mZzz.TableZzz])
	assert.True(t, ok)
	_ = a.TruncateTable(string(mZzz.TableZzz))
	seed := func() *ZzzMutator {
		x := NewZzzMutator(a)
		x.CreatedAt = int64(1)
		x.Coords = []any{1.1, 2.2}
		x.Name = "sample"
		x.HeightMeter = 1.5
		assert.True(t, x.DoInsert())
		assert.Greater(t, x.Id, uint64(0))
		return x
	}
	rec0 := seed()
	readInsert0 := NewZzzMutator(a)
	readInsert0.Id = rec0.Id
	assert.True(t, readInsert0.FindById())
	assert.Equal(t, rec0.Id, readInsert0.Id)
	assert.Equal(t, rec0.CreatedAt, readInsert0.CreatedAt)
	assert.Equal(t, rec0.Coords, readInsert0.Coords)
	assert.Equal(t, rec0.Name, readInsert0.Name)
	assert.Equal(t, rec0.HeightMeter, readInsert0.HeightMeter)
	assert.True(t, rec0.DoUpdateById())
	readUpdate0 := NewZzzMutator(a)
	readUpdate0.Id = rec0.Id
	assert.True(t, readUpdate0.FindById())
	assert.Equal(t, rec0.Id, readUpdate0.Id)
	assert.Equal(t, rec0.CreatedAt, readUpdate0.CreatedAt)
	assert.Equal(t, rec0.Coords, readUpdate0.Coords)
	assert.Equal(t, rec0.Name, readUpdate0.Name)
	assert.Equal(t, rec0.HeightMeter, readUpdate0.HeightMeter)
	assert.True(t, rec0.SetCreatedAt(int64(2)))
	assert.True(t, rec0.DoUpdateById())
	readUpdate0 = NewZzzMutator(a)
	readUpdate0.Id = rec0.Id
	assert.True(t, readUpdate0.FindById())
	assert.Equal(t, rec0.Id, readUpdate0.Id)
	assert.Equal(t, rec0.CreatedAt, readUpdate0.CreatedAt)
	assert.Equal(t, rec0.Coords, readUpdate0.Coords)
	assert.Equal(t, rec0.Name, readUpdate0.Name)
	assert.Equal(t, rec0.HeightMeter, readUpdate0.HeightMeter)
	assert.True(t, rec0.DoOverwriteById())
	readOverwrite0 := NewZzzMutator(a)
	readOverwrite0.Id = rec0.Id
	assert.True(t, readOverwrite0.FindById())
	assert.Equal(t, rec0.Id, readOverwrite0.Id)
	assert.Equal(t, rec0.CreatedAt, readOverwrite0.CreatedAt)
	assert.Equal(t, rec0.Coords, readOverwrite0.Coords)
	assert.Equal(t, rec0.Name, readOverwrite0.Name)
	assert.Equal(t, rec0.HeightMeter, readOverwrite0.HeightMeter)
	rows0 := rec0.FindOffsetLimit(0, 10, rec0.UniqueIndexId())
	assert.NotNil(t, rows0)
	arrRows0, _ := rec0.FindArrOffsetLimit(0, 10, rec0.UniqueIndexId())
	assert.NotNil(t, arrRows0)
	assert.GreaterOrEqual(t, rec0.Total(), int64(0))
	assert.True(t, rec0.DoDeletePermanentById())
	assert.False(t, rec0.FindById())
	rec1 := seed()
	readInsert1 := NewZzzMutator(a)
	readInsert1.Name = rec1.Name
	assert.True(t, readInsert1.FindByName())
	assert.Equal(t, rec1.Id, readInsert1.Id)
	assert.Equal(t, rec1.CreatedAt, readInsert1.CreatedAt)
	assert.Equal(t, rec1.Coords, readInsert1.Coords)
	assert.Equal(t, rec1.Name, readInsert1.Name)
	assert.Equal(t, rec1.HeightMeter, readInsert1.HeightMeter)
	assert.True(t, rec1.DoUpdateByName())
	readUpdate1 := NewZzzMutator(a)
	readUpdate1.Name = rec1.Name
	assert.True(t, readUpdate1.FindByName())
	assert.Equal(t, rec1.Id, readUpdate1.Id)
	assert.Equal(t, rec1.CreatedAt, readUpdate1.CreatedAt)
	assert.Equal(t, rec1.Coords, readUpdate1.Coords)
	assert.Equal(t, rec1.Name, readUpdate1.Name)
	assert.Equal(t, rec1.HeightMeter, readUpdate1.HeightMeter)
	assert.True(t, rec1.SetCreatedAt(int64(2)))
	assert.True(t, rec1.DoUpdateByName())
	readUpdate1 = NewZzzMutator(a)
	readUpdate1.Name = rec1.Name
	assert.True(t, readUpdate1.FindByName())
	assert.Equal(t, rec1.Id, readUpdate1.Id)
	assert.Equal(t, rec1.CreatedAt, readUpdate1.CreatedAt)
	assert.Equal(t, rec1.Coords, readUpdate1.Coords)
	assert.Equal(t, rec1.Name, readUpdate1.Name)
	assert.Equal(t, rec1.HeightMeter, readUpdate1.HeightMeter)
	assert.True(t, rec1.DoOverwriteByName())
	readOverwrite1 := NewZzzMutator(a)
	readOverwrite1.Name = rec1.Name
	assert.True(t, readOverwrite1.FindByName())
	assert.Equal(t, rec1.Id, readOverwrite1.Id)
	assert.Equal(t, rec1.CreatedAt, readOverwrite1.CreatedAt)
	assert.Equal(t, rec1.Coords, readOverwrite1.Coords)
	assert.Equal(t, rec1.Name, readOverwrite1.Name)
	assert.Equal(t, rec1.HeightMeter, readOverwrite1.HeightMeter)
	rows1 := rec1.FindOffsetLimit(0, 10, rec1.UniqueIndexName())
	assert.NotNil(t, rows1)
	arrRows1, _ := rec1.FindArrOffsetLimit(0, 10, rec1.UniqueIndexName())
	assert.NotNil(t, arrRows1)
	assert.GreaterOrEqual(t, rec1.Total(), int64(0))
	assert.True(t, rec1.DoDeletePermanentByName())
	assert.False(t, rec1.FindByName())
	deletedCheck1 := NewZzzMutator(a)
	deletedCheck1.Id = rec1.Id
	assert.False(t, deletedCheck1.FindById())
	u := NewZzzMutator(a)
	u.CreatedAt = int64(1)
	u.Coords = []any{1.1, 2.2}
	u.Name = "sample"
	u.HeightMeter = 1.5
	assert.True(t, u.DoUpsertById())
	assert.Greater(t, u.Id, uint64(0))
	uRead := NewZzzMutator(a)
	uRead.Id = u.Id
	assert.True(t, uRead.FindById())
	assert.Equal(t, u.Id, uRead.Id)
	assert.Equal(t, u.CreatedAt, uRead.CreatedAt)
	assert.Equal(t, u.Coords, uRead.Coords)
	assert.Equal(t, u.Name, uRead.Name)
	assert.Equal(t, u.HeightMeter, uRead.HeightMeter)
	assert.True(t, u.SetCreatedAt(int64(2)))
	assert.True(t, u.DoUpsertById())
	uRead = NewZzzMutator(a)
	uRead.Id = u.Id
	assert.True(t, uRead.FindById())
	assert.Equal(t, u.Id, uRead.Id)
	assert.Equal(t, u.CreatedAt, uRead.CreatedAt)
	assert.Equal(t, u.Coords, uRead.Coords)
	assert.Equal(t, u.Name, uRead.Name)
	assert.Equal(t, u.HeightMeter, uRead.HeightMeter)
	assert.True(t, u.DoDeletePermanentById())
}

