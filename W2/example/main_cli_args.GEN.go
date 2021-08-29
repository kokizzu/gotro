package main

import (
	"context"

	"os"

	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/domain"
)

func cliArgsRunner(args []string) {
	tracerCtx, span := conf.T.Start(context.Background(), args[0])
	defer span.End()

	var (
		vdomain = domain.NewDomain()
	)

	patterns := map[string]map[string]int{
		domain.PlayerChangeEmail_Url:    {},
		domain.PlayerChangePassword_Url: {},
		domain.PlayerConfirmEmail_Url:   {},
		domain.PlayerForgotPassword_Url: {},
		domain.PlayerList_Url:           {},
		domain.PlayerLogin_Url:          {},
		domain.PlayerLogout_Url:         {},
		domain.PlayerProfile_Url:        {},
		domain.PlayerRegister_Url:       {},
		domain.PlayerResetPassword_Url:  {},
		domain.PlayerUpdateProfile_Url:  {},
	}
	switch pattern := cliUrlPattern(args[0], patterns); pattern {

	case domain.PlayerChangeEmail_Url:
		in := domain.PlayerChangeEmail_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerChangeEmail(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerChangePassword_Url:
		in := domain.PlayerChangePassword_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerChangePassword(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerConfirmEmail_Url:
		in := domain.PlayerConfirmEmail_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerConfirmEmail(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerForgotPassword_Url:
		in := domain.PlayerForgotPassword_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerForgotPassword(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerList_Url:
		in := domain.PlayerList_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerList(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerLogin_Url:
		in := domain.PlayerLogin_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerLogin(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerLogout_Url:
		in := domain.PlayerLogout_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerLogout(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerProfile_Url:
		in := domain.PlayerProfile_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerProfile(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerRegister_Url:
		in := domain.PlayerRegister_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerRegister(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerResetPassword_Url:
		in := domain.PlayerResetPassword_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerResetPassword(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	case domain.PlayerUpdateProfile_Url:
		in := domain.PlayerUpdateProfile_In{}
		in.FromCli(os.Stdin, tracerCtx)
		out := vdomain.PlayerUpdateProfile(&in)
		out.ToCli(os.Stdout)
		in.ToCli(os.Stdout, &out)

	}
}
