package main

import (
	"example-complete/handler"
	"example-complete/handler/uSuperAdmin"
	"github.com/kokizzu/gotro/W"
)

var HANDLERS = map[string]W.Action{
	``: handler.Home,

	`login`:              handler.Login,
	`login/verify/:from`: handler.Login_Verify,
	`login/verify`:       handler.Login_Verify,
	`logout`:             handler.Logout,

	`superadmin/users`: uSuperAdmin.Users,
}
