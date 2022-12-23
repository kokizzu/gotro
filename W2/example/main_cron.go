package main

import (
	"github.com/jasonlvhit/gocron"

	"github.com/kokizzu/gotro/W2/example/domain"
)

func runCron(d *domain.Domain) {

	_ = gocron.Every(10).Minute().Do(func() {
		// TODO: do cron/periodic action here
		//d.CalculateRank()
	})
	gocron.Start()
}
