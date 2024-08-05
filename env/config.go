package env

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Configuration struct {
	LogLevel              string `env:"LOG_LEVEL"`
}


func GetEnvConfig() *Configuration {
	setupEnv()
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}

func setupEnv() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return
	}

	err := godotenv.Load(fmt.Sprintf("./env/.env.%s", strings.ToLower(appEnv)))
	if err == nil {
		return
	}

	_ = godotenv.Load()
}
