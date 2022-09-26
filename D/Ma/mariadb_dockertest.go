package Ma

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type MaDockerTest struct {
	Password string
	Image    string
	pool     *D.DockerTest
	Port     string
}

// https://hub.docker.com/_/mariadb
func (in *MaDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `mariadb`,
		Name:       `dockertest-mariadb-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`MYSQL_ROOT_PASSWORD=` + in.Password,
		},
	}
}

func (in *MaDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *MaDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *MaDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("3306/tcp")
	hostPort := in.pool.HostPort(in.Port)
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
