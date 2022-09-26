package Pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kokizzu/gotro/L"
)

type Adapter struct {
	*pgxpool.Pool
	Reconnect func() *pgxpool.Pool
}

func Connect1(user, pass, host, dbName string, port, maxConn int) *pgxpool.Pool {
	const connTpl = `postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d`
	connStr := fmt.Sprintf(connTpl,
		user,
		pass,
		host,
		port,
		dbName,
		maxConn,
	)

	db, err := pgxpool.New(context.Background(), connStr)
	L.PanicIf(err, `pgxpool.Connect `+connStr)
	return db
}
