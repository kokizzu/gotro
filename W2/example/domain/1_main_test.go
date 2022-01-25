package domain

import (
	"github.com/goccy/go-yaml"
	"github.com/kokizzu/goproc"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
	"github.com/kpango/fastime"
	"github.com/ory/dockertest/v3"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

type DockerComposeConf struct { // https://zhwt.github.io/yaml-to-go/
	Version  string `yaml:"version"`
	Services struct {
		Clickhouse1 struct {
			Image   string   `yaml:"image"`
			Ports   []string `yaml:"ports"`
			Ulimits struct {
				Nofile struct {
					Soft int `yaml:"soft"`
					Hard int `yaml:"hard"`
				} `yaml:"nofile"`
			} `yaml:"ulimits"`
		} `yaml:"clickhouse1"`
		Tarantool1 struct {
			Image       string `yaml:"image"`
			Environment struct {
				TARANTOOLUSERNAME     string `yaml:"TARANTOOL_USER_NAME"`
				TARANTOOLUSERPASSWORD string `yaml:"TARANTOOL_USER_PASSWORD"`
			} `yaml:"environment"`
			Volumes []string `yaml:"volumes"`
			Ports   []string `yaml:"ports"`
		} `yaml:"tarantool1"`
		Mailhog struct {
			Image string   `yaml:"image"`
			Ports []string `yaml:"ports"`
		} `yaml:"mailhog"`
	} `yaml:"services"`
}

func ParseDockerCompose() *DockerComposeConf {
	tryDirs := []string{`./`, `../`, `../../`}
	const fileName = `docker-compose.yml`
	var confFile string
	for _, tryDir := range tryDirs {
		confFile = tryDir + fileName
		bytes, err := ioutil.ReadFile(confFile)
		if err != nil {
			continue
		}
		// parse config
		dcConf := DockerComposeConf{}
		err = yaml.Unmarshal(bytes, &dcConf)
		if err != nil {
			log.Fatalf(`failed yaml.Unmarshal %s: %s`, confFile, err.Error())
		}
		return &dcConf
	}
	log.Fatalf(`failed read file %s`, fileName)
	return nil
}

func ConnectLocalDocker() *dockertest.Pool {
	endpoint := `http://localhost:2375`
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(`127.0.0.1`, `2375`), time.Second)
	if conn != nil {
		conn.Close()
	} else {
		endpoint = ``

		// intellij does not load env properly on run
		daemon := goproc.New()
		cmdId := daemon.AddCommand(&goproc.Cmd{
			Program:    `docker-machine`, // program to run
			Parameters: []string{`env`},  // command line arguments
			HideStdout: true,
			OnStdout: func(cmd *goproc.Cmd, line string) error { // optional
				const export = `export `
				if !S.StartsWith(line, export) {
					return nil
				}
				line = line[len(export):] // remove prefix
				conf := S.Split(line, `="`)
				L.Print(`setting env: `, line)
				v := conf[1]
				os.Setenv(conf[0], v[:len(v)-1])
				return nil
			},
		})
		daemon.Start(cmdId)
	}
	pool, err := dockertest.NewPool(endpoint)
	if err != nil {
		log.Fatalf("Could not connect to docker service: %s | %s | %s", err,
			os.Getenv(`DOCKER_HOST`),
			os.Getenv(`DOCKER_CERT_PATH`),
		)
	}
	return pool
}

func SpawnDocker(pool *dockertest.Pool, imageVer string, retryFunc func(res *dockertest.Resource) func() error) (cleaner func()) {
	img := strings.Split(imageVer, `:`) // image:version
	if len(img) != 2 {
		panic(`missing version on ` + img[0])
	}
	res, err := pool.Run(img[0], img[1], []string{})
	L.PanicIf(err, `pool.Run %v`, img)
	cleaner = func() {
		L.IsError(pool.Purge(res), `failed purge resource`, img)
	}
	if err := pool.Retry(retryFunc(res)); err != nil {
		log.Fatalf("Could not connect to %v docker port: %s ", img, err)
	}
	return
}

func TestMain(t *testing.M) {
	docker := ParseDockerCompose()
	L.Print(conf.LoadTestEnv())

	pool := ConnectLocalDocker()
	var code int
	defer func() {
		err := recover()
		if err != nil {
			L.Print(err, `catched panic`)
		}
		os.Exit(code)
	}()

	// try spawn and connect tarantool
	ttCleaner := SpawnDocker(pool, docker.Services.Tarantool1.Image, func(res *dockertest.Resource) func() error {
		return func() (err error) {

			defer func() {
				rec := recover()
				err, _ = rec.(error)
				L.Print(rec, err)
			}()
			conf.TARANTOOL_PORT = res.GetPort("3301/tcp")

			// override conf to use guest
			conf.TARANTOOL_USER = `guest`
			conf.TARANTOOL_PASS = ``

			L.Print(`connecting`, conf.TARANTOOL_HOST, conf.TARANTOOL_PORT)
			conf.ConnectTarantool()
			return
		}
	})
	defer ttCleaner()

	// try spawn and connect clickhouse
	chCleaner := SpawnDocker(pool, docker.Services.Clickhouse1.Image, func(res *dockertest.Resource) func() error {
		return func() (err error) {
			defer func() {
				rec := recover()
				err, _ = rec.(error)
				L.Print(rec, err)
			}()
			conf.CLICKHOUSE_PORT = res.GetPort("9000/tcp")
			L.Print(`connecting`, conf.CLICKHOUSE_HOST, conf.CLICKHOUSE_PORT)
			conf.ConnectClickhouse()
			return
		}
	})
	defer chCleaner()

	// try spawn and connect mailhog
	mhCleaner := SpawnDocker(pool, docker.Services.Mailhog.Image, func(res *dockertest.Resource) func() error {
		return func() (err error) {
			port := res.GetPort("1025/tcp")
			conf.MAILER_PORT = S.ToInt(port)
			L.Print(`connecting`, conf.MAILER_HOST, conf.MAILER_PORT)
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(DockerIp(), port), time.Second)
			if err == nil {
				conn.Close()
			}
			return err
		}
	})
	defer mhCleaner()

	// run migration first
	model.RunMigration()

	// ensure admin account exists
	d := NewDomain()
	user := wcAuth.NewUsersMutator(d.Taran)
	user.Id = 1
	if !user.FindById() {
		user.Email = conf.SuperAdmin
		user.SetEncryptPassword(user.Email)
		user.CreatedAt = fastime.UnixNow()
		if !user.DoUpdateById() {
			panic(`cannot create superadmin`)
		}
	}

	code = t.Run()
}

func DockerIp() string {
	dockerHost := os.Getenv(`DOCKER_HOST`)
	if dockerHost == `` {
		return `127.0.0.1`
	}
	dockerHost = S.RightOfLast(dockerHost, `/`)
	dockerHost = S.LeftOf(dockerHost, `:`)
	return dockerHost
}
