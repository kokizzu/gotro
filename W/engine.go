package W

import (
	"github.com/kokizzu/gotro/L"
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

type RequestModel struct {
	Id      int64
	Posts   *Posts
	Ajax    Ajax
	DbActor int64
}

type Ajax struct {
	M.SX
}

func NewAjax() Ajax {
	return Ajax{M.SX{}}
}

func (json Ajax) HasError() bool {
	if json.SX[`errors`] == nil {
		return false
	}
	errors := json.SX[`errors`].([]string)
	return len(errors) > 0
}

func (json Ajax) Info(msg string) {
	str, ok := (json.SX[`info`]).(string)
	if !ok {
		str = ``
	}
	if len(str) > 0 {
		str += "\n<br/>"
	}
	str += msg
	json.SX[`info`] = str
}

func (json Ajax) Error(msg string) string {
	if msg == `` {
		return msg
	}
	if json.SX[`errors`] == nil {
		json.SX[`errors`] = []string{}
	}
	errors := json.SX[`errors`].([]string)
	errors = append(errors, msg)
	json.SX[`errors`] = errors
	json.SX[`is_success`] = false
	L.Describe(`Ajax error`, msg)
	return msg
}
