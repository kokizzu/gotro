package Rp

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type RpDockerTest struct {
	Image string
	pool  *D.DockerTest
	Port  string
}

// https://hub.docker.com/_/mysql
/*
node=`docker ps | grep redpanda | cut -f 1 -d ' '`
docker exec -it redpanda1 rpk
*/
func (in *RpDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `vectorized/redpanda`,
		Name:       `dockertest-redpanda-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env:        []string{},
		Cmd: []string{
			`redpanda`, `start`,
			`--overprovisioned`,
			`--smp`, `1`,
			`--memory`, `1G`,
			`--reserve-memory`, `0M`,
			`--node-id`, `0`,
			`--check=false`,
		},
	}
}

func (in *RpDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *RpDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *RpDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("9092/tcp")
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
