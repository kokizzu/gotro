package main

import (
	"github.com/jasonlvhit/gocron"

	"example1/domain"
)

func runCron(d *domain.Domain) {

	_ = gocron.Every(10).Minute().Do(func() {
		// TODO: do cron/periodic action here
		//d.CalculateRank()
	})
	gocron.Start()
}
