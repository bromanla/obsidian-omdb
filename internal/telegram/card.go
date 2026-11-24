package telegram

import (
	"fmt"
	"html"
	"strings"

	"obsidian/omdb/internal/catalog"

	"gopkg.in/telebot.v4"
)

const (
	captionLimit = 1024
	messageLimit = 4096
)

func (h *handlers) sendCard(to telebot.Recipient, card catalog.Card, footer string, previewID int) error {
	options := htmlOptions(previewKeyboard(card.Kind, previewID))
	if url := knownValue(card.Poster, ""); url != "" {
		photo := &telebot.Photo{
			File:    telebot.FromURL(url),
			Caption: renderCard(card, footer, captionLimit),
		}
		if _, err := h.bot.Send(to, photo, options); err == nil {
			return nil
		}
	}
	_, err := h.bot.Send(to, renderCard(card, footer, messageLimit), options)
	return err
}

func editCard(c telebot.Context, card catalog.Card, footer string, keyboard *telebot.ReplyMarkup) error {
	options := htmlOptions(keyboard)
	if c.Message().Photo != nil {
		return c.EditCaption(renderCard(card, footer, captionLimit), options)
	}
	return c.Edit(renderCard(card, footer, messageLimit), options)
}

func renderCard(card catalog.Card, footer string, limit int) string {
	title := strings.Join(strings.Fields(knownValue(card.Title, "Unknown title")), " ")
	title = shorten(title, 256)
	year := shorten(knownValue(card.Year, "Unknown year"), 32)
	metadata := fmt.Sprintf("%s · %s · IMDb ", card.Kind.Label(), year)
	plotLimit := limit - len(title+"\n"+metadata+card.ImdbID+"\n\n"+footer)
	plot := shorten(knownValue(card.Plot, "No description available."), plotLimit)
	return fmt.Sprintf("<b>%s</b>\n%s<a href=\"https://www.imdb.com/title/%s/\">%s</a>\n\n%s%s",
		html.EscapeString(title), html.EscapeString(metadata), card.ImdbID, card.ImdbID,
		html.EscapeString(plot), html.EscapeString(footer),
	)
}

func shorten(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	words := strings.Fields(text)
	length := 0
	for i, word := range words {
		length += len(word) + 1
		if length > limit-len("…") {
			return strings.Join(words[:i], " ") + "…"
		}
	}
	return strings.Join(words, " ")
}

func knownValue(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "N/A" {
		return fallback
	}
	return value
}

func htmlOptions(keyboard *telebot.ReplyMarkup) *telebot.SendOptions {
	return &telebot.SendOptions{ParseMode: telebot.ModeHTML, DisableWebPagePreview: true, ReplyMarkup: keyboard}
}

func previewKeyboard(kind catalog.Kind, tmdbID int) *telebot.ReplyMarkup {
	if tmdbID == 0 {
		return nil
	}
	return &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{
		button("« Back", "back", kind, 0),
		button("Confirm", "save", kind, tmdbID),
	}}}
}
