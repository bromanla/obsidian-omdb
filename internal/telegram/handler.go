package telegram

import (
	"context"
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	"obsidian/omdb/internal/catalog"

	"gopkg.in/telebot.v4"
)

const (
	helpText      = "Send a movie or TV show title.\nYou can include a year, for example Dune [2021]."
	searchLine    = "\n\nSearch · "
	yearLine      = "\nYear · "
	savedMark     = "\n\n✅ Saved"
	saveError     = "\n\n⚠️ Couldn't save the file."
	maxQueryBytes = 512
)

func (h *handlers) handleMessage(c telebot.Context) error {
	text := strings.Join(strings.Fields(c.Text()), " ")
	if text == "" || strings.HasPrefix(text, "/") {
		return c.Send(helpText)
	}
	if len(text) > maxQueryBytes {
		return c.Send("That title is too long. Send a shorter title.")
	}

	keyboard := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{
		button(catalog.Movie.Label(), "search", catalog.Movie, 0),
		button(catalog.TVShow.Label(), "search", catalog.TVShow, 0),
	}}}
	return sendSearch(c, "Choose a type", parseQuery(text), keyboard)
}

func (h *handlers) handleCallback(c telebot.Context) error {
	if err := c.RespondText("Loading…"); err != nil {
		h.bot.OnError(err, c)
	}
	parts := strings.Split(c.Callback().Data, ":")
	if len(parts) != 3 || c.Message() == nil {
		return nil
	}
	kind := catalog.Kind(parts[1])
	if kind != catalog.Movie && kind != catalog.TVShow {
		return nil
	}
	tmdbID, err := strconv.Atoi(parts[2])
	if err != nil || tmdbID < 0 {
		return nil
	}
	_, query := splitSearchFooter(c.Text())
	if query.Title == "" {
		return nil
	}

	switch parts[0] {
	case "search", "back":
		err = h.search(c, query, kind, parts[0] == "back")
	case "preview":
		err = h.preview(c, query, kind, tmdbID)
	case "save":
		return h.save(c, query, kind, tmdbID)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	return c.Delete()
}

func (h *handlers) search(c telebot.Context, query catalog.Query, kind catalog.Kind, showList bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	results, err := h.catalog.Search(ctx, query, kind)
	if err != nil {
		h.bot.OnError(err, c)
		return c.Send("Couldn't complete the search.")
	}
	if len(results) == 0 {
		return c.Send("No results found.")
	}
	if len(results) == 1 && !showList {
		return h.preview(c, query, kind, results[0].ID)
	}

	keyboard := &telebot.ReplyMarkup{}
	for _, result := range results {
		label := fmt.Sprintf("%s (%s) · %s", result.Title, knownValue(result.Year, "Unknown year"), kind.Label())
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []telebot.InlineButton{
			button(label, "preview", kind, result.ID),
		})
	}
	return sendSearch(c, "Choose a result", query, keyboard)
}

func (h *handlers) preview(c telebot.Context, query catalog.Query, kind catalog.Kind, tmdbID int) error {
	card, err := h.loadCard(kind, tmdbID)
	if err != nil {
		h.bot.OnError(err, c)
		return c.Send("Couldn't load title details.")
	}
	return h.sendCard(c.Chat(), card, searchFooter(query), tmdbID)
}

func (h *handlers) save(c telebot.Context, query catalog.Query, kind catalog.Kind, tmdbID int) error {
	card, err := h.loadCard(kind, tmdbID)
	if err != nil {
		h.bot.OnError(err, c)
		return c.Send("Couldn't load title details.")
	}
	result, err := h.store.Save(card)
	if err != nil {
		h.bot.OnError(err, c)
		return editCard(c, card, saveError+searchFooter(query), previewKeyboard(kind, tmdbID))
	}
	if result.AlreadyExists {
		footer := "\n\nℹ️ File already exists" + searchFooter(query)
		return editCard(c, card, footer, previewKeyboard(kind, tmdbID))
	}

	log.Printf("[INFO] Telegram user %d saved %q", c.Sender().ID, result.Filename)
	err = editCard(c, card, savedMark, &telebot.ReplyMarkup{})
	for _, id := range h.users {
		if id == c.Sender().ID {
			continue
		}
		if sendErr := h.sendCard(&telebot.User{ID: id}, card, savedMark, 0); sendErr != nil {
			h.bot.OnError(sendErr, c)
		}
	}
	return err
}

func (h *handlers) loadCard(kind catalog.Kind, tmdbID int) (catalog.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return h.catalog.Details(ctx, kind, tmdbID)
}

func sendSearch(c telebot.Context, heading string, query catalog.Query, keyboard *telebot.ReplyMarkup) error {
	return c.Send("<b>"+heading+"</b>"+html.EscapeString(searchFooter(query)), htmlOptions(keyboard))
}

func button(label, action string, kind catalog.Kind, tmdbID int) telebot.InlineButton {
	return telebot.InlineButton{Text: label, Data: fmt.Sprintf("%s:%s:%d", action, kind, tmdbID)}
}

func parseQuery(text string) catalog.Query {
	text = strings.TrimSpace(text)
	title, year, _ := strings.CutLast(text, "[")
	title = strings.TrimSpace(title)
	year, closed := strings.CutSuffix(year, "]")
	if title != "" && closed && len(year) == 4 && strings.Trim(year, "0123456789") == "" {
		return catalog.Query{Title: title, Year: year}
	}
	return catalog.Query{Title: text}
}

// Keep the query in the message so Back works after a restart.
func searchFooter(query catalog.Query) string {
	footer := searchLine + query.Title
	if query.Year != "" {
		footer += yearLine + query.Year
	}
	return footer
}

func splitSearchFooter(text string) (string, catalog.Query) {
	body, footer, _ := strings.CutLast(text, searchLine)
	title, year, _ := strings.Cut(footer, yearLine)
	return body, catalog.Query{Title: title, Year: year}
}
