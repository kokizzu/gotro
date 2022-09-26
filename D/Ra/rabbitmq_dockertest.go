package Ra

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type RaDockerTest struct {
	User           string
	Password       string
	Image          string
	Port           string
	ManagementPort string
	DSN            string
	pool           *D.DockerTest
}

// https://hub.docker.com/_/rabbitmq
func (in *RaDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `rabbitmq`,
		Name:       `dockertest-rabbitmq-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`RABBITMQ_DEFAULT_USER=` + in.User, // for management
			`RABBITMQ_DEFAULT_PASS=` + in.Password,
		},
	}
}

func (in *RaDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *RaDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
	if in.User == `` {
		in.User = `rauser`
	}
	if in.Password == `` {
		in.Password = `rapass`
	}
}

func (in *RaDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.ManagementPort = res.GetPort("8080/tcp")
	in.Port = res.GetPort("5672/tcp")
	hostPort := in.pool.HostPort(in.Port)
	userPass := in.User + `:` + in.Password // default: guest:guest
	in.DSN = `amqp://` + userPass + `@` + hostPort
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
