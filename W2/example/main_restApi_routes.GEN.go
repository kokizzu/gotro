package main

import (
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/domain"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
	//"github.com/kokizzu/gotro/S"
)

func webApiInitRoutes(app *fiber.App) *domain.Domain {
	var (
		vdomain = domain.NewDomain()
	)

	app.All(conf.API_PREFIX+domain.PlayerChangeEmail_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerChangeEmail_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerChangeEmail_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerChangeEmail(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerChangePassword_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerChangePassword_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerChangePassword_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerChangePassword(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerConfirmEmail_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerConfirmEmail_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerConfirmEmail_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerConfirmEmail(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerForgotPassword_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerForgotPassword_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerForgotPassword_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerForgotPassword(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerList_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerList_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerList_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerList(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerLogin_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerLogin_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerLogin_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerLogin(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerLogout_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerLogout_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerLogout_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerLogout(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerProfile_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerProfile_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerProfile_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerProfile(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerRegister_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerRegister_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerRegister_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerRegister(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerResetPassword_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerResetPassword_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerResetPassword_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerResetPassword(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	app.All(conf.API_PREFIX+domain.PlayerUpdateProfile_Url, func(ctx *fiber.Ctx) error {
		url := domain.PlayerUpdateProfile_Url
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		in := domain.PlayerUpdateProfile_In{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := vdomain.PlayerUpdateProfile(&in)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		return in.ToFiberCtx(ctx, out)
	})

	return vdomain
}
