package Es

import (
	"context"
	"github.com/kokizzu/gotro/L"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

type Adapter struct {
	*elastic.Client
	Reconnect func() *elastic.Client
}

func Connect1(url string, debug bool) *elastic.Client {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetSniff(false),
	}
	if debug {
		opts = append(opts,
			elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
		)
	}
	esClient, err := elastic.NewClient(opts...)
	L.PanicIf(err, `elastic.NewClient`)
	_, _, err = esClient.Ping(url).Do(context.Background())
	L.PanicIf(err, `elastic.Ping`)
	return esClient
}
