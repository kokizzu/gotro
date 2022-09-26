package Ms

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type MiDockerTest struct {
	Image       string
	Username    string
	Password    string
	pool        *D.DockerTest
	Port        string
	ConsolePort string
}

// https://hub.docker.com/r/minio/minio
func (in *MiDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `quay.io/minio/minio`,
		Name:       `dockertest-minio-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`MINIO_ROOT_USER=` + in.Username,
			`MINIO_ROOT_PASSWORD=` + in.Password,
		},
		Cmd: []string{
			`server`,
			`/data`,
			`--console-address`, `:9001`,
		},
	}
}

func (in *MiDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *MiDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
	if in.Username == `` {
		in.Username = `miuser`
	}
	if in.Password == `` {
		in.Password = `mipass`
	}
}

func (in *MiDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.ConsolePort = res.GetPort("9001/tcp")
	in.Port = res.GetPort("9000/tcp")
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
