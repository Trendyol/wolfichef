package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	Gitlab GitlabConfig `yaml:"gitlab"`
	Server ServerConfig `yaml:"server"`
}

func (cfg *AppConfig) Setup() {
	err := cleanenv.ReadConfig("configs/config.yml", cfg)
	if err != nil {
		log.Printf("error while reading the config from file: %+v\n", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Printf("error while reading the config from env: %+v\n", err)
	}
	cfg.Server.Setup()
}
