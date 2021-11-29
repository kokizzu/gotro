package domain

//go:generate gomodifytags -file health.go -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported --skip-unexported -w -file health.go
//go:generate replacer 'Id" form' 'Id,string" form' type health.go
//go:generate replacer 'json:"id"' 'json:id,string" form' type health.go
//go:generate replacer 'By" form' 'By,string" form' type health.go

type (
	Health_In struct {
		RequestCommon
	}
	Health_Out struct {
		ResponseCommon
	}
)

const Health_Url = `/Health`

func (d *Domain) Health(in *Health_In) (out Health_Out) {
	return
}
