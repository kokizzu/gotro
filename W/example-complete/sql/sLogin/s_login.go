package sLogin

import (
	"errors"
	"example-complete/sql"
	"example-complete/sql/sResponse"
	"example-complete/sql/tUsers"
	"fmt"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"net/url"
)

// TODO: kalau pakai Firebase Auth, ntar import https://github.com/robbert229/jwt

// TODO: change to correct one (see console.developer.google.com

var OAUTH_URLS = []string{
	`http://TODO_CHANGE_DOMAIN`,
	`https://TODO_CHANGE_DOMAIN`,
	`http://www.TODO_CHANGE_DOMAIN`,
	`https://www.TODO_CHANGE_DOMAIN`,
	`http://local.TODO_CHANGE_DOMAIN`,
}
var GPLUS_OAUTH_PROVIDERS map[string]*oauth2.Config
var USERINFO_ENDPOINT string

const GPLUS_CLIENTID = `1060129198780-q0eco7m81bsj0ip9n119ukm72ebntiv1.apps.googleusercontent.com`
const GPLUS_CLIENTSECRET = `P4wEX6PZgL0q9wGt-UrIOjWO`

const RESET_MINUTE = 20

var Z func(string) string
var ZZ func(string) string
var ZJ func(string) string
var ZB func(bool) string
var ZI func(int64) string
var ZLIKE func(string) string
var ZT func(strs ...string) string
var PG *Pg.RDBMS

func init() {
	Z = S.Z
	ZB = S.ZB
	ZZ = S.ZZ
	ZJ = S.ZJJ
	ZI = S.ZI
	ZT = S.ZT
	ZLIKE = S.ZLIKE
	PG = sql.PG
}

// credential for OpenID
// https://console.developers.google.com/apis/credentials?project=example-complete-example-cron

type fbConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// 2016-11-23 Prayogo
func fetchJson(url string) (W.Ajax, error) {
	ajax := W.NewAjax()
	L.Print(url)
	resp, err := http.Get(url)
	if ajax.ErrorIf(err, sql.ERR_201_FAILED_OAUTH_EXCHANGE) {
		return ajax, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if ajax.ErrorIf(err, sql.ERR_201_FAILED_OAUTH_EXCHANGE) {
		return ajax, err
	}
	body_str := string(body)
	m := S.JsonToMap(body_str)
	err2 := X.ToS(m[`error`])
	if L.CheckIf(err2 != ``, `fetchJson %# v`, m) {
		return ajax, fmt.Errorf(`%s`, err2)
	}
	err3 := X.ToS(m[`type`])
	if L.CheckIf(err3 == `OAuthException`, `fetchJson %# v`, m) {
		return ajax, fmt.Errorf(`%s`, body_str)
	}
	//L.Describe(m)
	//L.Describe(string(body))
	return W.Ajax{M.SX(m)}, nil
}

// 2016-01-10 Prayogo
func RetrieveGoogleUserInfo(provider *oauth2.Config, access_token *oauth2.Token) (res W.Ajax, err error) {
	res = W.NewAjax()
	client := provider.Client(oauth2.NoContext, access_token)
	if USERINFO_ENDPOINT == `` {
		// no need to refetch userinfo_endpoint
		response, err := client.Get(`https://accounts.google.com/.well-known/openid-configuration`)
		if err != nil {
			return res, err
		}
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return res, err
		}
		json_body := S.JsonToMap(string(body))
		USERINFO_ENDPOINT = X.ToS(json_body[`userinfo_endpoint`])
	}
	response, err := client.Get(USERINFO_ENDPOINT)
	if err != nil {
		return res, err
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return res, err
	}
	json := S.JsonToMap(string(body)) // example: {"email":"","email_verified":true,"family_name":"","gender":"","given_name":"","locale":"en-GB","name":"","picture":"http://","profile":"http://","sub":"number"};
	return W.Ajax{json}, nil
}

// 2016-11-23 Prayogo
func (f *fbConfig) AuthCodeURL(state string) string {
	url2, err := url.Parse(`https://www.facebook.com/v2.8/dialog/oauth`)
	L.PanicIf(err, sql.ERR_201_FAILED_OAUTH_EXCHANGE)
	parameters := url.Values{}
	parameters.Add(`display`, `page`)
	parameters.Add(`client_id`, f.ClientID)
	parameters.Add(`redirect_uri`, f.RedirectURL)
	parameters.Add(`scope`, A.StrJoin(f.Scopes, `,`))
	parameters.Add(`state`, state)
	url2.RawQuery = parameters.Encode()
	url1, err := url.Parse(`https://www.facebook.com/login.php`)
	L.PanicIf(err, sql.ERR_201_FAILED_OAUTH_EXCHANGE)
	parameters = url.Values{}
	parameters.Add(`skip_api_login`, `1`)
	parameters.Add(`api_key`, f.ClientID)
	parameters.Add(`signed_next`, `1`)
	parameters.Add(`next`, url2.String())
	url1.RawQuery = parameters.Encode()
	return url1.String()
}

func init() {
	GPLUS_OAUTH_PROVIDERS = map[string]*oauth2.Config{} // yahoo tidak bisa multiple domain (harus dibuat 1-1), tidak support IP
	for _, url := range OAUTH_URLS {
		GPLUS_OAUTH_PROVIDERS[url] = &oauth2.Config{
			ClientID:     GPLUS_CLIENTID,
			ClientSecret: GPLUS_CLIENTSECRET,
			RedirectURL:  url + `/login/verify`,
			Scopes: []string{
				`openid`,
				`email`,
				`profile`,
			},
			Endpoint: google.Endpoint,
		}
	}
}

// tutorial: http://golangtutorials.blogspot.com/2011/11/oauth2-3-legged-authorization-in-go-web.html
// https://developers.google.com/identity/protocols/OpenIDConnect
// get G+ OAuth provider and domain csrf
// 2016-11-08 Prayogo
func GetGPlusOAuth(ctx *W.Context) *oauth2.Config {
	return GPLUS_OAUTH_PROVIDERS[ctx.Host()]
}

// handle G+ oauth login
// 2016-07-26 Prayogo
func GPlusExchangeInfo(provider *oauth2.Config, gets *W.QueryParams) (W.Ajax, error) {
	token, err := provider.Exchange(oauth2.NoContext, gets.GetStr(`code`))
	if err != nil {
		return W.NewAjax(), err
	}
	return RetrieveGoogleUserInfo(provider, token)
}

func API_All_Logout(ctx *W.Context) {
	ajax := sResponse.NewAjax()
	user_id := ctx.Session.GetInt(`user_id`)
	if user_id > 0 && !ctx.IsWebMaster() {
		// TODO: update last login
	}
	ctx.Session.Logout()
	ctx.AppendJson(ajax.SX)
}

func AccessLevel(email string, id int64) M.SX {
	query := `SELECT COALESCE((
			SELECT group_id
			FROM users 
			WHERE id = ` + ZI(id) + `
		),0)`
	group_id := PG.QInt(query)
	query2 := `SELECT COALESCE((
			SELECT name
			FROM groups
			WHERE id = ` + ZI(group_id) + `
		),'')`
	group := PG.QStr(query2)
	res := M.SX{
		`id`:      id,
		`user_id`: id,
		`email`:   email,
		`level`: M.SX{
			`group`:         group,
			`company`:       group,
			`group_id`:      group_id,
			`company_id`:    group_id,
			`is_backoffice`: group == `Administrator`,
			`page`: M.SB{
				`guest`:      true,
				`superadmin`: group == `Administrator`,
			},
		},
	}
	return res
}

func API_All_Login(ctx *W.Context) {
	posts := ctx.Posts()
	ident := posts.GetStr(`email`) // or phone
	pass := posts.GetStr(`pass`)
	ajax := sResponse.NewAjax()
	if ajax.HasError() {
		ctx.AppendJson(ajax.SX)
		return
	}
	id := int64(0)
	id = tUsers.FindID_ByIdent_ByPass(ident, pass)
	logged := false
	if id > 0 {
		//tUsers.UpdateLastLogin(id)
		ctx.Session.Login(AccessLevel(ident, id))
		logged = true
		ajax.Set(`logged`, id)
	}
	if !logged {
		ajax.Error(sql.ERR_301_WRONG_USERNAME_OR_PASSWORD)
		T.RandomSleep()
	}
	ctx.AppendJson(ajax.SX)
	//L.Describe(ajax)
}

func API_All_VerifyOAuth(ctx *W.Context) {
	rm := sResponse.Prepare(ctx, `Verify OAuth`, false)
	_ = rm
	params := ctx.QueryParams()

	csrf := ctx.Session.StateCSRF()
	ncsrf := params.GetStr(`state`)
	var err error
	if ncsrf != csrf {
		err = errors.New(sql.ERR_306_CSRF_STATE + ncsrf + ` <> ` + csrf)
	} else {
		var json W.Ajax
		source := ``
		id := int64(0)
		switch ctx.ParamStr(`from`) {
		default:
			g_provider := GetGPlusOAuth(ctx)
			if g_provider == nil {
				err = errors.New(sql.ERR_206_MISSING_OAUTH_PROVIDER)
				break
			}
			json, err = GPlusExchangeInfo(g_provider, params)
			// example: {"email":"","email_verified":true,"family_name":"","gender":"","given_name":"","locale":"en-GB","name":"","picture":"http://","profile":"http://","sub":"number"};
			if err != nil {
				err = errors.New(sql.ERR_201_FAILED_OAUTH_EXCHANGE + err.Error())
				break
			}
			email := json.GetStr(`email`)
			id = tUsers.FindID_ByEmail(email)
			if id == 0 {
				err = errors.New(sql.ERR_305_EMAIL_NOT_REGISTERED + email)
				break
			}
			//tUsers.UpdateLastLogin(id)
			ctx.Session.Login(AccessLevel(email, id))
			source = `Google`
		}
		if err == nil {
			ctx.Render(`login/oauth`, M.SX{
				`data`:          json,
				`redirect_path`: ``,
				`user_id`:       id,
				`webmaster`:     ctx.Engine.WebMasterAnchor,
				`source`:        source,
			})
		}
	}
	if err != nil {
		ctx.Error(403, `OAuth Failed: `+err.Error())
	}
}
