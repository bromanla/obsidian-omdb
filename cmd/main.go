package main

import (
	"log"

	"obsidian/omdb/internal/catalog"
	"obsidian/omdb/internal/config"
	"obsidian/omdb/internal/omdb"
	"obsidian/omdb/internal/telegram"
	"obsidian/omdb/internal/template"
	"obsidian/omdb/internal/tmdb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] failed to load configuration: %v", err)
	}

	omdbClient := omdb.New(cfg.OmdbToken, nil)
	tmdbClient := tmdb.New(cfg.TmdbToken, cfg.TmdbLanguage, nil)
	catalogService := catalog.New(tmdbClient, omdbClient)

	store, err := template.New(cfg.OutputDir)
	if err != nil {
		log.Fatalf("[FATAL] failed to initialize the file store: %v", err)
	}

	bot, err := telegram.New(cfg.TelegramToken, cfg.TelegramUsers, catalogService, store)
	if err != nil {
		log.Fatalf("[FATAL] failed to start telegram bot: %v", err)
	}

	log.Printf("[INFO] Starting Telegram bot for %d user(s)", len(cfg.TelegramUsers))
	bot.Start()
}
