package domain

import (
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/lexid"
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
	default:
		out.SetError(400, GuestExternalAuthProviderNotSet)
	}

	if in.Redirect {
		out.SetRedirect(out.Link)
	}

	return
}
