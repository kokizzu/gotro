package zCrud

import (
	"sync"

	"github.com/kokizzu/gotro/S"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file form.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type form.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type form.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type form.go
//go:generate farify doublequote --file form.go

type DataType string

type InputType string

type Validation string

const (
	DataTypeString DataType = `string`
	DataTypeInt    DataType = `int`
	DataTypeFloat  DataType = `float`
	DataTypeIntArr DataType = `intArr`

	InputTypeText     InputType = `text`
	InputTypeTextArea InputType = `textarea`
	InputTypeEmail    InputType = `email`
	InputTypePassword InputType = `password`
	InputTypeCombobox InputType = `combobox`
	InputTypeCheckbox InputType = `checkbox`
	InputTypeNumber   InputType = `number`
	InputTypeHidden   InputType = `hidden`
	InputTypeDateTime InputType = `datetime`

	ValidationRequired Validation = `required`
	ValidationMinLen   Validation = `minLen`
	ValidationMaxLen   Validation = `maxLen`
	ValidationRegex    Validation = `regex`
)

type Field struct {
	Name        string    `json:"name" form:"name" query:"name" long:"name" msg:"name"`
	Label       string    `json:"label" form:"label" query:"label" long:"label" msg:"label"`
	Description string    `json:"description" form:"description" query:"description" long:"description" msg:"description"`
	DataType    DataType  `json:"type" form:"type" query:"type" long:"type" msg:"type"`
	InputType   InputType `json:"inputType" form:"inputType" query:"inputType" long:"inputType" msg:"inputType"`
	ReadOnly    bool      `json:"readOnly" form:"readOnly" query:"readOnly" long:"readOnly" msg:"readOnly"`

	Validations map[Validation]any `json:"validations" form:"validations" query:"validations" long:"validations" msg:"validations"`

	// fixed value for combobox/select, must be under 5 rows
	Ref []string `json:"ref" form:"ref" query:"ref" long:"ref" msg:"ref"`
	// endpoint to find the combobox reference, if combobox/select source for autocomplete is too large
	RefEndpoint string `json:"refEndpoint" form:"refEndpoint" query:"refEndpoint" long:"refEndpoint" msg:"refEndpoint"`
}

type Meta struct {
	Fields []Field `json:"fields" form:"fields" query:"fields" long:"fields" msg:"fields"`

	mutex        sync.Mutex
	cachedSelect string
}

func (m *Meta) ToSelect() string {
	if m.cachedSelect > `` {
		return m.cachedSelect
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, f := range m.Fields {
		// our internal, so safe from sql injection
		m.cachedSelect += `, ` + S.QQ(f.Name)
	}
	// assume more than 1
	m.cachedSelect = m.cachedSelect[1:]
	return m.cachedSelect
}
