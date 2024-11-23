package main

import (
	"log"
	"wb_tech_l0/cmd/app"
)

func main() {
	a := app.App{}
	if err := a.Init(); err != nil {
		log.Panicf("Failed to init applicaation: %s", err.Error())
	}
	if err := a.Start(); err != nil {
		panic(err)
	}
}
