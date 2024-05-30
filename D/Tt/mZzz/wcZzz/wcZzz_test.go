package wcZzz

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/D/Tt/mZzz"
	"github.com/kokizzu/gotro/D/Tt/mZzz/rqZzz"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kpango/fastime"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool/v2"
)

var globalPool *dockertest.Pool

func prepareDb(onReady func(db *tarantool.Connection) int) {
	const dockerRepo = `tarantool/tarantool`
	const dockerVer = `3.1`
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

func TestAutoIncrementAndCRUD(t *testing.T) {
	a := &Tt.Adapter{Connection: dbConn, Reconnect: reconnect}
	t.Run(`test zzz table`, func(t *testing.T) {
		ok := a.UpsertTable(mZzz.TableZzz, mZzz.TarantoolTables[mZzz.TableZzz])
		assert.True(t, ok)
	})
	t.Run(`test insert auto increment`, func(t *testing.T) {
		zzz := NewZzzMutator(a)
		now := fastime.Now().Unix()
		zzz.CreatedAt = now
		lastCoord := []any{12.34, 56.78}
		zzz.Coords = lastCoord
		ok := zzz.DoInsert()
		assert.True(t, ok)
		assert.Greater(t, zzz.Id, uint64(0))
		t.Run(`upsert`, func(t *testing.T) {
			z4 := zzz
			lastCoord = []any{23.45, 6.7}
			z4.SetCoords(lastCoord)
			ok := z4.DoUpsertById()
			assert.True(t, ok)
			t.Run(`select`, func(t *testing.T) {
				z5 := rqZzz.NewZzz(a)
				z5.Id = zzz.Id
				ok := z5.FindById()
				assert.True(t, ok)
				assert.Equal(t, lastCoord, z5.Coords)
			})
		})
		t.Run(`update`, func(t *testing.T) {
			zzz.SetName(`foo`)
			//ok := zzz.DoUpdateById()
			res, err := dbConn.Do(tarantool.NewUpdateRequest(zzz.SpaceName()).
				Key(tarantool.UintKey{uint(zzz.Id)}).
				//Key(A.X{zzz.Id}).
				//Operations(tarantool.NewOperations().Assign(3, zzz.Name))).Get()
				Operations(zzz.mutations)).Get()
			assert.NoError(t, err)
			fmt.Println(res)
			assert.True(t, ok)
			t.Run(`select`, func(t *testing.T) {
				z2 := NewZzzMutator(a)
				z2.Id = zzz.Id
				ok := z2.FindById()
				assert.True(t, ok)
				assert.Equal(t, `foo`, z2.Name)
				assert.Equal(t, lastCoord, z2.Coords)
				t.Run(`delete`, func(t *testing.T) {
					ok := z2.DoDeletePermanentById()
					assert.True(t, ok)
					t.Run(`selectMissing`, func(t *testing.T) {
						z3 := NewZzzMutator(a)
						z3.Id = zzz.Id
						ok := z3.FindById()
						assert.False(t, ok)
					})
				})
			})
		})
	})
}
