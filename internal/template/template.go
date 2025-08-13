package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"obsidian/omdb/internal/config"
)

type Client struct {
	tmpl *template.Template
}

type Data interface {
	Sanitize() string
}

func New() (*Client, error) {
	funcMap := template.FuncMap{
		"split": strings.Split,
		"trim":  strings.TrimSpace,
		"now": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	tmpl, err := template.
		New("template.md").
		Funcs(funcMap).
		ParseFiles("template.md")

	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Client{tmpl}, nil
}

func (c *Client) generateTemplate(data any) (string, error) {
	var buf bytes.Buffer

	if err := c.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (c *Client) saveTemplate(content string, filename string) error {
	folder := config.Get().ObsidianPath

	if fi, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("folder does not exist: %s", folder)
		}
		return fmt.Errorf("failed to check folder: %w", err)
	} else if !fi.IsDir() {
		return fmt.Errorf("destination is not a directory: %s", folder)
	}

	path := filepath.Join(folder, filename+".md")
	file, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("fUnable to create file: %w", err)
	}
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Client) Run(data Data) error {
	content, err := c.generateTemplate(data)
	if err != nil {
		return err
	}

	filename := data.Sanitize()
	return c.saveTemplate(content, filename)
}
