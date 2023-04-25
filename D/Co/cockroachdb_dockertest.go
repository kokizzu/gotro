package Ca

import (
	"net"
	"time"

	"github.com/ory/dockertest/v3"

	"github.com/kokizzu/gotro/D"
)

type CoDockerTest struct {
	Image          string
	Port           string
	ManagementPort string
	pool           *D.DockerTest
}

// https://hub.docker.com/r/cockroachdb/cockroach
/*
node=`docker ps | grep cockroach | cut -f 1 -d ' '`
docker exec -it $node cockroach sql --insecure
*/
func (in *CoDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
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

func (in *CoDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.ManagementPort = res.GetPort(`8080/tcp`)
	in.Port = res.GetPort("5432/tcp")
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
