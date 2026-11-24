package template

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"obsidian/omdb/internal/catalog"
)

//go:embed assets/movie.md
var templateContent string

var invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

type Store struct {
	tmpl   *template.Template
	folder string
	mu     sync.Mutex
}

type SaveResult struct {
	Filename      string
	AlreadyExists bool
}

func New(folder string) (*Store, error) {
	funcMap := template.FuncMap{
		"split": strings.Split,
		"trim":  strings.TrimSpace,
		"now": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	tmpl, err := template.
		New("movie.md").
		Funcs(funcMap).
		Parse(templateContent)

	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	return &Store{tmpl: tmpl, folder: folder}, nil
}

func (s *Store) Save(card catalog.Card) (SaveResult, error) {
	filename, err := markdownFilename(card.Title)
	if err != nil {
		return SaveResult{}, err
	}

	content, err := s.renderMarkdown(card)
	if err != nil {
		return SaveResult{}, err
	}
	return s.writeFile(content, filename)
}

func (s *Store) renderMarkdown(card catalog.Card) ([]byte, error) {
	var content bytes.Buffer
	if err := s.tmpl.Execute(&content, card); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}
	return content.Bytes(), nil
}

func (s *Store) writeFile(content []byte, filename string) (SaveResult, error) {
	path := filepath.Join(s.folder, filename)
	// Finish writing before another handler can report that the file exists.
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if errors.Is(err, os.ErrExist) {
		return SaveResult{Filename: filename, AlreadyExists: true}, nil
	}
	if err != nil {
		return SaveResult{}, fmt.Errorf("create file %q: %w", filename, err)
	}
	_, writeErr := file.Write(content)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return SaveResult{}, fmt.Errorf("write file %q: %w", filename, errors.Join(err, os.Remove(path)))
	}
	return SaveResult{Filename: filename}, nil
}

func markdownFilename(title string) (string, error) {
	name := strings.TrimSpace(invalidFilenameChars.ReplaceAllString(title, ""))
	name = strings.Trim(name, ". ")
	if name == "" {
		return "", errors.New("title has no valid filename characters")
	}
	return name + ".md", nil
}
