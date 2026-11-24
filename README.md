# 🎥 Obsidian OMDB

## About

I’m continuing a series of small projects aimed at learning the Go language and its ecosystem.

For a long time, I’ve been using Obsidian to keep personal notes. While exploring available resources, I came across a video by [FromSergio](https://www.youtube.com/watch?v=t-hKCgGhQuk), where he explains how to search and add movies from OMDb into your vault.

Although I used this solution for some time, my preference for simple, clear, and “clean” tools turned out to be more important. That’s how the idea was born to move this functionality into a standalone service with convenient interaction through a Telegram bot.

At first glance, the solution might seem excessive, but I have two solid reasons:

1. I’m building it primarily for personal use
2. I have plans for further development of the project

For synchronization, I use Syncthing on a server that’s always online. The application is designed with this setup in mind. New movies will be added to the library and synchronized via Syncthing to all other devices.

## Run with Docker Compose

With Docker and Compose installed, clone the repository and create the configuration and output directory.

```sh
git clone https://github.com/bromanla/obsidian-omdb
cd obsidian-omdb
cp .env.example .env
mkdir -p data
```

Fill in the three API tokens and your Telegram user ID in `.env`, then build and start the bot.

```sh
docker compose up --build -d
```

Files are saved to `./data` on the host. To use another directory, create it and set `HOST_OUTPUT_DIR` in `.env`. Compose mounts the host directory at `/data` and sets `OUTPUT_DIR=/data` in the container.

Omit `-d` to run in the foreground. For a background run, use `docker compose logs -f` to follow the logs and `docker compose down` to stop the bot.

## Configuration

| Variable          | Description                                                                          |
| ----------------- | ------------------------------------------------------------------------------------ |
| `OMDB_TOKEN`      | OMDb API key for looking up details by IMDb ID                                       |
| `TMDB_TOKEN`      | TMDB API Read Access Token                                                           |
| `TELEGRAM_TOKEN`  | Telegram Bot API token                                                               |
| `TELEGRAM_USERS`  | Allowed Telegram user IDs, separated by commas                                       |
| `TMDB_LANGUAGE`   | Language for titles and descriptions. Defaults to `en-US`. For Russian, use `ru-RU`. |
| `OUTPUT_DIR`      | Directory for files when running the bot locally. Defaults to `/data`.               |
| `HOST_OUTPUT_DIR` | Directory on the host mounted by Docker Compose. Defaults to `./data`.               |

## Use the bot

1. Open a private chat with the bot and send `/start`.
2. Send a title, optionally ending with a year in square brackets, such as `Dune [2021]`.
3. Choose **Movie** or **TV show**, then select a result if several are found.
4. Press **Confirm** to save the card or **Back** to return to the results.

After saving, the card shows `✅ Saved`. Other allowed users who have opened a chat with the bot receive a copy. If a file with the same name already exists, the card shows `ℹ️ File already exists`.

## Local development

Install Go 1.27.4 or later and Make. Create `.env` from `.env.example` if needed, fill in the credentials, and set `OUTPUT_DIR=./data`. Run these commands from the project root.

```sh
mkdir -p data
make dev
```

On its first run, `make dev` installs the pinned version of Air in `.air/bin`. Air loads `.env` and rebuilds the bot when Go files, the Markdown template, or module files change. Restart `make dev` after editing `.env` or `.air.toml`. [Air](https://github.com/air-verse/air)

Use `make fmt` to format the Go files and `make` to list commands. For `go run ./cmd` or a direct binary launch, export the settings in your shell first. Use a separate bot token if a deployed instance is already polling Telegram.

## API requests in Bruno

Run `./bin/setup-bruno.sh` to make the credentials in `.env` available to the collection. In Bruno, choose **Open collection** and select `bruno/`.
