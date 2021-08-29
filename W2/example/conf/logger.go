package conf

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLogger() *zerolog.Logger {
	const DirTrim = PROJECT_NAME
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: `2006-01-02 15:04:05`, PartsOrder: []string{
		zerolog.TimestampFieldName,
		zerolog.LevelFieldName,
		zerolog.CallerFieldName, // clickable on IntelliJ if first or starts with ./
		zerolog.MessageFieldName,
	}, FormatCaller: func(i interface{}) string {
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
	return &zlog.Logger
}
