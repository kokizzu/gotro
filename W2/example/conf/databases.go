package conf

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool"
)

func ConnectTarantool() *tarantool.Connection {
	hostPort := fmt.Sprintf(`%s:%s`,
		TARANTOOL_HOST,
		TARANTOOL_PORT,
	)
	taran, err := tarantool.Connect(hostPort, tarantool.Opts{
		User: TARANTOOL_USER,
		Pass: TARANTOOL_PASS,
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
