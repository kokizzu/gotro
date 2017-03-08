package W

import (
	"github.com/kokizzu/gotro/X"
	"github.com/valyala/fasthttp"
)

type QueryParams struct {
	*fasthttp.Args
}

func (p *QueryParams) GetInt(key string) int64 {
	return X.ToI(p.Peek(key))
}

func (p *QueryParams) GetStr(key string) string {
	return X.ToS(p.Peek(key))
}

func (p *QueryParams) GetFloat(key string) float64 {
	return X.ToF(p.Peek(key))
}
