package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"golang.org/x/sync/errgroup"

	"example2/conf"
	"example2/domain"
	"example2/model"
	"example2/model/xMailer"
	"example2/presentation"
)

func main() {
	log := conf.InitLogger()
	conf.LoadEnv()

	validArgs := `run, web, migrate`
	if len(os.Args) < 2 {
		log.Fatal().Msg(`must have at least one argument with: ` + validArgs)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	var closers []func() error
	var err error

	// mailer
	var mailer xMailer.Mailer
	eg.Go(func() error {
		mailerCfg := conf.EnvMailer()
		fmt.Println(`mailer: ` + mailerCfg.DefaultMailer)
		switch mailerCfg.DefaultMailer {
		case `dockermailserver`:
			dms := xMailer.Dockermailserver{DockermailserverConf: conf.EnvDockermailserver()}
			L.PanicIf(dms.Connect(), `Dockermailserver.Connect`)
			mailer.SendMailFunc = dms.SendEmail
		default: // use mailhog
			mh, err := xMailer.NewMailhog(conf.EnvMailhog())
			L.PanicIf(err, `NewMailhog`)
			mailer.SendMailFunc = mh.SendEmail
		}
		return nil
	})

	// connect to tarantool
	var tConn *Tt.Adapter
	eg.Go(func() error {
		tConf := conf.EnvTarantool()
		tConn, err = tConf.Connect()
		if tConn != nil {
			closers = append(closers, tConn.Close)
			fmt.Println(`tarantool connected: ` + tConf.DebugStr())
		}
		return err
	})

	// connect to clickhouse
	var cConn *Ch.Adapter
	eg.Go(func() error {
		cConf := conf.EnvClickhouse()
		cConn, err = cConf.Connect()
		if cConn != nil {
			closers = append(closers, cConn.Close)
			fmt.Println(`clickhouse connected: ` + cConf.DebugStr())
		}
		return err
	})

	L.PanicIf(eg.Wait(), `eg.Wait`) // if error, make sure no error on: docker compose up
	for _, closer := range closers {
		closer := closer
		defer closer()
	}

	oauth := conf.EnvOauth()
	d := &domain.Domain{
		AuthOltp: tConn,
		AuthOlap: cConn,
		StorOltp: tConn,
		Mailer: xMailer.Mailer{
			SendMailFunc: mailer.SendMailFunc,
		},
		IsBgSvc: false,
		Oauth:   oauth,
		Log:     log,

		UploadDir: conf.UploadDir(),

		Superadmins: conf.EnvSuperAdmins(),
	}
	d.InitTimedBuffer()
	defer d.CloseTimedBuffer()

	mode := S.ToLower(os.Args[1])

	// check table existence
	if mode != `migrate` {
		L.Print(`verifying table schema, if failed, run: go run main.go migrate`)
		model.VerifyTables(tConn, cConn)
	}

	switch mode {
	case `run`:
		cmd := &presentation.CLI{
			Domain: d,
		}
		cmd.Run(os.Args[2:], log)
	case `web`:
		ws := &presentation.WebServer{
			Domain: d,
			Cfg:    conf.EnvWebConf(),
		}
		ws.Start(log)
	case `migrate`:
		model.RunMigration(tConn, cConn)
	default:
		log.Fatal().Msg(`must start with: ` + validArgs)
	}
}
