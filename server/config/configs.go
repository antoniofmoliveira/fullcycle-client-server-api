package config

import (
	"log"

	"github.com/joho/godotenv"
)

type LocalConfig struct {
	Port  string
	Host  string
	DbUrl string
}

func NewConfig() LocalConfig {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	localConfig := LocalConfig{
		Port:  env["PORT"],
		Host:  env["HOST"],
		DbUrl: env["DB_URL"],
	}

	return localConfig
}
