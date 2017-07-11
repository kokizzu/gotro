package main

import (
	"example-complete/sql"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
	"gitlab.com/kokizzu/gokil/T"
	"time"
)

var VERSION string
var PROJECT_NAME string
var DOMAIN string

var dead_count M.SI
var last_hour int64

func init() {
	dead_count = M.SI{}
}

func CheckUrlAndVideos() {

	W.Mailers[``].SendBCC([]string{sql.DEBUGGER_EMAIL, sql.SUPPORT_EMAIL}, `[`+sql.PROJECT_NAME+`] Archive Status`, `
Dear Administrator,<br/>
<br/>
these videos and urls has been archived, <br/>
<br/>
CHANGE_ME
Best Regards,<br/>
Automated Software<br/>
`)
}

func main() {
	sql.PROJECT_NAME = PROJECT_NAME
	sql.DOMAIN = DOMAIN
	defer sql.ErrorReport(0, `example-cron Internal Server Error: `+VERSION)
	// event loop
	now := time.Now()
	last_minutely_event := now
	last_15minutely_event := now
	last_hourly_event := now
	var dur time.Duration
	CheckUrlAndVideos()
	for {
		//len := sql.QueueLen(mNotifications.REDIS_KEY)
		//if len == 0 {
		time.Sleep(time.Second)
		now := time.Now()
		// send follower count notification
		// send no new content reminder
		now = time.Now()
		dur = now.Sub(last_hourly_event)
		if dur.Minutes() > 60 {
			// do something
			last_hourly_event = now
		}
		// upload to youtube
		dur = now.Sub(last_15minutely_event)
		if dur.Minutes() > 30 {
			// do something
			CheckUrlAndVideos()
			last_15minutely_event = now
		}
		// send public chat notification
		now = time.Now()
		dur = now.Sub(last_minutely_event)
		if dur.Seconds() > 60 {
			// do something
			last_minutely_event = now
		}
		//	continue
		//}
		//task, is_err := sql.QueuedMSX(mNotifications.REDIS_KEY)
		//if is_err {
		//	continue
		//}
		//// TODO: handle send notif
		//_ = task
		//sql.DequeueMSX(mNotifications.REDIS_KEY)
	}
}
