package Tt

import (
	"fmt"

	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool"
)

type Adapter struct {
	*tarantool.Connection
	Reconnect func() *tarantool.Connection
}

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
