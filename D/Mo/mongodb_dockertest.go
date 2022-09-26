package Mo

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type MoDockerTest struct {
	User     string
	Password string
	Image    string
	pool     *D.DockerTest
	Port     string
}

// https://hub.docker.com/_/mongo
/*
docker-compose exec --user root mongo1 /bin/bash
docker-compose exec --user root mongo1 mongo
*/
func (in *MoDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `mongo`,
		Name:       `dockertest-mongo-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`MONGO_INITDB_ROOT_USERNAME=` + in.User,
			`MONGO_INITDB_ROOT_PASSWORD=` + in.Password,
		},
	}
}

func (in *MoDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *MoDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *MoDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("27017/tcp")
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
