package main

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model"

	"github.com/joho/godotenv"
	//"github.com/lightstep/otel-launcher-go/launcher"
)

var VERSION = ``
var log *zerolog.Logger

func main() {
	conf.VERSION = VERSION
	log = conf.InitLogger()

	err := godotenv.Overload(`.env`)
	L.PanicIf(err, `godotenv.Load .env`)
	err = godotenv.Overload(`.env.override`)
	L.PanicIf(err, `godotenv.Load .env.override`)

	args := os.Args
	if len(args) < 2 {
		L.Print(`must start with: run, web/rest, grpc, or cron as first argument`)
		return
	}
	cliMode := args[1] == `run` || args[1] == `cmd` || args[1] == `cli`

	conf.SERVICE_MODE = B.ToS(cliMode)
	conf.LoadFromEnv(cliMode)

	model.RunMigration()

	//if !cliMode {
	//	// telemetry
	//	//ls := launcher.ConfigureOpentelemetry(
	//	//	launcher.WithAccessToken(conf.LIGHTSTEP_ACCESS_TOKEN),
	//	//	launcher.WithServiceName(conf.PROJECT_NAME),
	//	//	launcher.WithServiceVersion("v0.0.1"), // TODO: get from command line $(git rev-parse --verify HEAD)
	//	//	launcher.WithResourceAttributes(map[string]string{
	//	//		"service.mode":           conf.SERVICE_MODE,
	//	//		"deployment.environment": conf.ENV,
	//	//	}),
	//	//)
	//	//defer ls.Shutdown()
	//
	//	// gops
	//	err := agent.Listen(agent.Options{})
	//	L.IsError(err, `gops agent.Listen`)
	//}
	//
	//switch args[1] {
	//case `run`, `cli`, `cmd`:
	//	// start command line mode
	//	if len(args) < 2 {
	//		L.Print(`2nd parameter must be the url segments domainRole/action`)
	//		return
	//	}
	//	cliArgsRunner(args)
	//	// 2nd param: role/action
	//	// input/output always from stdin/stdout
	//	// all logs should be to stderr
	//	// other args:
	//	//  --auth=cookie,jwt,impersonation-key
	//	//  --in-format=json,queryString
	//	//  --upload=key=filename
	//	//  --out-format=json,yaml,table,csv
	//	// TODO: codegen the available commands/options from d*
	//case `web`, `rest`:
	//	// start web mode (also websocket)
	//	// TODO: overwrite logger to json format
	//	NewWebApi()
	//case `grpc`:
	//	// TODO: overwrite logger to json format
	//	// start grpc server, use twitch's twirp instead of google's
	//case `cron`:
	//	d := domain.NewDomain()
	//	runCron(d)
	//	d.WaitInterrupt()
	//	// TODO: overwrite logger to json format
	//	// start cron runner if 1st param = cron
	//	// must set 1 tarantool key with autoexpire and defer unset
	//	// if only 1 instance may run at the same time
	//case `help`:
	//	// TODO: generate list of all arguments and their input/output
	//}
}
