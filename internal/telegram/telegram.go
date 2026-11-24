package telegram

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"obsidian/omdb/internal/catalog"
	"obsidian/omdb/internal/template"

	"gopkg.in/telebot.v4"
)

type handlers struct {
	bot     *telebot.Bot
	catalog *catalog.Service
	store   *template.Store
	users   []int64
}

func New(token string, users []int64, catalogService *catalog.Service, store *template.Store) (*telebot.Bot, error) {
	redact := func(err error) string {
		return strings.ReplaceAll(err.Error(), token, "[redacted]")
	}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Client: &http.Client{Timeout: 20 * time.Second},
		OnError: func(err error, _ telebot.Context) {
			log.Printf("[ERROR] Telegram: %s", redact(err))
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect to Telegram: %s", redact(err))
	}

	registerHandlers(bot, users, catalogService, store)
	return bot, nil
}

func registerHandlers(bot *telebot.Bot, users []int64, catalogService *catalog.Service, store *template.Store) {
	h := &handlers{
		bot:     bot,
		catalog: catalogService,
		store:   store,
		users:   users,
	}
	bot.Use(h.allowPrivateUsers)

	bot.Handle("/start", h.handleMessage)
	for _, endpoint := range []string{
		telebot.OnText, telebot.OnMedia, telebot.OnContact, telebot.OnLocation,
		telebot.OnVenue, telebot.OnDice, telebot.OnGame,
	} {
		bot.Handle(endpoint, h.handleMessage)
	}
	bot.Handle(telebot.OnCallback, h.handleCallback)
}

func (h *handlers) allowPrivateUsers(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		chat, sender := c.Chat(), c.Sender()
		if chat == nil || chat.Type != telebot.ChatPrivate || sender == nil || !slices.Contains(h.users, sender.ID) {
			return nil
		}
		return next(c)
	}
}
