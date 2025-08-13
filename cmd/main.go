package main

import (
	"log"
	"obsidian/omdb/internal/config"
	"obsidian/omdb/internal/omdb"
	"obsidian/omdb/internal/telegram"
	"obsidian/omdb/internal/template"
)

func main() {
	cfg := config.Get()
	log.Println("[DEBUG] load application")

	omdbClient := omdb.New(cfg.OmdbKey, nil)

	templateClient, err := template.New()
	if err != nil {
		log.Fatalf("[FATAL] failed to init template: %v", err)
	}

	bot, err := telegram.New(omdbClient, templateClient)
	if err != nil {
		log.Fatalf("[FATAL] failed to start telegram bot: %v", err)
	}

	log.Println("[DEBUG] start telegram bot")
	bot.Start()
}
