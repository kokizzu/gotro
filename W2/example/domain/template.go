package domain

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file template.go
//go:generate replacer 'Id" form' 'Id,string" form' type template.go
//go:generate replacer 'json:"id"' 'json:id,string" form' type template.go
//go:generate replacer 'By" form' 'By,string" form' type template.go
// go:generate msgp -tests=false -file template.go -o template__MSG.GEN.go

// copy this template if need new API

type (
	XXX_In struct {
		RequestCommon
	}
	XXX_Out struct {
		ResponseCommon
	}
)

const XXX_Url = `/XXX`

func (d *Domain) XXX(in *XXX_In) (out XXX_Out) {
	// TODO: continue this
	return
}
