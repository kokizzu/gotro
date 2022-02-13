package Rd

import (
	"errors"
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
	const dockerRepo = `redislabs/rejson`
	const dockerVer = `2.0.6`
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

	if err := globalPool.Retry(func() error {
		var err error
		// m, _ := resource.Container.NetworkSettings.Ports[docker.Port(rdPort+`/tcp`)]
		// log.Println(m)
		// return errors.New(`nothing`)

		s := fmt.Sprintf(connStr, resource.GetPort(rdPort))

		reconnect = func() *RedisSession {
			rs := NewRedisSession(s, ``, 0, `nihao`)
			return rs
		}
		rc := reconnect()
		if rc == nil {
			return errors.New(`empty connection`)
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

var reconnect func() *RedisSession
var redisConn rueidis.Client

func TestMain(m *testing.M) {
	// _, err := rueidis.NewClient(rueidis.ClientOption{
	// 	InitAddress: []string{`127.0.0.1:6379`},
	// })
	// log.Println(err)
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
	t.Run(`set string`, func(t *testing.T) {
		key := `esteh`
		value := `panas`
		rds.SetStr(key, value)
		must.Equal(value, rds.GetStr(key))
	})
	t.Run(`must error`, func(t *testing.T) {
		key := `megu`
		val := `kaboom`
		xVal := `kaphew`
		rds.SetStr(key, val)
		must.NotEqual(xVal, rds.GetStr(key))
	})

}
