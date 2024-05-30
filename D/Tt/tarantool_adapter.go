package Tt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool/v2"
)

type Adapter struct {
	*tarantool.Connection
	Reconnect func() *tarantool.Connection
}

func (a *Adapter) RetryDo(op tarantool.Request, times ...int) ([]any, error) {
	count := 3
	if len(times) == 1 {
		count = times[0]
	}
	for {
		res, err := a.Connection.Do(op).Get()
		if err != nil {
			var e tarantool.ClientError
			if errors.As(err, &e) && e.Code == 16385 { // using closed connection (0x4001)
				a.Connection = a.Reconnect()
				count--
				if count == 0 {
					return res, err
				}
				continue
			}
		}
		return res, err
	}
}

func (a *Adapter) RetryDoResp(op tarantool.Request, times ...int) (tarantool.Response, error) {
	count := 3
	if len(times) == 1 {
		count = times[0]
	}
	for {
		future, err := a.Connection.Do(op).GetResponse()
		if err != nil {
			var e tarantool.ClientError
			if errors.As(err, &e) && e.Code == 16385 { // using closed connection (0x4001)
				a.Connection = a.Reconnect()
				count--
				if count == 0 {
					return future, err
				}
				continue
			}
		}
		return future, err
	}
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
