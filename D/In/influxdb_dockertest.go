package In

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type InDockerTest struct {
	User     string
	Password string
	Image    string
}

// https://hub.docker.com/_/influxdb
func (in *InDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `influxdb`,
		Name:       `dockertest-influxdb-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`INFLUX_USERNAME=` + in.User, // for management
			`INFLUX_PASSWORD=` + in.Password,
		},
	}
}

func (in *InDockerTest) Image18(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `1.8.10`)
}

func (in *InDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
	if in.User == `` {
		in.User = `inuser`
	}
	if in.Password == `` {
		in.Password = `inpass`
	}
}

func (in *InDockerTest) ConnectCheck(res *dockertest.Resource) (dsn string, err error) {
	port := res.GetPort("8086/tcp")
	hostPort := `127.0.0.1:` + port
	userPass := in.User + `:` + in.Password
	dsn = `amqp://` + userPass + `@` + hostPort
	// using net Dial instead of proper driver
	var conn net.Conn
	conn, err = net.DialTimeout("tcp", hostPort, 1*time.Second)
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	}
	return
}
