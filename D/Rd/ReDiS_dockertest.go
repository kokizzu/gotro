package Rd

import (
	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
	"github.com/rueian/rueidis"
)

type RdDockerTest struct {
	// username and password must use conf https://stackoverflow.com/a/70808430/1620210
	Database int
	Image    string
	pool     *D.DockerTest
	Port     string
}

// https://hub.docker.com/_/redis
func (in *RdDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `redis`,
		Name:       `dockertest-redis-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env:        []string{},
	}
}

func (in *RdDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *RdDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *RdDockerTest) ConnectCheck(res *dockertest.Resource) (rueidis.Client, error) {
	in.Port = res.GetPort("6379/tcp")
	hostPort := in.pool.HostPort(in.Port)
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{hostPort},
		SelectDB:    in.Database,
	})
	return conn, err
}
