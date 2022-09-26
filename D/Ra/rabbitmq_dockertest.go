package Ra

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type RaDockerTest struct {
	User     string
	Password string
	Image    string
}

// https://hub.docker.com/_/rabbitmq
func (in *RaDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
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

func (in *RaDockerTest) ConnectCheck(res *dockertest.Resource) (dsn string, managementPort string, err error) {
	managementPort = res.GetPort("8080/tcp")
	port := res.GetPort("5672/tcp")
	hostPort := `127.0.0.1:` + port
	userPass := in.User + `:` + in.Password // default: guest:guest
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
