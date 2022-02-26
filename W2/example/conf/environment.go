package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yahoo"
)

const PROJECT_NAME = `github.com/kokizzu/gotro/W2/example` // must be the same as go.mod
const API_PREFIX = `/api`
const MEDIA_SUBDIR = `media/`
const MAIL_VIEWS_DIR = `/views`

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

	DOCKER_TLS_VERIFY   string
	DOCKER_HOST         string
	DOCKER_CERT_PATH    string
	DOCKER_MACHINE_NAME string

	OAUTH_URLS          []string
	OAUTH_CALLBACK_PATH string

	GPLUS_CLIENTID     string
	GPLUS_CLIENTSECRET string
	GPLUS_SCOPES       []string

	YAHOO_APPID        string
	YAHOO_CLIENTID     string
	YAHOO_CLIENTSECRET string
	YAHOO_SCOPES       []string

	GITHUB_CLIENTID     string
	GITHUB_CLIENTSECRET string
	GITHUB_SCOPES       []string

	STEAM_APPID        string
	STEAM_CLIENTSECRET string
	STEAM_ENDPOINT     oauth2.Endpoint

	TWITTER_CLIENTID     string
	TWITTER_CLIENTSECRET string
	TWITTER_SCOPES       []string

	FACEBOOK_APPID     string
	FACEBOOK_APPSECRET string
	FACEBOOK_SCOPES    []string

	GPLUS_OAUTH_PROVIDERS    map[string]*oauth2.Config
	YAHOO_OAUTH_PROVIDERS    map[string]*oauth2.Config
	GITHUB_OAUTH_PROVIDERS   map[string]*oauth2.Config
	STEAM_OAUTH_PROVIDERS    map[string]*oauth2.Config
	TWITTER_OAUTH_PROVIDERS  map[string]*oauth2.Config
	FACEBOOK_OAUTH_PROVIDERS map[string]*oauth2.Config

	GPLUS_USERINFO_ENDPOINT string
)

func strArr(envName string, separator string) []string {
	str := os.Getenv(envName)
	return S.Split(str, separator)
}

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

	// OpenID/OAuth
	OAUTH_URLS = strArr(`OAUTH_URLS`, `,`)
	OAUTH_CALLBACK_PATH = os.Getenv(`OAUTH_CALLBACK_PATH`)

	GPLUS_CLIENTID = os.Getenv(`GPLUS_CLIENTID`)
	GPLUS_CLIENTSECRET = os.Getenv(`GPLUS_CLIENTSECRET`)
	GPLUS_SCOPES = strArr(`GPLUS_SCOPES`, `,`)

	GPLUS_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		GPLUS_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     GPLUS_CLIENTID,
			ClientSecret: GPLUS_CLIENTSECRET,
			RedirectURL:  url + OAUTH_CALLBACK_PATH,
			Scopes:       GPLUS_SCOPES,
			Endpoint:     google.Endpoint,
		}
	}

	YAHOO_APPID = os.Getenv(`YAHOO_APPID`)
	YAHOO_CLIENTID = os.Getenv(`YAHOO_CLIENTID`)
	YAHOO_CLIENTSECRET = os.Getenv(`YAHOO_CLIENTSECRET`)
	YAHOO_SCOPES = strArr(`YAHOO_SCOPES`, `,`)

	YAHOO_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		YAHOO_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     YAHOO_CLIENTID,
			ClientSecret: YAHOO_CLIENTSECRET,
			RedirectURL:  url + OAUTH_CALLBACK_PATH,
			Scopes:       YAHOO_SCOPES,
			Endpoint:     yahoo.Endpoint,
		}
	}

	GITHUB_CLIENTID = os.Getenv(`GITHUB_CLIENTID`)
	GITHUB_CLIENTSECRET = os.Getenv(`GITHUB_CLIENTSECRET`)
	GITHUB_SCOPES = strArr(`GITHUB_SCOPES`, `,`)

	GITHUB_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		GITHUB_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     GITHUB_CLIENTID,
			ClientSecret: GITHUB_CLIENTSECRET,
			RedirectURL:  url + OAUTH_CALLBACK_PATH,
			Scopes:       GITHUB_SCOPES,
			Endpoint:     github.Endpoint,
		}
	}

	STEAM_APPID = os.Getenv(`STEAM_APPID`)
	STEAM_CLIENTSECRET = os.Getenv(`STEAM_CLIENTSECRET`)
	STEAM_ENDPOINT = oauth2.Endpoint{
		AuthURL: fmt.Sprintf(
			"https://steamcommunity.com/oauth/login?response_type=token&client_id=%s",
			STEAM_APPID),
		TokenURL: fmt.Sprintf(
			"http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s",
			STEAM_CLIENTSECRET, STEAM_APPID),
	}

	STEAM_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		STEAM_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     STEAM_APPID,
			ClientSecret: STEAM_CLIENTSECRET,
			RedirectURL:  url + OAUTH_CALLBACK_PATH,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://steamcommunity.com/oauth/login",                           // ?response_type=token&client_id=client_id_here&state=whatever_you_want
				TokenURL: "https://api.steampowered.com/ISteamUserOAuth/GetTokenDetails/v1/", // ?access_token=token
			},
		}
	}

	TWITTER_CLIENTID = os.Getenv(`TWITTER_CLIENTID`)
	TWITTER_CLIENTSECRET = os.Getenv(`TWITTER_CLIENTSECRET`)
	TWITTER_SCOPES = strArr(`TWITTER_SCOPES`, `,`)

	TWITTER_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		TWITTER_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     TWITTER_CLIENTID,
			ClientSecret: TWITTER_CLIENTSECRET,
			RedirectURL:  url + OAUTH_CALLBACK_PATH,
			Scopes:       TWITTER_SCOPES,
			Endpoint: oauth2.Endpoint{
				AuthURL: "https://twitter.com/i/oauth2/authorize?code_challenge_method=plain&code_challenge=CODE_CHALLENGE",
				// the rest not needed, since twitter oauth are not normal
				// we exchange manually
			},
		}
	}

	STEAM_CLIENTSECRET = os.Getenv(`STEAM_CLIENTSECRET`)

	STEAM_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		STEAM_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     "", // TODO: fill these
			ClientSecret: STEAM_CLIENTSECRET,
			Endpoint: oauth2.Endpoint{
				AuthURL: "https://steamcommunity.com/oauth/login",
				// the rest not needed, since steam oauth give access_token directly
			},
			RedirectURL: "",
			Scopes:      nil,
		}
	}

	FACEBOOK_APPID = os.Getenv(`FACEBOOK_APPID`)
	FACEBOOK_APPSECRET = os.Getenv(`FACEBOOK_APPSECRET`)
	FACEBOOK_SCOPES = strArr(`FACEBOOK_SCOPES`, `,`)

	FACEBOOK_OAUTH_PROVIDERS = map[string]*oauth2.Config{}
	for _, url := range OAUTH_URLS {
		FACEBOOK_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     FACEBOOK_APPID,
			ClientSecret: FACEBOOK_APPSECRET,
			RedirectURL:  url + OAUTH_CALLBACK_PATH,
			Scopes:       FACEBOOK_SCOPES,
			Endpoint:     facebook.Endpoint,
		}
	}

	if len(ignoreBinary) == 0 {
		if !L.FileExists(WEBAPI_EXEPATH) {
			L.PanicIf(errors.New(`binary must be exists`), WEBAPI_EXEPATH)
		}
	}
}

// hard dependency does not need to return error, panic is ok

// loads .env file even when the binary/test not in project's root directory
// returns project's root directory (where `.env` should be located)
func LoadTestEnv() string {
	for z := 0; z < 4; z++ {
		dir := strings.Repeat(`../`, z)
		err := godotenv.Load(dir + `.env`)
		if err == nil {
			cwd, _ := os.Getwd()
			for i := 0; i < z; i++ {
				cwd = S.LeftOfLast(cwd, "/")
			}
			LoadFromEnv(true)
			return cwd
		}
	}
	return ``
}

// TODO: change L.PanicIf to use zerolog?
