package Tt

import (
	"context"
	"fmt"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool/v2"
)

type Adapter struct {
	*tarantool.Connection
	Reconnect func() *tarantool.Connection
}

// NewAdapter create new tarantool adapter
// adapter contains helper methods for schema manipulation and query execution
func NewAdapter(connectFunc func() *tarantool.Connection) *Adapter {
	return &Adapter{
		Reconnect:  connectFunc,
		Connection: connectFunc(),
	}
}

// Connect1 is example of connect function
// to connect on terminal locally, use:
// tarantoolctl connect user:password@localhost:3301
func Connect1(host, port, user, pass string) *tarantool.Connection {
	hostPort := fmt.Sprintf(`%s:%s`,
		host,
		port,
	)
	taran, err := tarantool.Connect(context.Background(), tarantool.NetDialer{
		Address:  hostPort,
		User:     user,
		Password: pass,
	}, tarantool.Opts{
		Timeout: 8 * time.Second,
	})
	L.PanicIf(err, `tarantool.Connect `+hostPort)
	return taran
}
