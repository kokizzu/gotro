package Ca

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type CoDockerTest struct {
	Image string
}

// https://hub.docker.com/r/cockroachdb/cockroach
func (in *CoDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `cockroachdb/cockroach`,
		Name:       `dockertest-cockroach-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env:        []string{},
		Cmd: []string{
			`start-single-node`,
			`--insecure`,
			`--accept-sql-without-tls`,
		},
	}
}

func (in *CoDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *CoDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *CoDockerTest) ConnectCheck(res *dockertest.Resource) (hostPort string, managementPort string, err error) {
	managementPort = res.GetPort(`8080/tcp`)
	port := res.GetPort("5432/tcp")
	hostPort = `127.0.0.1:` + port
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
