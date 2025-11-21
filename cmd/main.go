package main

import (
	"log"
	"webServerGo/internal/app"
	"webServerGo/internal/pkg/cfg"
)

func main() {
	conf, err := cfg.New()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(conf.Port)
}
