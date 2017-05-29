package main

import (
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/handler"
	"github.com/kokizzu/gotro/W/example-complex/handler/hBackoffice"
)

var ROUTERS = map[string]W.Action{
	``: handler.Home,

	`backoffice/users`: hBackoffice.Users,

	`login`:              handler.Login,
	`login/forgot`:       handler.Login_Forgot,
	`login/reset/:key`:   handler.Login_Reset,
	`login/verify/:from`: handler.Login_Verify,
	`login/verify`:       handler.Login_Verify,
	`logout`:             handler.Logout,
}
