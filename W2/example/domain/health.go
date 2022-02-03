package domain

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kpango/fastime"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"math"
)

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
