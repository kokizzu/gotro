package mLogin

import (
	"errors"
	"fmt"
	fb "github.com/huandu/facebook"
	"github.com/kokizzu/goauth2/oauth"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/T"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/model/mUsers"
	"github.com/kokizzu/gotro/X"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"luwes/sql"
	"luwes/sql/nResponse"
	"net/http"
	"net/url"
)

// TODO: kalau pakai Firebase Auth, ntar import https://github.com/robbert229/jwt

// TODO: change to correct one (see console.developer.google.com

var OAUTH_URLS = []string{
	`http://test.localhost`,
	`https://test.localhost`,
}
var GPLUS_OAUTH_PROVIDERS map[string]*oauth2.Config
var FB_OAUTH_PROVIDERS map[string]*fbConfig
var YT_OAUTH_PROVIDERS map[string]*oauth.Config
var USERINFO_ENDPOINT string

const GPLUS_CLIENTID = `*********************************************.apps.googleusercontent.com`
const GPLUS_CLIENTSECRET = `************************`

const FB_API_VER = `v2.8`
const FB_APPID = `****************`
const FB_APPSECRET = `********************************`
const FB_TOKEN_ENDPOINT = `https://graph.facebook.com/` + FB_API_VER + `/oauth/access_token`

const FB_AK_VER = `v1.1`
const FB_AP_SECRET = `********************************` // client token: not used
const FB_AK_TOKEN_ENDPOINT = `https://graph.accountkit.com/` + FB_AK_VER + `/access_token`
const FB_AK_USER_ENDPOINT = `https://graph.accountkit.com/` + FB_AK_VER + `/me`

const RESET_MINUTE = 20

// credential for OpenID
// https://console.developers.google.com/apis/credentials

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

// 2017-01-10 Prayogo
func RetrieveFacebookUserInfo(access_token string) (W.Ajax, error) {
	params := url.Values{}
	params.Add(`access_token`, access_token)
	params.Add(`fields`, `id,name,email,birthday,gender,hometown,languages,location,religion,picture.type(large)`)
	return fetchJson(`https://graph.facebook.com/v2.8/me?` + params.Encode())
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

// 2016-11-23 Prayogo
func (f *fbConfig) ExchangeInfo(code string) (W.Ajax, error) {
	fb.Debug = fb.DEBUG_ALL
	params := url.Values{}
	params.Add(`client_id`, f.ClientID)
	params.Add(`redirect_uri`, f.RedirectURL)
	params.Add(`client_secret`, f.ClientSecret)
	params.Add(`code`, code)
	m, err := fetchJson(FB_TOKEN_ENDPOINT + `?` + params.Encode())
	if err != nil {
		return m, err
	}
	tok := m.GetStr(`access_token`)
	return RetrieveFacebookUserInfo(tok)
}

func init() {
	GPLUS_OAUTH_PROVIDERS = map[string]*oauth2.Config{} // yahoo tidak bisa multiple domain (harus dibuat 1-1), tidak support IP
	FB_OAUTH_PROVIDERS = map[string]*fbConfig{}
	YT_OAUTH_PROVIDERS = map[string]*oauth.Config{}
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
		FB_OAUTH_PROVIDERS[url] = &fbConfig{
			ClientID:     FB_APPID,
			ClientSecret: FB_APPSECRET,
			RedirectURL:  url + `/login/verify/fb`,
			Scopes: []string{
				`user_about_me`,
				`user_birthday`,
				`user_education_history`,
				`user_hometown`,
				`user_location`,
				`user_religion_politics`,
				`user_relationships`,
				`user_website`,
				`email`,
				//www.facebook.com/dialog/oauth?display=page&client_id=365761640428516&redirect_uri=https%3A%2F%2Fluweswatersensor.com%2Flogin%2Fverify%2Ffb&scope=user_about_me,user_birthday,user_education_history,user_hometown,user_location,user_religion_politics,user_relationships,user_website,email
			},
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

// get FB OAuth provider and domain csrf
// 2016-07-26 Prayogo
func GetFBOAuth(ctx *W.Context) *fbConfig {
	return FB_OAUTH_PROVIDERS[ctx.Host()]
}

// get OAuth provider and domain csrf
func GetYTOAuth(ctx *W.Context) *oauth.Config {
	return YT_OAUTH_PROVIDERS[ctx.Host()]
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

// 2016-11-04 Prayogo, when using facebook SDK (2.8+)
func CheckFacebook(posts *W.QueryParams) (map[string]interface{}, error) {
	fb.Version = FB_API_VER
	return fb.Get(`/`+posts.GetStr(`userID`), fb.Params{
		`fields`:       `name,email,birthday,gender,hometown,languages,location,religion,picture`,
		`access_token`: posts.GetStr(`accessToken`),
	})
}

// 2017-01-12 Prayogo
func CheckFacebookAccountKit(phone, code string, ajax W.Ajax) (is_phone bool, json M.SX) {
	if len(phone) == 0 || phone[0] != '+' {
		return
	}
	is_phone = true
	// exchange code from user to access_token
	app_access_token := A.StrJoin([]string{`AA`, FB_APPID, FB_AP_SECRET}, `|`)
	params := url.Values{}
	params.Add(`grant_type`, `authorization_code`)
	params.Add(`code`, code)
	params.Add(`access_token`, app_access_token)
	token_exchange_url := FB_AK_TOKEN_ENDPOINT + `?` + params.Encode()
	var err error
	json2, err := fetchJson(token_exchange_url)
	if err != nil {
		ajax.Error(sql.ERR_207_FB_AK_TOKEN_EXCHANGE_ERROR + err.Error())
		L.Describe(ajax)
		return
	}
	L.Describe(`json2`, json2)
	user_access_token := json2.GetStr(`access_token`)
	expires_at := json.GetStr(`expires_at`)
	ajax.Set(`expires_at`, expires_at)
	user_id := json2.GetStr(`id`)
	ajax.Set(`user_id`, user_id)
	// phone from FB AK
	me_endpoint_url := FB_AK_USER_ENDPOINT + `?access_token=` + user_access_token
	json3, err := fetchJson(me_endpoint_url)
	if err != nil {
		ajax.Error(sql.ERR_208_FB_AK_USER_INFO_ERROR + err.Error())
		return
	}
	L.Describe(`json3`, json3)

	// flattend and reformat phone and email
	json = M.SX{}
	m := json3.GetMSX(`phone`)
	if len(m) > 0 {
		json[`phone`] = `+` + X.ToS(m[`country_prefix`]) + ` ` + X.ToS(m[`national_number`])
	}
	m = json3.GetMSX(`email`)
	if len(m) > 0 {
		json[`email`] = X.ToS(m[`address`])
	}
	// other info from graph API, required for subscriber (not possible, FB AK <> FB GraphAPI)
	//posts := W.NewPosts()
	//posts.Add(`userID`, user_id)
	//posts.Add(`accessToken`, user_access_token)
	//json4, err := CheckFacebook(posts)
	//L.Print(`json4`, json4)
	//for k, v := range json4 {
	//	json[k] = v
	//}
	L.Describe(`json`, json)
	return
}

func API_All_Logout(ctx *W.Context) {
	ajax := nResponse.NewAjax()
	user_id := ctx.Session.GetInt(`user_id`)
	if user_id > 0 && !ctx.IsWebMaster() {
		// TODO: update last login
	}
	ctx.Session.Logout()
	ctx.AppendJson(ajax.SX)
}

func API_All_Login(ctx *W.Context) {
	posts := ctx.Posts()
	ident := posts.GetStr(`email`) // or phone
	pass := posts.GetStr(`pass`)
	ajax := nResponse.NewAjax()
	is_phone, json := CheckFacebookAccountKit(ident, pass, ajax)
	if ajax.HasError() {
		ctx.AppendJson(ajax.SX)
		return
	}
	id := int64(0)
	if is_phone {
		ident = X.ToS(json[`phone`])
		id = mUsers.FindID_ByPhone(ident)
	} else {
		id = mUsers.FindID_ByIdent_ByPass(ident, pass)
	}
	logged := false
	if id > 0 {
		mUsers.UpdateLastLogin(id)
		ctx.Session.Login(M.SX{
			`email`:   ident,
			`user_id`: id,
			`level`:   M.SX{},
		})
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
	rm := nResponse.Prepare(ctx, `Verify OAuth`, false)
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
		case `fb`:
			f_provider := GetFBOAuth(ctx)
			if f_provider == nil {
				err = errors.New(sql.ERR_206_MISSING_OAUTH_PROVIDER)
				break
			}
			json, err = f_provider.ExchangeInfo(params.GetStr(`code`))
			// example: { "name": "Kiswono Prayogo", "email": "kiswono@gmail.com", "gender":  "male", "picture": { "data": { "is_silhouette": false, "url": "https://",}, }, "id": "561039484102125" }
			if err != nil {
				err = errors.New(sql.ERR_201_FAILED_OAUTH_EXCHANGE + err.Error())
				break
			}
			L.Print(params)
			email := json.GetStr(`email`)
			id = mUsers.FindID_ByEmail(email)
			if id == 0 {
				err = errors.New(sql.ERR_305_EMAIL_NOT_REGISTERED + email)
				break
			}
			mUsers.UpdateLastLogin(id)
			ctx.Session.Login(M.SX{
				`email`:   email,
				`user_id`: id,
				`level`:   M.SX{},
			})
			source = `Facebook`
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
			id = mUsers.FindID_ByEmail(email)
			if id == 0 {
				err = errors.New(sql.ERR_305_EMAIL_NOT_REGISTERED + email)
				break
			}
			mUsers.UpdateLastLogin(id)
			ctx.Session.Login(M.SX{
				`email`:   email,
				`user_id`: id,
				`level`:   M.SX{},
			})
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

func API_All_LoginReset(ctx *W.Context) {
	posts := ctx.Posts()
	email := posts.GetStr(`email`)
	pass := posts.GetStr(`pass`)
	ajax := nResponse.NewAjax()
	key := ctx.ParamStr(`key`)
	id := W.Sessions.GetStr(`reset-link:` + key)
	if id == `` || id == `0` {
		ajax.Error(sql.ERR_004_INVALID_RESET_LINK)
		ctx.AppendJson(ajax.SX)
		T.RandomSleep()
		return
	}
	e_id := int64(0)
	e_id = mUsers.FindID_ByEmail(email)
	if e_id == 0 {
		ajax.Error(sql.ERR_005_INVALID_RESET_EMAIL)
		ctx.AppendJson(ajax.SX)
		T.RandomSleep()
		return
	}
	if I.ToS(e_id) != id {
		ajax.Error(sql.ERR_006_INVALID_RESET_ASSOC)
		ctx.AppendJson(ajax.SX)
		T.RandomSleep()
		return
	}
	mUsers.PG_W.DoTransaction(func(tx *Pg.Tx) string {
		var dm *Pg.Row
		rm := W.NewRequestModel_ById_ByDbActor_ByAjax(id, id, ajax)
		dm = Pg.NewRow(tx, mUsers.TABLE, rm)
		dm.Set_UserPassword(pass)
		if dm.UpdateRow() > 0 {
			ajax.OverwriteInfo(`password reset successful, you may login now`)
		}
		W.Sessions.Del(`reset-link:` + key)
		W.Sessions.Del(`reset-req:` + id)
		return ajax.LastError()
	})
	ctx.AppendJson(ajax.SX)
}

func API_All_LoginForgot(ctx *W.Context) {
	posts := ctx.Posts()
	ident := posts.GetStr(`ident`)
	email := posts.GetStr(`email`)
	ajax := nResponse.NewAjax()
	id := int64(0)
	id = mUsers.FindID_ByCompactName_ByEmail(ident, email)
	if id == 0 {
		ajax.Error(sql.ERR_302_NAME_EMAIL_COMBINATION_NOT_FOUND)
		ctx.AppendJson(ajax.SX)
		T.RandomSleep()
		return
	}
	key := `reset-req:` + I.ToS(id)
	ttl := W.Sessions.Expiry(key)
	ttl_str := I.ToS(ttl)
	if ttl > 0 {
		ajax.Error(sql.ERR_303_TOO_SOON_RESET + ttl_str + `s`)
		ctx.AppendJson(ajax.SX)
		T.RandomSleep()
		return
	}
	full_name, emails := ``, []string{}
	full_name, emails = mUsers.Name_Emails_ByID(id)
	r := S.RandomCB63(8)
	W.Sessions.FadeStr(key, r, 60*RESET_MINUTE)
	url := ctx.Host() + `/login/reset/` + r
	validity := `The reset password link only valid for ` + I.ToS(RESET_MINUTE) + ` minutes since the request received.`
	err := ctx.Engine.SendMailSync(``, emails, `[CHANGEME] Reset Password Link`, `Dear `+full_name+`,<br/>
<br/>
We've received a request to send a reset password link for your account at `+T.HumanStr()+`<br/>
from `+ctx.Session.UserAgent+`<br/>
<br/>
To initiate the password reset process for your CHANGEME Account, click the link below:<br/>
<br/>
<a href='`+url+`' />`+url+`</a><br/>
<br/>
If clicking the link above doesn't work, please copy and paste the URL in a new browser window instead.<br/>
`+validity+`<br/>
<br/>
If you've received this mail in error, it's likely that another user entered your email address by mistake while trying to reset a password. If you didn't initiate the request, you don't need to take any further action and can safely disregard this email.<br/>`)
	if err == `` {
		W.Sessions.FadeInt(`reset-link:`+r, id, 60*RESET_MINUTE)
		ajax.Info(`A reset password link has been sent to all of your registered e-mails, please check your e-mails to reset your password. ` + validity)
		mUsers.UpdateLastForgotPassword(id)
	} else {
		ajax.Error(sql.ERR_304_FAILED_SEND_RESET_EMAIL + sql.SUPPORT_EMAIL) // TODO: replace hardcoded when finance done
	}
	ctx.AppendJson(ajax.SX)
}
