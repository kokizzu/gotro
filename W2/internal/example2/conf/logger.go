package conf

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kokizzu/lexid"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *zerolog.Logger {
	const DirTrim = PROJECT_NAME
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if VERSION == `` {
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: `2006-01-02 15:04:05`, PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName, // clickable on IntelliJ if first or starts with ./
			zerolog.MessageFieldName,
		}, FormatCaller: func(i any) string {
			var c string
			if cc, ok := i.(string); ok {
				c = cc
			}
			if len(c) > 0 {
				// vscode use full path, "./" only work on correct pwd
				pos := strings.LastIndex(c, DirTrim)
				if pos > 0 {
					pos += len(DirTrim)
					if len(c) > pos {
						c = c[pos:]
						c = `.` + c
					}
				}
			}
			return c
		}}
		zlog.Logger = zlog.Output(ConsoleWriter).With().Timestamp().Stack().Caller().Logger()
	} else {
		zlog.Logger = zerolog.New(&lumberjack.Logger{
			Filename: `all.log`,
			MaxSize:  128, // MB
		}).With().Timestamp().Stack().Logger()
	}
	return &zlog.Logger
}

type logFields struct {
	ID         string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	StatusCode int
	Latency    float64
	Error      error
	Stack      []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", lf.ID).
		Str("remote_ip", lf.RemoteIP).
		Str("host", lf.Host).
		Str("path", lf.Path).
		Int("status_code", lf.StatusCode).
		Float64("latency", lf.Latency).
		Str("tag", "request")

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}

// Logger Middleware requestid + logger + recover for request traceability
func Logger(log zerolog.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		rid := c.Get(fiber.HeaderXRequestID)
		if rid == `` {
			rid = `r-` + lexid.ID()
			c.Set(fiber.HeaderXRequestID, rid)
		}

		fields := &logFields{
			ID:       rid,
			RemoteIP: c.IP(),
			Method:   c.Method(),
			Host:     c.Hostname(),
			Path:     c.Path(),
		}

		defer func() {
			rvr := recover()

			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				fields.Error = err
				fields.Stack = debug.Stack()

				c.Status(http.StatusInternalServerError)
				_ = c.JSON(map[string]any{
					"status": http.StatusText(http.StatusInternalServerError),
				})
			}

			fields.StatusCode = c.Response().StatusCode()
			fields.Latency = time.Since(start).Seconds()

			switch {
			case rvr != nil:
				log.Error().EmbedObject(fields).Msg("panic recover")
			case fields.StatusCode >= 500:
				log.Error().EmbedObject(fields).Msg("server error")
			case fields.StatusCode >= 400:
				log.Error().EmbedObject(fields).Msg("client error")
			case fields.StatusCode >= 300:
				log.Warn().EmbedObject(fields).Msg("redirect")
			case fields.StatusCode >= 200:
				log.Info().EmbedObject(fields).Msg("success")
			case fields.StatusCode >= 100:
				log.Info().EmbedObject(fields).Msg("informative")
			default:
				log.Warn().EmbedObject(fields).Msg("unknown status")
			}
		}()

		return c.Next()
	}
}
