package conf

import (
	"github.com/joho/godotenv"
	"github.com/kokizzu/gotro/L"
)

var VERSION = ``

const PROJECT_NAME = `Example2`

func IsDebug() bool {
	return VERSION == ``
}

func LoadEnv() {
	dirRetryList := []string{``, `../`, `../../`, `../../../`}
	for _, dirPrefix := range dirRetryList {
		envFile := dirPrefix + `.env`
		err := godotenv.Overload(envFile)
		if err == nil {
			envOverrideFile := dirPrefix + `.env.override`
			err = godotenv.Overload(envOverrideFile)
			L.PanicIf(err, `godotenv.Load .env.override`)
			return
		}
	}
	panic(`cannot load .env and .env.override`)
}
