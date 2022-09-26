package My

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type MyDockerTest struct {
	Password string
	Image    string
	pool     *D.DockerTest
	Port     string
}

// https://hub.docker.com/_/mysql
/*
docker-compose exec --user root mysql1 /bin/bash
docker-compose exec --user root mysql1 mysql -u root -p
*/
func (in *MyDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `mysql`,
		Name:       `dockertest-mariadb-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`MYSQL_ROOT_PASSWORD=` + in.Password,
		},
	}
}

func (in *MyDockerTest) Image57(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `5.7`)
}

func (in *MyDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *MyDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *MyDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
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
