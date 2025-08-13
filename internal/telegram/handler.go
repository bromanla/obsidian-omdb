package telegram

import (
	"context"
	"fmt"
	"obsidian/omdb/internal/omdb"
	"obsidian/omdb/internal/template"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

type Handlers struct {
	omdb     *omdb.Client
	template *template.Client
}

func NewHandlers(omdbClient *omdb.Client, templateClient *template.Client) *Handlers {
	return &Handlers{
		omdb:     omdbClient,
		template: templateClient,
	}
}

func (h *Handlers) MovieCommand(c telebot.Context) error {
	query := c.Message().Payload
	if query == "" {
		return c.Send("Please provide a movie name or IMDb ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movies, err := h.omdb.Find(ctx, query)
	if err != nil {
		return c.Send("Error searching for movies: " + err.Error())
	}

	if len(movies) == 0 {
		return c.Send("No results found for: " + query)
	}

	if len(movies) > 5 {
		movies = movies[:5]
	}

	replyMarkup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	for _, m := range movies {
		btn := replyMarkup.Data(m.Header(), m.ImdbID)
		rows = append(rows, replyMarkup.Row(btn))
	}

	replyMarkup.Inline(rows...)
	msg := fmt.Sprintf("Search results for: %s", query)

	return c.Send(msg, replyMarkup)
}

func (h *Handlers) Callback(c telebot.Context) error {
	data := c.Callback().Data
	data = strings.TrimPrefix(data, "\f")

	var err error

	if strings.HasPrefix(data, "tt") {
		err = h.preview(c, data)
	}

	if imdbId, ok := strings.CutPrefix(data, "confirm_"); ok {
		err = h.complete(c, imdbId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *Handlers) preview(c telebot.Context, imdbId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movie, err := h.omdb.GetByID(ctx, imdbId)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Error fetching movie details: " + err.Error(),
		})
	}

	confirmBtn := telebot.InlineButton{
		Text: "Confirm",
		Data: "confirm_" + movie.ImdbID,
	}
	keyboard := [][]telebot.InlineButton{{confirmBtn}}
	replyMarkup := &telebot.ReplyMarkup{InlineKeyboard: keyboard}

	msg := fmt.Sprintf(
		"🎬 *%s* (%s)\n\n"+
			"🆔 IMDb ID: %s\n"+
			"📽 Type: %s\n\n"+
			"📖 Plot:\n%s",
		movie.Title,
		movie.Year,
		movie.ImdbID,
		movie.Type,
		movie.Plot,
	)

	photo := &telebot.Photo{
		File:    telebot.FromURL(movie.Poster),
		Caption: msg,
	}

	return c.Edit(photo, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: replyMarkup,
	})
}

func (h *Handlers) complete(c telebot.Context, imdbId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movie, err := h.omdb.GetByID(ctx, imdbId)
	if err != nil {
		return c.Send("Error fetching movie details: " + err.Error())
	}

	if err := h.template.Run(movie); err != nil {
		return c.Send("Failed to save movie: " + err.Error())
	}

	msg := c.Callback().Message.Caption + "\n\n ✅ Saved"
	err = c.EditCaption(msg, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: nil,
	})
	if err != nil {
		return err
	}

	return nil
}
