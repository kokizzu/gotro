package W

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

////////// Ajax

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
		str += S.WebBR
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

func (json Ajax) ErrorIf(err error, msg string) bool {
	if !L.IsError(err, msg) {
		return false
	}
	if msg == `` {
		return true
	}
	if json.SX[`errors`] == nil {
		json.SX[`errors`] = []string{}
	}
	errors := json.SX[`errors`].([]string)
	errors = append(errors, msg)
	json.SX[`errors`] = errors
	json.SX[`is_success`] = false
	return true
}

func (json Ajax) OverwriteInfo(msg string) {
	json.SX[`info`] = msg
}

func (json Ajax) LastError() string {
	if json.SX[`errors`] == nil {
		return ``
	}
	errors := json.SX[`errors`].([]string)
	if len(errors) == 0 {
		return ``
	}
	return errors[len(errors)-1]
}
