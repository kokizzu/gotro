package conf

import (
	"errors"
	"os"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
)

const PROJECT_NAME = `example` // must be the same as go.mod
const API_PREFIX = `/api`
const SKIN_SUBDIR = `skins/`
const EVENT_SUBDIR = `events/`
const COIN_NAME = `$DIH Coin`
const MAIL_VIEWS_DIR = `/views`
const AVATAR_PREFIX = `/avatars/`

var (
	// ALL_CAPS = from .env
	// split per service when microservice
	TARANTOOL_HOST string
	TARANTOOL_PORT string
	TARANTOOL_USER string
	TARANTOOL_PASS string
	ROOT_DIR       string
	// TODO: add more configs

	CLICKHOUSE_HOST string
	CLICKHOUSE_PORT string
	CLICKHOUSE_USER string
	CLICKHOUSE_PASS string

	LIGHTSTEP_ACCESS_TOKEN string

	WEBAPI_HOSTPORT string
	WEBAPI_EXEPATH  string

	MINIO_ENDPOINT   string
	MINIO_ACCESS_KEY string
	MINIO_SECRET_KEY string
	MINIO_USESSL     bool

	MAILER_HOST string
	MAILER_PORT int
	MAILER_USER string
	MAILER_PASS string

	// CameCase = from command line
	SERVICE_MODE string
	ENV          string
	DEBUG_MODE   bool

	UPLOAD_DIR string
	UPLOAD_URI string
)

func LoadFromEnv(ignoreBinary ...interface{}) {
	// TODO: change to loop from struct's tag and print the result
	TARANTOOL_HOST = S.IfEmpty(os.Getenv(`TARANTOOL_HOST`), `127.0.0.1`)
	TARANTOOL_PORT = S.IfEmpty(os.Getenv(`TARANTOOL_PORT`), `3301`)
	TARANTOOL_USER = S.IfEmpty(os.Getenv(`TARANTOOL_USER`), `guest`)
	TARANTOOL_PASS = os.Getenv(`TARANTOOL_PASS`)

	CLICKHOUSE_HOST = S.IfEmpty(os.Getenv(`CLICKHOUSE_HOST`), `127.0.0.1`)
	CLICKHOUSE_PORT = S.IfEmpty(os.Getenv(`CLICKHOUSE_PORT`), `9000`)
	CLICKHOUSE_USER = os.Getenv(`CLICKHOUSE_USER`)
	CLICKHOUSE_PASS = os.Getenv(`CLICKHOUSE_PASS`)

	LIGHTSTEP_ACCESS_TOKEN = os.Getenv(`LIGHTSTEP_ACCESS_TOKEN`)

	WEBAPI_HOSTPORT = S.IfEmpty(os.Getenv(`WEBAPI_HOSTPORT`), `:3000`)
	WEBAPI_EXEPATH = S.IfEmpty(os.Getenv(`WEBAPI_EXEPATH`), `./`+PROJECT_NAME+`.exe`)

	MINIO_ENDPOINT = os.Getenv(`MINIO_ENDPOINT`)
	MINIO_ACCESS_KEY = os.Getenv(`MINIO_ACCESS_KEY`)
	MINIO_SECRET_KEY = os.Getenv(`MINIO_SECRET_KEY`)
	MINIO_USESSL = os.Getenv(`MINIO_USESSL`) == `true`

	MAILER_HOST = os.Getenv(`MAILER_HOST`)
	MAILER_PORT = S.ToInt(os.Getenv(`MAILER_PORT`))
	MAILER_USER = os.Getenv(`MAILER_USER`)
	MAILER_PASS = os.Getenv(`MAILER_PASS`)

	pwd, _ := os.Getwd()
	UPLOAD_DIR = pwd + S.IfEmpty(os.Getenv(`UPLOAD_DIR`), `/svelte/dist/upload/`)
	UPLOAD_URI = S.IfEmpty(os.Getenv(`UPLOAD_URI`), `/upload/`)

	ENV = S.IfEmpty(os.Getenv(`ENV`), `dev`)
	if ENV == `dev` {
		DEBUG_MODE = true
	}

	if len(ignoreBinary) == 0 {
		if !L.FileExists(WEBAPI_EXEPATH) {
			L.PanicIf(errors.New(`binary must be exists`), WEBAPI_EXEPATH)
		}
	}
}

// hard dependency does not need to return error, panic is ok
// TODO: change L.PanicIf to use zerolog?
