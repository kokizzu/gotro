package Ts

import (
	"net"
	"time"

	"github.com/ory/dockertest/v3"

	"github.com/kokizzu/gotro/D"
)

type TsDockerTest struct {
	ApiKey string
	Image  string
	Port   string
	pool   *D.DockerTest
}

// ImageVersion https://hub.docker.com/r/typesense/typesense
/*
default empty api key = 123
*/
func (in *TsDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `typesense/typesense`,
		Name:       `dockertest-typesense-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env:        []string{},
		Cmd: []string{
			`--data-dir /data `,
			`--api-key=` + in.ApiKey,
		},
	}
}

func (in *TsDockerTest) ImageStable(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `0.23.1`)
}

func (in *TsDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *TsDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("8108/tcp")
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
