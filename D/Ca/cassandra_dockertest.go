package Ca

import (
	"net"
	"time"

	"github.com/kokizzu/gotro/D"
	"github.com/ory/dockertest/v3"
)

type CaDockerTest struct {
	Image string
	Port  string
	pool  *D.DockerTest
}

// https://hub.docker.com/_/cassandra
func (in *CaDockerTest) ImageVersion(pool *D.DockerTest, version string) *dockertest.RunOptions {
	in.pool = pool
	in.SetDefaults(version)
	return &dockertest.RunOptions{
		Repository: `cassandra`,
		Name:       `dockertest-cassandra-` + pool.Uniq,
		Tag:        in.Image,
		NetworkID:  pool.Network.ID,
		Env: []string{
			`JVM_EXTRA_OPTS=-Dcassandra.skip_wait_for_gossip_to_settle=0 -Dcassandra.load_ring_state=false -Dcassandra.initial_token=1 -Dcassandra.num_tokens=nil -Dcassandra.allocate_tokens_for_local_replication_factor=nil`,
			`CASSANDRA_BROADCAST_ADDRESS=`,
			`CASSANDRA_CLUSTER_NAME=cluster1`,
			`CASSANDRA_DC=dc1`,
			`CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch=`,
		},
	}
}

func (in *CaDockerTest) Image3(pool *D.DockerTest) *dockertest.RunOptions {
	return in.ImageVersion(pool, `3`)
}

func (in *CaDockerTest) SetDefaults(img string) {
	if in.Image == `` {
		in.Image = img
	}
}

func (in *CaDockerTest) ConnectCheck(res *dockertest.Resource) (err error) {
	in.Port = res.GetPort("9042/tcp")
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
