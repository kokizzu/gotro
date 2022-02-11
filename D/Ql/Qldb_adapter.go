package Ql

import (
	"context"
	"errors"
	"github.com/amzn/ion-go/ion"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/v2/qldbdriver"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
)

// https://docs.aws.amazon.com/qldb/latest/developerguide/console_QLDB.html#partiql-editor-ref-tips
// https://docs.aws.amazon.com/qldb/latest/developerguide/driver-quickstart-golang.html

type Adapter struct {
	*qldbdriver.QLDBDriver
	Reconnect func() *qldbdriver.QLDBDriver
}

func (a *Adapter) Shutdown() {
	a.QLDBDriver.Shutdown(context.Background())
}

func Connect1(keyId, secret, region, ledger string) *qldbdriver.QLDBDriver {
	conf := aws.NewConfig()
	conf.WithRegion(region)
	if keyId == `` || secret == `` {
		// AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
		// AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
		conf.WithCredentials(credentials.NewEnvCredentials())
	} else {
		conf.WithCredentials(credentials.NewStaticCredentials(keyId, secret, ``))
	}

	awsSession, err := session.NewSession(conf)
	L.PanicIf(err, `session.NewSession `+region)
	qldbSession := qldbsession.New(awsSession)

	driver, err := qldbdriver.New(
		ledger,
		qldbSession,
		func(options *qldbdriver.DriverOptions) {
			options.LoggerVerbosity = qldbdriver.LogInfo
		})
	L.PanicIf(err, `qldbdriver.New `+ledger)
	return driver
}

func (a *Adapter) QMapArray(query string, eachRowFunc func(row M.SX) (exitEarly bool)) bool {
	_, err := a.Execute(context.Background(), func(txn qldbdriver.Transaction) (interface{}, error) {
		tables, err := txn.Execute(query)
		if L.IsError(err, `QMapArray.txn.Execute: `+query) {
			return nil, err
		}
		for tables.Next(txn) {
			ionBinary := tables.GetCurrentData()
			row := M.SX{}
			err := ion.Unmarshal(ionBinary, &row)
			if L.IsError(err, `QMapArray.ion.Unmarshall: `+query) {
				return nil, err
			}
			if eachRowFunc(row) {
				return nil, errors.New(`QMapArray.eachRowFunc.exitEarly`)
			}
		}
		return nil, nil
	})
	return err == nil
}

// TODO: add insert/update/delete example
