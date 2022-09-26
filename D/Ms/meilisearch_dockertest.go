package Ms

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type MsDockerTest struct {
	MasterKey string
	Image     string
	pool      *D.DockerTest
	Port      string
}

// https://hub.docker.com/r/getmeili/meilisearch
func (in *MsDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `getmeili/meilisearch`,
		Name:       `dockertest-meili-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env:        []string{},
		Cmd:        []string{`--master-key=` + in.MasterKey},
	}
}

func (in *MsDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *MsDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *MsDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("7700/tcp")
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
