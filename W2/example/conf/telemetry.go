package conf

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var T trace.Tracer

func init() {
	T = otel.Tracer(PROJECT_NAME)
}
