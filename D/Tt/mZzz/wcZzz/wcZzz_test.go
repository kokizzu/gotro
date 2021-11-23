package wcZzz

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/D/Tt/mZzz"
	"github.com/kokizzu/gotro/L"
	"github.com/kpango/fastime"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/tarantool/go-tarantool"
	"log"
	"os"
	"testing"
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

func TestAutoIncrement(t *testing.T) {
	a := &Tt.Adapter{dbConn, reconnect}
	t.Run(`test zzz table`, func(t *testing.T) {
		ok := a.UpsertTable(mZzz.TableZzz, mZzz.TarantoolTables[mZzz.TableZzz])
		assert.True(t, ok)
	})
	t.Run(`test insert auto increment`, func(t *testing.T) {
		zzz := NewZzzMutator(a)
		now := fastime.Now().Unix()
		zzz.CreatedAt = now
		ok := zzz.DoInsert()
		assert.True(t, ok)
		assert.Greater(t, zzz.Id, uint64(0))
	})
}
