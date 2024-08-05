package main

import (
	"log"

	app_log "github.com/Allan-Nava/go-wire-fuego-scafffold/app/lib/log"
	"github.com/Allan-Nava/go-wire-fuego-scafffold/deps"
	"github.com/Allan-Nava/go-wire-fuego-scafffold/env"
)

func main() {

	config := env.GetEnvConfig()
	logger := app_log.NewLogger(config.LogLevel)
	app, err := deps.InjectApp(config, logger)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Routes(logger).Run()
	log.Fatal(err)
}
