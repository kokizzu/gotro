package W

import "github.com/kokizzu/gotro/M"

var Mailers map[string]*SmtpConfig // used in mailer.go
var Webmasters M.SS                // used in engine.go
var Sessions SessionConnector      // from session.go
var Routes map[string]Action
var Assets [][2]string // []{{`js css /js /css`,`filename`}, ...}
var Filters []Action
