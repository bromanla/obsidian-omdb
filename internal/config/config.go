package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	OmdbKey       string `env:"OMDB_KEY" env-required:"true"`
	TelegramKey   string `env:"TELEGRAM_KEY" env-required:"true"`
	TelegramAdmin int64  `env:"TELEGRAM_ADMIN" env-required:"true"`
	ObsidianPath  string `env:"OBSIDIAN_PATH" env-required:"true"`
}

var cfg Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARNING] .env file not found, using environment variables")
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Panicln("[ERROR] failed to read environment variables:", err)
	}
}

func Get() *Config {
	return &cfg
}
