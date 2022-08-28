package Rd

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/rueian/rueidis"
	"github.com/stretchr/testify/assert"
)

var globalPool *dockertest.Pool

func prepareDb(onReady func(rd rueidis.Client) int) {
	const dockerRepo = `redis`
	const dockerVer = `7.0.4`
	const rdPort = `6379/tcp`
	const connStr = `127.0.0.1:%s`
	var err error
	if globalPool == nil {
		globalPool, err = dockertest.NewPool(``)
		if err != nil {
			log.Printf("Could not connect to docker: %s\n", err)
			return
		}
	}

	resource, err := globalPool.RunWithOptions(&dockertest.RunOptions{
		Repository: dockerRepo,
		Tag:        dockerVer,
		// Cmd:        []string{`-p 6379:6379`},
	})
	if err != nil {
		log.Printf("Could not start resource: %s\n", err)
		return
	}
	var rd rueidis.Client

	retries := 0
	if err := globalPool.Retry(func() (err error) {
		// m, _ := resource.Container.NetworkSettings.Ports[docker.Port(rdPort+`/tcp`)]
		// log.Println(m)
		// return errors.New(`nothing`)

		s := fmt.Sprintf(connStr, resource.GetPort(rdPort))

		reconnect = func() (*RedisSession, error) {
			retries++
			return TryRedisSession(s, `kl234j23095125125125`, 0, `prefix1`)
		}
		var rc *RedisSession
		rc, err = reconnect()
		if err != nil {
			log.Printf("attempt %d %s\n", retries, err)
			return err
		}
		rd = rc.Pool
		return err
	}); err != nil {
		log.Printf("Could not connect to docker: %s\n", err)
		return
	}

	code := onReady(rd)
	if err := globalPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

var reconnect func() (*RedisSession, error)
var redisConn rueidis.Client

func TestMain(m *testing.M) {
	prepareDb(func(db rueidis.Client) int {
		redisConn = db
		if db != nil {
			return m.Run()
		}
		log.Println(`failed when init db`)
		return 0
	})
}

func TestBasic_Operation(t *testing.T) {
	rds := RedisSession{Pool: redisConn, Prefix: `duar`}
	must := assert.New(t)
	const key1 = `esteh`
	const val1 = `panas`
	t.Run(`setStr`, func(t *testing.T) {
		rds.SetStr(key1, val1)
		must.Equal(val1, rds.GetStr(key1))
	})
	t.Run(`setInt`, func(t *testing.T) {
		const val2 int64 = 123
		rds.SetInt(key1, val2)
		must.Equal(val2, rds.GetInt(key1))
	})
}
