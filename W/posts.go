package W

import (
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/valyala/fasthttp"
)

type Posts struct {
	*fasthttp.Args
	M.SS
}

func (p *Posts) GetJsonMap(key string) M.SX {
	return S.JsonToMap(p.GetStr(key))
}

func (p *Posts) IsSet(key string) bool {
	_, ok := p.SS[key]
	return ok
}

func (p *Posts) GetJsonStrArr(key string) []string {
	return S.JsonToStrArr(p.GetStr(key))
}

func (p *Posts) GetJsonIntArr(key string) []int64 {
	return S.JsonToIntArr(p.GetStr(key))
}

func (p *Posts) FromContext(ctx *Context) {
	p.Args = ctx.RequestCtx.PostArgs()
	p.SS = M.SS{}
	p.Args.VisitAll(func(k, v []byte) {
		p.SS[string(k)] = string(v)
	})
}

// censor the password string, also when length is too long
func (p *Posts) String() string {
	return p.SS.PrettyFunc(` | `, func(key, val string) string {
		if len(val) > 64 {
			return val[:64] + `...`
		}
		return S.IfElse(key == `pass` || key == `password`, S.Repeat(`*`, len(val)), val)
	})
}
