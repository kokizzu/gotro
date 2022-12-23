package domain

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/rqAuth"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/saAuth"

	chBuffer "github.com/kokizzu/ch-timed-buffer"
	"github.com/kpango/fastime"

	"github.com/kokizzu/gotro/L"
)

type Domain struct {
	// add dependencies (tarantool, clickhouse, meilisearch)
	Taran     *Tt.Adapter
	Click     *Ch.Adapter
	chBuffers map[Ch.TableName]*chBuffer.TimedBuffer
	waitGroup *sync.WaitGroup
}

func (d *Domain) InitClickhouseBuffer(preparators map[Ch.TableName]chBuffer.Preparator) {
	for tableName, preparator := range preparators {
		chb := chBuffer.NewTimedBuffer(d.Click.DB, 30000, 1*time.Second, preparator)
		chb.IgnoreInterrupt = true
		d.chBuffers[tableName] = chb
		d.waitGroup.Add(1)
	}
}

func (d *Domain) WaitInterrupt() {
	interrupt := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGKILL)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGQUIT)

	<-interrupt
	L.Print(`caught signal`, interrupt)
}

func (d *Domain) handleTermSignal() {
	d.WaitInterrupt()
	for tableName := range d.chBuffers {
		go func(tableName Ch.TableName) {
			chb := d.chBuffers[tableName]
			chb.Close()
			<-chb.WaitFinalFlush
			L.Print(`done waiting: ` + tableName)
			d.waitGroup.Done()
		}(tableName)
	}
	d.waitGroup.Wait()
	// TODO: how to wait overseer?
	os.Exit(0)
}

type AnalyticsRow interface {
	SqlInsertParam() []any
	TableName() Ch.TableName
}

func (d *Domain) Statistics(row AnalyticsRow) {
	tableName := row.TableName()
	res := d.chBuffers[tableName]
	if res != nil {
		res.Insert(row.SqlInsertParam())
		return
	}
	panic(`did you forgot to register InitClickhouseBuffer preparators for ` + string(tableName))
}

func NewDomain() *Domain {
	d := &Domain{
		Taran: &Tt.Adapter{Connection: conf.ConnectTarantool(), Reconnect: conf.ConnectTarantool},
		Click: &Ch.Adapter{DB: conf.ConnectClickhouse(), Reconnect: conf.ConnectClickhouse},
	}
	d.waitGroup = &sync.WaitGroup{}
	d.chBuffers = map[Ch.TableName]*chBuffer.TimedBuffer{}
	d.InitClickhouseBuffer(saAuth.Preparators)
	// add more preparators if there's new clickhouse tables on model

	go d.handleTermSignal()
	return d
}

func (d *Domain) mustAdmin(token string, userAgent string, out *ResponseCommon) *conf.Session {
	sess := d.mustLogin(token, userAgent, out)
	if sess == nil {
		return nil
	}
	if !conf.Admins[sess.Email] {
		out.SetError(403, `must be admin`)
		return nil
	}
	return sess
}

func (d *Domain) mustLogin(token string, userAgent string, out *ResponseCommon) *conf.Session {
	sess := &conf.Session{}
	if token == `` {
		out.SetError(400, `missing session token`)
		return nil
	}
	if !sess.Decrypt(token, userAgent) {
		out.SetError(400, `invalid session token`) // if got this, possibly wrong userAgent-sessionToken pair
		return nil
	}
	if sess.ExpiredAt <= fastime.UnixNow() {
		out.SetError(400, `token expired`)
		return nil
	}

	session := rqAuth.NewSessions(d.Taran)
	session.SessionToken = token
	if !(token == conf.AdminTestSessionToken) {
		if !session.FindBySessionToken() {
			out.SetError(400, `session missing from database, wrong env?`)
			return nil
		}
		if session.ExpiredAt <= fastime.UnixNow() {
			out.SetError(403, `session expired or logged out`)
			return nil
		}
	}
	return sess
}
