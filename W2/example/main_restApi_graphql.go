package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/domain"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type GraphqlRequest struct {
	domain.RequestCommon
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
	Mutation      string `json:"mutation"`
}

type GraphqlResponse struct {
	domain.ResponseCommon
	*graphql.Result
}

func webApiInitGraphql(app *fiber.App) {
	const url = `/graphql`
	app.All(url, func(ctx *fiber.Ctx) error {
		tracerCtx, span := conf.T.Start(ctx.Context(), url, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		isGet := string(ctx.Request().Header.Method()) == http.MethodGet
		in := GraphqlRequest{}
		if err := webApiParseInput(ctx, &in.RequestCommon, &in, url); err != nil {
			return err
		}
		in.FromFiberCtx(ctx, tracerCtx)
		out := GraphqlResponse{}
		if isGet {
			ctx.WriteString(graphqlTemplate)
		}
		handleGraphql(&in, &out)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		err := in.ToFiberCtx(ctx, out)
		if isGet {
			ctx.Set(`content-type`, `text/html; charset=utf-8`)
		}
		return err
	})
}

const graphqlTemplate = `<!DOCTYPE html>
<html>
	<head>
		<!-- Copyright (c) 2021 GraphQL Contributors -->
		<style>
			body {
				height: 100%;
				margin: 0;
				width: 100%;
				overflow: hidden;
			}
			
			#graphiql {
				height: 100vh;
			}
		</style>
		<script crossorigin src="https://unpkg.com/react@16/umd/react.development.js"></script>
		<script crossorigin src="https://unpkg.com/react-dom@16/umd/react-dom.development.js"></script>
		<link rel="stylesheet" href="https://unpkg.com/graphiql/graphiql.min.css" />
	</head>
	<body>
		<div id="graphiql">Loading...</div>
		<script src="https://unpkg.com/graphiql/graphiql.min.js"	type="application/javascript"></script>
		<script>
			var fetcher = GraphiQL.createFetcher({
				url: '/graphql',
			});
			ReactDOM.render(
				React.createElement(GraphiQL, {
					fetcher: fetcher,
					headerEditorEnabled: true,
					defaultVariableEditorOpen: true,
				}),
				document.getElementById('graphiql'),
			);
		</script>
	</body>
</html>`
