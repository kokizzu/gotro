package Ql

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/v2/qldbdriver"
	"github.com/kokizzu/gotro/L"
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
