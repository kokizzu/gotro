package conf

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/S"
	"github.com/tarantool/go-tarantool/v2"
)

type TarantoolConf struct {
	User string
	Pass string
	Host string
	Port int
}

func EnvTarantool() TarantoolConf {
	return TarantoolConf{
		User: os.Getenv("TARANTOOL_USER"),
		Pass: os.Getenv("TARANTOOL_PASS"),
		Host: os.Getenv("TARANTOOL_HOST"),
		Port: S.ToInt(os.Getenv("TARANTOOL_PORT")),
	}
}

var ErrConnectTarantool = errors.New(`TarantoolConf) Connect`)

func (c TarantoolConf) Connect() (a *Tt.Adapter, err error) {
	hostPort := fmt.Sprintf("%s:%d", c.Host, c.Port)
	var taran *tarantool.Connection
	connectFunc := func() *tarantool.Connection {
		taran, err = tarantool.Connect(context.Background(), tarantool.NetDialer{
			Address:  hostPort,
			User:     c.User,
			Password: c.Pass,
		}, tarantool.Opts{
			Timeout: 8 * time.Second,
		})
		return taran
	}
	taran = connectFunc()
	if taran == nil {
		return nil, WrapError(ErrConnectTarantool, err)
	}
	_, err = taran.Do(tarantool.NewPingRequest()).Get()
	if err != nil {
		return nil, WrapError(ErrConnectTarantool, err)
	}
	a = &Tt.Adapter{
		Connection: taran,
		Reconnect:  connectFunc,
	}
	return
}

func (c TarantoolConf) DebugStr() string {
	return fmt.Sprintf(`%s@%s:%d`, c.User, c.Host, c.Port)
}
