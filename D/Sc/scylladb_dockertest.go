package Sc

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type ScDockerTest struct {
	Image string
	Port  string
	pool  *D.DockerTest
}

// https://hub.docker.com/r/scylladb/scylla
/*
node=`docker ps | grep /scylla: | cut -f 1 -d ' '`
docker exec -it $node cqlsh
*/
func (in *ScDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `scylladb/scylla`,
		Name:       `dockertest-scylla-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env:        []string{},
		Cmd: []string{
			`--smp`, `1`,
			`--memory`, `750M`,
			`--overprovisioned`, `1`,
			`--api-address`, `0.0.0.0`,
			`--developer-mode`, `1`,
		},
	}
}

func (in *ScDockerTest) Image3(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `3`)
}

func (in *ScDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *ScDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("9042/tcp")
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
