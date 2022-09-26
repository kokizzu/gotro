package Pg

import (
	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type PgDockerTest struct {
	User     string
	Password string
	Database string
	Image    string
}

func (in *PgDockerTest) Image96(pool *D.DockerTest) *dockertest.RunOptions {
	in.SetDefaults()
	return &dockertest.RunOptions{
		Repository: `postgres`,
		Name:       `dockertest-postgres-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`POSTGRES_USER=` + in.User,
			`POSTGRES_PASSWORD=` + in.Password,
			`POSTGRES_DB=` + in.Database,
		},
	}
}

func (in *PgDockerTest) SetDefaults() {
	if in.Image == `` {
		in.Image = `9.6.14-alpine`
	}
	if in.User == `` {
		in.User = `pguser`
	}
	if in.Password == `` {
		in.Password = `pgpass`
	}
	if in.Database == `` {
		in.Database = `pgdb`
	}
}
