package domain

import (
	"net"
	"time"

	chBuffer "github.com/kokizzu/ch-timed-buffer"
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/M"
	"github.com/rs/zerolog"

	"example2/conf"
	"example2/model/mAuth"
	"example2/model/mAuth/saAuth"
	"example2/model/xMailer"
)

type Domain struct {
	AuthOltp *Tt.Adapter
	AuthOlap *Ch.Adapter

	StorOltp *Tt.Adapter

	Mailer xMailer.Mailer

	IsBgSvc bool // long-running program

	// 3rd party
	Oauth conf.OauthConf

	// oauth related cache
	googleUserInfoEndpointCache string

	// timed buffer
	authLogs *chBuffer.TimedBuffer

	// logger
	Log *zerolog.Logger

	// list of superadmin emails
	Superadmins M.SB
	UploadDir   string
}

// will run in background if background service
func (d *Domain) runSubtask(subTask func()) {
	if d.IsBgSvc {
		go subTask()
	} else {
		subTask()
	}
}

func (d *Domain) InitTimedBuffer() {
	d.authLogs = chBuffer.NewTimedBuffer(d.AuthOlap.DB, 100_000, 1*time.Second, saAuth.Preparators[mAuth.TableActionLogs])
}

func (d *Domain) WaitTimedBufferFinalFlush() {
	<-d.authLogs.WaitFinalFlush
	d.Log.Debug().Msg(`timed buffer flushed`)
}

var defaultIP4 = net.ParseIP(`0.0.0.0`).To4()
var defaultIP6 = net.ParseIP(`0:0:0:0:0:0:0:0`).To16()

func (d *Domain) InsertActionLog(in *RequestCommon, out *ResponseCommon) bool {
	ip := net.ParseIP(in.IpAddress)
	ip4 := ip.To4()
	if ip4 == nil {
		ip4 = defaultIP4
	}
	ip6 := ip.To16()
	if ip6 == nil {
		ip6 = defaultIP6
	}
	row := saAuth.ActionLogs{
		CreatedAt:  in.TimeNow(),
		RequestId:  in.RequestId,
		ActorId:    in.SessionUser.UserId,
		Action:     in.Action,
		StatusCode: int16(out.StatusCode),
		Traces:     out.Traces(),
		Error:      out.Error,
		IpAddr4:    ip4,
		IpAddr6:    ip6,
		UserAgent:  in.UserAgent,
		TenantCode: in.SessionUser.TenantCode,
		Latency:    in.Latency(),
	}
	return d.authLogs.Insert([]any{
		row.CreatedAt,
		row.RequestId,
		row.ActorId,
		row.Action,
		row.StatusCode,
		row.Traces,
		row.Error,
		row.IpAddr4,
		row.IpAddr6,
		row.UserAgent,
		row.Latency,
		row.TenantCode,
		row.RefId,
	})
}

func (d *Domain) CloseTimedBuffer() {
	go d.authLogs.Close()
	d.WaitTimedBufferFinalFlush()
}
