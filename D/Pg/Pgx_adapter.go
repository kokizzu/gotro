package Pg

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Adapter struct {
	*pgxpool.Pool
	Reconnect func() *pgxpool.Pool
}
