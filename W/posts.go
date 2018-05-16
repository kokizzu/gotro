package W

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/valyala/fasthttp"
)

type Posts struct {
	*fasthttp.Args
	M.SS
}

func (p *Posts) GetBool(key string) bool {
	val := p.SS[key]
	return !(val == `` || val == `0` || val == `f` || val == `false`)
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

func (p *Posts) GetJsonObjArr(key string) []map[string]interface{} {
	return S.JsonToObjArr(p.GetStr(key))
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
	mf, err := ctx.RequestCtx.MultipartForm()
	if err == nil {
		for k, v := range mf.Value {
			p.SS[k] = v[0]
		}
	} else {
		L.Print(`Error Parsing Post Data: ` + err.Error())
	}
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
