package Ch

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kokizzu/gotro/L"
)

type Adapter struct {
	*sql.DB
	Reconnect func() *sql.DB
}

func Connect1(user, pass, host, port, dbName string, debug bool) *sql.DB {
	hostPort := host + `:` + port
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{hostPort},
		Auth: clickhouse.Auth{
			Database: dbName, // `default`
			Username: user,   // `default`
			Password: pass,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		Settings: clickhouse.Settings{
			`max_execution_time`: 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: debug,
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	err := conn.Ping()
	L.PanicIf(err, `clickhouse.OpenDB `+user+`@`+hostPort)
	return conn
}

func ConnectLocal(host, port string) *sql.DB {
	connStr := fmt.Sprintf("tcp://%s:%s?debug=true",
		host,
		port,
	)
	click, err := sql.Open(`clickhouse`, connStr)
	L.PanicIf(err, `sql.Open `+connStr)
	return click
}
