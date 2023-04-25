package Ss

import (
	"net"
	"time"

	"github.com/ory/dockertest/v3"

	"github.com/kokizzu/gotro/D"
)

type SsDockerTest struct {
	StrongPassword string
	Image          string
	pool           *D.DockerTest
	Port           string
}

// https://hub.docker.com/_/microsoft-mssql-server
/*
docker-compose exec --user root sqlserver1 /bin/bash
docker-compose exec --user root sqlserver1 mongo
*/
func (in *SsDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `mcr.microsoft.com/mssql/server`,
		Name:       `dockertest-mssql-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`ACCEPT_EULA=Y`,
			`SA_PASSWORD=` + in.StrongPassword, // username: sa
		},
	}
}

func (in *SsDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *SsDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *SsDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("1433/tcp")
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
