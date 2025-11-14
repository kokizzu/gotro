package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/lexid"
	"golang.org/x/oauth2"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestExternalAuth.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestExternalAuth.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestExternalAuth.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestExternalAuth.go
//go:generate farify doublequote --file GuestExternalAuth.go

type (
	GuestExternalAuthIn struct {
		RequestCommon
		Provider string `json:"provider" form:"provider" query:"provider" long:"provider" msg:"provider"`
		Redirect bool   `json:"redirect" form:"redirect" query:"redirect" long:"redirect" msg:"redirect"`
	}
	GuestExternalAuthOut struct {
		ResponseCommon
		Link string `json:"link" form:"link" query:"link" long:"link" msg:"link"`

		// these for manual client-side oauth link generation
		ClientID    string   `json:"clientId,string" form:"clientId" query:"clientId" long:"clientId" msg:"clientId"`
		RedirectUrl string   `json:"redirectUrl" form:"redirectUrl" query:"redirectUrl" long:"redirectUrl" msg:"redirectUrl"`
		Scopes      []string `json:"scopes" form:"scopes" query:"scopes" long:"scopes" msg:"scopes"`
		CsrfState   string   `json:"csrfState" form:"csrfState" query:"csrfState" long:"csrfState" msg:"csrfState"`
	}
)

const (
	GuestExternalAuthAction = `guest/externalAuth`

	GuestExternalAuthProviderNotSet = `oauth provider not set`
	GuestExternalAuthInvalidUrl     = `oauth provider invalid url`
)

func (d *Domain) GuestExternalAuth(in *GuestExternalAuthIn) (out GuestExternalAuthOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	csrfState := in.Provider + `|`
	if in.SessionToken == `` {
		in.SessionToken = `TEMP__` + lexid.ID()
		out.SessionToken = in.SessionToken
	}
	csrfState += S.Left(in.SessionToken, 20)

	switch in.Provider {
	case OauthGoogle:
		provider := d.Oauth.Google[in.Host]
		if provider == nil {
			out.SetError(400, GuestExternalAuthInvalidUrl)
			return
		}
		out.Link = provider.AuthCodeURL(csrfState)
		out.ClientID = provider.ClientID
		out.RedirectUrl = provider.RedirectURL
		out.Scopes = provider.Scopes
		out.CsrfState = csrfState

	case OauthTwitter:
		provider := d.Oauth.Twitter[in.Host]
		if provider == nil {
			out.SetError(400, GuestExternalAuthInvalidUrl)
			return
		}

		codeVerifier := csrfState

		hash := sha256.Sum256([]byte(codeVerifier))
		codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

		scopes := strings.Join(provider.Scopes, " ")

		out.Link = provider.AuthCodeURL(csrfState,
			oauth2.SetAuthURLParam("redirect_uri", provider.RedirectURL),
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("scope", scopes),
		)

		out.ClientID = provider.ClientID
		out.RedirectUrl = provider.RedirectURL
		out.Scopes = provider.Scopes
		out.CsrfState = csrfState
	default:
		out.SetError(400, GuestExternalAuthProviderNotSet)
	}

	if in.Redirect {
		out.SetRedirect(out.Link)
	}

	return
}
