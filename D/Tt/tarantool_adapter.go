package Tt

import (
	"fmt"

	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/L"
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
	taran, err := tarantool.Connect(hostPort, tarantool.Opts{
		User: user,
		Pass: pass,
	})
	L.PanicIf(err, `tarantool.Connect `+hostPort)
	return taran
}
