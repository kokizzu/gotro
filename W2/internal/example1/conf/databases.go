package conf

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool/v2"
)

func ConnectTarantool() *tarantool.Connection {
	hostPort := fmt.Sprintf(`%s:%s`,
		TARANTOOL_HOST,
		TARANTOOL_PORT,
	)
	taran, err := tarantool.Connect(context.Background(), tarantool.NetDialer{
		Address:  hostPort,
		User:     TARANTOOL_USER,
		Password: TARANTOOL_PASS,
	}, tarantool.Opts{
		Timeout: 8 * time.Second,
	})
	L.PanicIf(err, `tarantool.Connect `+hostPort)
	return taran
}

func ConnectClickhouse() *sql.DB {
	connStr := fmt.Sprintf("tcp://%s:%s", // ?debug=true",
		CLICKHOUSE_HOST,
		CLICKHOUSE_PORT,
	)
	click, err := sql.Open(`clickhouse`, connStr)
	L.PanicIf(err, `sql.Open `+connStr)
	return click
}
