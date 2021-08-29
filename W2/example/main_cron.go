package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/kokizzu/gotro/W2/example/domain"
)

func runCron(d *domain.Domain) {

	// TODO: increment delay based on total player
	gocron.Every(10).Minute().Do(func() {
		//d.CalculateRank()
	})
	gocron.Start()
}
