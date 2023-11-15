package domain

import (
	"io"
	"net/http"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

const (
	OauthGoogle = `google`
)

func fetchJsonMap(client *http.Client, url string, res *ResponseCommon) (json M.SX) {
	resp, err := client.Get(url)
	if L.IsError(err, `failed fetch url %s`, url) {
		L.Print(err)
		res.SetError(500, `failed fetch url`)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if L.IsError(err, `failed read body`) {
		L.Print(err)
		res.SetError(500, `failed read body`)
		return
	}
	bodyStr := string(body)
	L.Print(bodyStr)
	json = S.JsonToMap(bodyStr)
	L.Describe(json)
	err2 := json.GetStr(`error`)
	if L.CheckIf(err2 != ``, `fetchJsonMap %s: %#v`, err2, json) {
		res.SetError(500, `error key set from json response`)
		return
	}
	err3 := json.GetStr(`type`)
	if L.CheckIf(err3 == `OAuthException`, `fetchJsonMap %s: %#v`, err3, json) {
		res.SetError(500, `object type from json respons is OAuthException`)
		return
	}
	return
}
