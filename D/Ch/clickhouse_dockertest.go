package Ch

import (
	"crypto/tls"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type ChDockerTest struct {
	User     string
	Password string
	Database string
	Image    string
	Port     string
	pool     *D.DockerTest
}

// https://hub.docker.com/r/clickhouse/clickhouse-server
func (in *ChDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `clickhouse/clickhouse-server`,
		Name:       `dockertest-clickhouse-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`CLICKHOUSE_USER=` + in.User,
			`CLICKHOUSE_PASSWORD=` + in.Password,
			`CLICKHOUSE_DB=` + in.Database,
			`CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1`,
		},
	}
}

func (in *ChDockerTest) ImageLatest(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `latest`)
}

func (in *ChDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
	if in.User == `` {
		in.User = `chuser`
	}
	if in.Password == `` {
		in.Password = `chpass`
	}
	if in.Database == `` {
		in.Database = `chdb`
	}
}

func (in *ChDockerTest) ConnectCheck(res *dockertest.Resource) (driver.Conn, error) {
	in.Port = res.GetPort("9000/tcp")
	hostPort := in.pool.HostPort(in.Port)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{hostPort},
		Auth: clickhouse.Auth{
			Database: in.Database,
			Username: in.User,
			Password: in.Password,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		Settings: clickhouse.Settings{
			`max_execution_time`: 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:           true,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	return conn, err
}
