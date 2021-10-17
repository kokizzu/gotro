package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/domain"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type Inputs struct {
	OperationName string                 `json:"operationName" form:"operationName" query:"operationName"`
	Query         string                 `json:"query" form:"query" query:"query"`
	Mutation      string                 `json:"mutation" form:"mutation" query:"mutation"`
	Variables     map[string]interface{} `json:"variables" form:"variables" query:"variables"`
}

const RequestCommonKey = `RC`

type GraphqlRequest struct {
	domain.RequestCommon
	Inputs
}

type GraphqlResponse struct {
	domain.ResponseCommon
	Inputs
	*graphql.Result
}

func webApiInitGraphql(app *fiber.App, d *domain.Domain) {
	const url = `/graphql`

	graphqlSchema := initGraphqlSchemaResolver(d)

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
			ctx.Set(`content-type`, `text/html; charset=utf-8`)
			return nil
		}
		params := graphql.Params{
			Context:        context.WithValue(ctx.Context(), RequestCommonKey, &in.RequestCommon),
			Schema:         graphqlSchema,
			RequestString:  in.Query,
			OperationName:  in.OperationName,
			VariableValues: in.Variables,
		}
		out.Result = graphql.Do(params)
		L.Describe(out)
		out.ToFiberCtx(ctx, &in.RequestCommon, &in)
		out.Inputs = in.Inputs
		L.Describe(out)
		err := ctx.JSON(out)
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
