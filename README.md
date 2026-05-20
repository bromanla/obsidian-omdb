# 🎥 Obsidian OMDB

## About

I’m continuing a series of small projects aimed at learning the Go language and its ecosystem.

For a long time, I’ve been using Obsidian to keep personal notes. While exploring available resources, I came across a video by [FromSergio](https://www.youtube.com/watch?v=t-hKCgGhQuk), where he explains how to search and add movies from OMDb into your vault.

Although I used this solution for some time, my preference for simple, clear, and “clean” tools turned out to be more important. That’s how the idea was born to move this functionality into a standalone service with convenient interaction through a Telegram bot.

At first glance, the solution might seem excessive, but I have two solid reasons:

1. I’m building it primarily for personal use
2. I have plans for further development of the project

For synchronization, I use Syncthing on a server that’s always online. The application is designed with this setup in mind. New movies will be added to the library and synchronized via Syncthing to all other devices.

## 🌐 Environment Variables

Create a `.env` file in your project folder with these:

| Variable           | Description                           |
| ------------------ | ------------------------------------- |
| OMDB_KEY           | OMDb API key for movie searches       |
| TELEGRAM_KEY       | Telegram Bot API token                |
| TELEGRAM_ADMIN     | Telegram user ID for access whitelist |
| OBSIDIAN_DATA_PATH | Path to Obsidian vault on the host    |

## 🐳 Docker Compose

1. Clone repository:

```sh
git clone https://github.com/bromanla/obsidian-omdb
cd obsidian-omdb
```

2. Create `.env` file:

```sh
cp .env.example .env
```

3. Start service:

```sh
docker compose up -d
```

Your Obsidian vault will be mounted into the container automatically.

---

## 🐳 Docker

You can also run the container directly without Docker Compose:

```sh
docker build -t omdb-bot .
```

```sh
docker run -d \
  --name omdb-bot \
  --restart unless-stopped \
  --env-file .env \
  -v ./data:/data/obsidian \
  omdb-bot
```

### ⚙️ Usage

1. In Telegram, use: /movie <query> (e.g., /movie "The Matrix", /movie tt0133093).
2. Select movie from inline buttons.
3. Preview details and poster.
4. Confirm to save Markdown note to Obsidian vault.
