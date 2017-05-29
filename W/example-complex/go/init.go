package main

import "github.com/kokizzu/gotro/W/example-complex/model"

func main() {
	model.PG_W.CreateBaseTable(`users`, `users`)
	model.PG_W.CreateBaseTable(`todos`, `users`) // 2nd table
}
