package Ch

import (
	"database/sql"
)

type Adapter struct {
	*sql.DB
	Reconnect func() *sql.DB
}
