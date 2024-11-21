package main

import "wb_tech_l0/cmd/app"

func main() {
	a := app.App{}
	if err := a.Start(); err != nil {
		panic(err)
	}
}
