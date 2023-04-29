package D

import (
	"log"
	"net"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/kokizzu/gotro/L"
)

type DockerTest struct {
	Pool      *dockertest.Pool
	Network   *docker.Network
	Uniq      string
	Resources []*dockertest.Resource
}

func InitDockerTest(endpoint string) *DockerTest {
	dockerUniq := time.Now().Format("20060102-150405")
	pool, err := dockertest.NewPool(endpoint)
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	network, err := pool.Client.CreateNetwork(docker.CreateNetworkOptions{
		Name: `dockertest-` + dockerUniq,
	})
	if err != nil {
		log.Fatalf("Could not create docker network: %s", err)
	}
	return &DockerTest{
		Pool:    pool,
		Network: network,
		Uniq:    dockerUniq,
	}
}

func (d *DockerTest) Cleanup() {
	for _, res := range d.Resources {
		L.IsError(d.Pool.Purge(res), `failed Purge: `+res.Container.Name)
	}
	netId := d.Network.ID
	L.IsError(d.Pool.Client.RemoveNetwork(netId), `failed RemoveNetwork: `+netId)
}

func (d *DockerTest) Spawn(options *dockertest.RunOptions, checkFunc func(res *dockertest.Resource) error) {
	res, err := d.Pool.RunWithOptions(options)
	label := options.Repository + `:` + options.Tag
	L.PanicIf(err, `failed RunWithOptions: `+label)
	err = d.Pool.Retry(func() error {
		return checkFunc(res)
	})
	L.PanicIf(err, `failed Pool.Retry: `+label)
	d.Resources = append(d.Resources, res)
}

func (d *DockerTest) HostPort(port string) string {
	// TODO: fetch remote docker if any (eg. docker machine env for mac/windows)
	return net.JoinHostPort(`127.0.0.1`, port)
}
