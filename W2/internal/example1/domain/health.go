package domain

import (
	"math"

	"github.com/kokizzu/lexid"
	"github.com/kokizzu/rand"
	"github.com/kpango/fastime"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	"example1/model/mAuth/rqAuth"
	"example1/model/mAuth/wcAuth"
	"github.com/kokizzu/gotro/L"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file health.go
//go:generate replacer -afterprefix 'Id" form' 'Id,string" form' type health.go
//go:generate replacer -afterprefix 'json:"id"' 'json:"id,string"' type health.go
//go:generate replacer -afterprefix 'By" form' 'By,string" form' type health.go
//go:generate farify doublequote --file health.go

type (
	Health_In struct {
		RequestCommon
	}
	Health_Out struct {
		ResponseCommon
		CpuPercent  float32 `json:"cpuPercent" form:"cpuPercent" query:"cpuPercent" long:"cpuPercent" msg:"cpuPercent"`
		RamPercent  float32 `json:"ramPercent" form:"ramPercent" query:"ramPercent" long:"ramPercent" msg:"ramPercent"`
		DiskPercent float32 `json:"diskPercent" form:"diskPercent" query:"diskPercent" long:"diskPercent" msg:"diskPercent"`
	}
)

const Health_Url = `/Health`

var healthLastTime int64
var healthLastResult [3]float32

func round(f float64) float32 {
	return float32(math.Round(f*100)) / 100
}

func (d *Domain) Health(in *Health_In) (out Health_Out) {
	now := fastime.Now().Unix()
	if now == healthLastTime { // throttle request 1 same result per second
		out.CpuPercent = healthLastResult[0]
		out.RamPercent = healthLastResult[1]
		out.DiskPercent = healthLastResult[2]
		return
	}
	healthLastTime = now
	defer func() {
		healthLastResult[0] = out.CpuPercent
		healthLastResult[1] = out.RamPercent
		healthLastResult[2] = out.DiskPercent
	}()

	proc, err := cpu.Percent(0, false)
	if L.IsError(err, `cpu.Percent`) {
		out.SetError(500, `failed fetch cpu usage`)
	} else {
		out.CpuPercent = round(proc[0])
	}

	ram, err := mem.VirtualMemory()
	if L.IsError(err, `mem.VirtualMemory`) {
		out.SetError(500, `failed fetch ram usage`)
	} else {
		out.RamPercent = round(ram.UsedPercent)
	}

	root, err := disk.Usage(`/`)
	if L.IsError(err, `disk.Usage`) {
		out.SetError(500, `failed fetch disk usage`)
	} else {
		out.DiskPercent = round(root.UsedPercent)
	}
	return
}

type (
	LoadTestWrite_In struct {
		RequestCommon
	}
	LoadTestWrite_Out struct {
		ResponseCommon
		Ok bool   `json:"ok" form:"ok" query:"ok" long:"ok" msg:"ok"`
		Id uint64 `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	}
)

const LoadTestWrite_Url = `/LoadTestWrite`

var userIds = make([]uint64, 0, 1024*1024*32) // 32 million ids

func (d *Domain) LoadTestWrite(_ *LoadTestWrite_In) (out LoadTestWrite_Out) {
	user := wcAuth.NewUsersMutator(d.Taran)
	user.Email = lexid.ID() + `@localhost`
	user.Password = `123`
	out.Ok = user.DoInsert()
	if user.Id > 0 {
		out.Id = user.Id
		userIds = append(userIds, user.Id)
	}
	return
}

type (
	LoadTestRead_In struct {
		RequestCommon
		Id uint64 `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	}
	LoadTestRead_Out struct {
		ResponseCommon
		User rqAuth.Users `json:"user" form:"user" query:"user" long:"user" msg:"user"`
	}
)

const LoadTestRead_Url = `/LoadTestRead`

func (d *Domain) LoadTestRead(_ *LoadTestRead_In) (out LoadTestRead_Out) {
	user := rqAuth.NewUsers(d.Taran)
	user.Id = userIds[rand.Intn(len(userIds))]
	user.FindById()
	out.User = *user
	return
}

type (
	LoadHello_In struct {
		RequestCommon
	}
	LoadHello_Out struct {
		ResponseCommon
		Hello string `json:"hello" form:"hello" query:"hello" long:"hello" msg:"hello"`
	}
)

const LoadHello_Url = `/LoadHello`

func (d *Domain) LoadHello(_ *LoadHello_In) (out LoadHello_Out) {
	out.Hello = `world`
	return
}
