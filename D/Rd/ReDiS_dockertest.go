package Rd

import (
	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
	"github.com/rueian/rueidis"
)

type RdDockerTest struct {
	// must use conf https://stackoverflow.com/a/70808430/1620210
	//User     string
	//Password string
	Database int
	Image    string
}

// https://hub.docker.com/_/redis
func (in *RdDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
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
	port := res.GetPort("6379/tcp")
	hostPort := `127.0.0.1:` + port
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{hostPort},
		SelectDB:    in.Database,
	})
	return conn, err
}
