package conf

import (
	"fmt"
	"os"

	"github.com/kokizzu/gotro/S"
)

type WebConf struct {
	Port           int
	WebProtoDomain string
}

func EnvWebConf() WebConf {
	webProtoDomain := os.Getenv("WEB_PROTO_DOMAIN")
	if webProtoDomain == `` {
		webProtoDomain = `http://localhost:1235`
	}

	return WebConf{
		Port:           S.ToInt(os.Getenv("WEB_PORT")),
		WebProtoDomain: webProtoDomain,
	}
}

func (w WebConf) ListenAddr() string {
	return ":" + fmt.Sprint(w.Port)
}
