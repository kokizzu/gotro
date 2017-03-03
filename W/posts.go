package W

import (
	"github.com/kokizzu/gotro/M"
)

type Posts struct {
}

func (p *Posts) GetStr(key string) (res string) {
	return res // TODO: continue this
}

func (p *Posts) GetInt(key string) (res int64) {
	return res // TODO: continue this
}

func (p *Posts) GetFloat(key string) (res float64) {
	return res // TODO: continue this
}

func (p *Posts) GetJsonMap(key string) (res M.SX) {
	res = M.SX{}
	return res // TODO: continue this
}

func (p *Posts) GetJsonStrArr(key string) (res []string) {
	res = []string{}
	return res // TODO: continue this
}

func (p *Posts) GetJsonIntArr(key string) (res []int64) {
	res = []int64{}
	return res // TODO: continue this
}

func (p *Posts) FromContext(ctx *Context) {
	// TODO: continue this
}
