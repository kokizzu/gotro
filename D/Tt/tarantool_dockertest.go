package Tt

import (
	"github.com/ory/dockertest/v3"
	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/D"
)

type TtDockerTest struct {
	User     string
	Password string
	Image    string
	Port     string
	pool     *D.DockerTest
}

// ImageVersion https://hub.docker.com/r/tarantool/tarantool
/*
tarantoolctl connect 3301
*/
func (in *TtDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `tarantool/tarantool`,
		Name:       `dockertest-tarantool-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`TARANTOOL_USER_NAME=` + in.User,
			`TARANTOOL_USER_PASSWORD=` + in.Password,
		},
	}
}

func (in *TtDockerTest) Image3(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `3`)
}

func (in *TtDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *TtDockerTest) ConnectCheck(res *dockertest.Resource) (taran *tarantool.Connection, err error) {
	in.Port = res.GetPort("3301/tcp")
	hostPort := in.pool.HostPort(in.Port)
	taran, err = tarantool.Connect(hostPort, tarantool.Opts{
		User: in.User,
		Pass: in.Password,
	})
	return
}
