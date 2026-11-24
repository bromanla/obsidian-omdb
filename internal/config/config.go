package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the bot's startup settings.
type Config struct {
	// OMDb supplies metadata for titles identified by IMDb ID.
	OmdbToken string

	// TMDB supplies localized titles and descriptions.
	TmdbToken    string
	TmdbLanguage string

	TelegramToken string
	TelegramUsers []int64 // Users allowed to search for titles and save notes.

	// Directory where the bot saves Markdown files.
	OutputDir string
}

// Load reads and validates the settings needed at startup.
func Load() (Config, error) {
	var env envReader

	cfg := Config{
		OmdbToken:     env.string("OMDB_TOKEN"),
		TmdbToken:     env.string("TMDB_TOKEN"),
		TmdbLanguage:  env.string("TMDB_LANGUAGE", "en-US"),
		TelegramToken: env.string("TELEGRAM_TOKEN"),
		TelegramUsers: env.int64List("TELEGRAM_USERS"),
		OutputDir:     env.string("OUTPUT_DIR", "/data"),
	}

	if err := env.err(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

type envReader struct {
	errs []error
}

func (env *envReader) err() error {
	return errors.Join(env.errs...)
}

func (env *envReader) string(name string, fallback ...string) string {
	value, ok := os.LookupEnv(name)
	if len(fallback) > 0 {
		if !ok {
			return fallback[0]
		}
		return value
	}

	value = strings.TrimSpace(value)
	if value == "" {
		env.errs = append(env.errs, fmt.Errorf("%s is required", name))
	}
	return value
}

func (env *envReader) int64List(name string) []int64 {
	value := env.string(name)
	if value == "" {
		// The required string lookup already recorded the error.
		return nil
	}

	var numbers []int64
	for part := range strings.SplitSeq(value, ",") {
		number, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			env.errs = append(env.errs, fmt.Errorf("%s must contain comma-separated integers", name))
			return nil
		}
		numbers = append(numbers, number)
	}

	return numbers
}
