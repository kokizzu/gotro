package conf

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/S"
)

type ClickhouseConf struct {
	User   string
	Pass   string
	DB     string
	Host   string
	Port   int
	UseSsl bool
}

func EnvClickhouse() ClickhouseConf {
	return ClickhouseConf{
		User:   os.Getenv("CLICKHOUSE_USER"),
		Pass:   os.Getenv("CLICKHOUSE_PASS"),
		DB:     os.Getenv("CLICKHOUSE_DB"),
		Host:   os.Getenv("CLICKHOUSE_HOST"),
		Port:   S.ToInt(os.Getenv("CLICKHOUSE_PORT")),
		UseSsl: os.Getenv("CLICKHOUSE_USE_SSL") == "true",
	}
}

var ErrConnectClickhouse = errors.New(`ClickhouseConf) Connect`)

func (c ClickhouseConf) Connect() (a *Ch.Adapter, err error) {
	hostPort := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conf := &clickhouse.Options{
		Addr: []string{hostPort},
		Auth: clickhouse.Auth{
			Database: c.DB,
			Username: c.User,
			Password: c.Pass,
		},
		Settings: clickhouse.Settings{
			`max_execution_time`: 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		//Debug: IsDebug(),
	}
	if c.UseSsl {
		conf.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	connectFunc := func() *sql.DB {
		conn := clickhouse.OpenDB(conf)
		conn.SetMaxIdleConns(5)
		conn.SetMaxOpenConns(10)
		conn.SetConnMaxLifetime(time.Hour)
		return conn
	}
	conn := connectFunc()
	err = conn.Ping()
	if err != nil {
		return nil, WrapError(ErrConnectClickhouse, err)
	}
	a = &Ch.Adapter{
		DB:        conn,
		Reconnect: connectFunc,
	}
	return a, nil

}

func (c ClickhouseConf) DebugStr() string {
	return fmt.Sprintf("%s@%s:%d", c.User, c.Host, c.Port)
}
