package domain

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W2/example/conf"
)

type (
	UserExternalLogin_In struct {
		RequestCommon
		Provider string
	}
	UserExternalLogin_Out struct {
		ResponseCommon
		Link string
	}
)

const UserExternalLogin_Url = `/UserExternalLogin`

func (d *Domain) UserExternalLogin(in *UserExternalLogin_In) (out UserExternalLogin_Out) {
	switch in.Provider {
	case `google`:
		gProvider := conf.GPLUS_OAUTH_PROVIDERS[in.Host]
		if gProvider == nil {
			out.SetError(500, `host not configured with oauth: `+in.Host)
			return
		}
		out.Link = gProvider.AuthCodeURL(in.Provider)
		fmt.Println(out.Link)
	case `yahoo`:
		gProvider := conf.YAHOO_OAUTH_PROVIDERS[in.Host]
		if gProvider == nil {
			out.SetError(500, `host not configured with oauth: `+in.Host)
			return
		}
		out.Link = gProvider.AuthCodeURL(in.Provider)
		fmt.Println(out.Link)
	default:
		out.SetError(400, `provider not set: `+in.Provider)
	}
	return
}

func fetchJson(client *http.Client, url string, res *ResponseCommon) (json M.SX) {
	resp, err := client.Get(url)
	if L.IsError(err, `failed fetch url %s`, url) {
		res.SetError(500, `failed fetch url`)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if L.IsError(err, `failed read body`) {
		res.SetError(500, `failed read body`)
		return
	}
	bodyStr := string(body)
	json = S.JsonToMap(bodyStr)
	L.Describe(json)
	err2 := json.GetStr(`error`)
	if L.CheckIf(err2 != ``, `fetchJson %#v`, json) {
		res.SetError(500, err2)
		return
	}
	err3 := json.GetStr(`type`)
	if L.CheckIf(err3 == `OAuthException`, `fetchJson %#v`, json) {
		res.SetError(500, err3)
		return
	}
	return
}

type (
	UserOauth_In struct {
		RequestCommon
		State string
		Code  string
	}
	UserOauth_Out struct {
		ResponseCommon
		Dummy interface{}
	}
)

const UserOauth_Url = `/UserOauth`

func (d *Domain) UserOauth(in *UserOauth_In) (out UserOauth_Out) {
	switch in.State {
	case `google`:
		gProvider := conf.GPLUS_OAUTH_PROVIDERS[in.Host]
		if gProvider == nil {
			out.SetError(500, `host not configured with oauth: `+in.Host)
			return
		}
		token, err := gProvider.Exchange(in.TracerContext, in.Code)
		if err != nil {
			out.SetError(500, `failed exchange oauth token`)
			return
		}
		client := gProvider.Client(in.TracerContext, token)
		if conf.GPLUS_USERINFO_ENDPOINT == `` {
			// no need to refetch userinfo_endpoint
			json := fetchJson(client, `https://accounts.google.com/.well-known/openid-configuration`, &out.ResponseCommon)
			conf.GPLUS_USERINFO_ENDPOINT = json.GetStr(`userinfo_endpoint`)
		}
		out.Dummy = fetchJson(client, conf.GPLUS_USERINFO_ENDPOINT, &out.ResponseCommon)
		// example: {"email":"","email_verified":true,"family_name":"","gender":"","given_name":"","locale":"en-GB","name":"","picture":"http://","profile":"http://","sub":"number"};
	case `yahoo`:
		gProvider := conf.YAHOO_OAUTH_PROVIDERS[in.Host]
		if gProvider == nil {
			out.SetError(500, `host not configured with oauth: `+in.Host)
			return
		}
		token, err := gProvider.Exchange(in.TracerContext, in.Code)
		if err != nil {
			out.SetError(500, `failed exchange oauth token`)
			return
		}
		client := gProvider.Client(in.TracerContext, token)
		if conf.YAHOO_USERINFO_ENDPOINT == `` {
			// no need to refetch userinfo_endpoint
			json := fetchJson(client, `https://api.login.yahoo.com/openid/v1/userinfo`, &out.ResponseCommon)
			conf.YAHOO_USERINFO_ENDPOINT = json.GetStr(`userinfo_endpoint`)
		}
		out.Dummy = fetchJson(client, conf.YAHOO_USERINFO_ENDPOINT, &out.ResponseCommon)
		// example: {"email":"","email_verified":true,"family_name":"","gender":"","given_name":"","locale":"en-GB","name":"","picture":"http://","profile":"http://","sub":"number"};
	default:
		out.SetError(400, `provider not set: `+in.State)
	}
	return
}
