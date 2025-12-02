package conf

import (
	"os"

	"github.com/kokizzu/gotro/S"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OauthConf struct {
	Urls         []string
	GoogleScope  []string
	TwitterScope []string

	Google  map[string]*oauth2.Config
	Twitter map[string]*oauth2.Config
}

func EnvOauth() (res OauthConf) {
	res.Urls = S.Split(os.Getenv(`OAUTH_URLS`), `,`)

	res.Google = map[string]*oauth2.Config{}
	res.GoogleScope = S.Split(os.Getenv(`OAUTH_GOOGLE_SCOPES`), `,`)

	res.Twitter = map[string]*oauth2.Config{}
	res.TwitterScope = S.Split(os.Getenv(`OAUTH_TWITTER_SCOPES`), `,`)
	for _, url := range res.Urls {
		res.Google[url] = &oauth2.Config{
			ClientID:     os.Getenv(`OAUTH_GOOGLE_CLIENT_ID`),
			ClientSecret: os.Getenv(`OAUTH_GOOGLE_CLIENT_SECRET`),
			RedirectURL:  url + `/guest/oauthCallback`,
			Scopes:       res.GoogleScope,
			Endpoint:     google.Endpoint,
		}

		res.Twitter[url] = &oauth2.Config{
			ClientID:     os.Getenv(`OAUTH_TWITTER_CLIENT_ID`),
			ClientSecret: os.Getenv(`OAUTH_TWITTER_CLIENT_SECRET`),
			RedirectURL:  url + `/guest/oauthCallback`,
			Scopes:       res.TwitterScope,
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://twitter.com/i/oauth2/authorize",
				TokenURL:  "https://api.twitter.com/2/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		}
	}
	return
}
