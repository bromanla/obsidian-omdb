package telegram

import (
	"obsidian/omdb/internal/config"
	"obsidian/omdb/internal/omdb"
	"obsidian/omdb/internal/template"
	"time"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func New(omdbClient *omdb.Client, templateClient *template.Client) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  config.Get().TelegramKey,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return bot, err
	}

	handler := NewHandlers(omdbClient, templateClient)

	// bot.Use(middleware.Logger())
	bot.Use(middleware.Whitelist(config.Get().TelegramAdmin))

	bot.Handle("/movie", handler.MovieCommand)
	bot.Handle(telebot.OnCallback, handler.Callback)

	return bot, nil
}
