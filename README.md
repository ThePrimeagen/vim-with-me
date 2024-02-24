# Vim With Me

Send the most common chat messages to neovim through TCP.

## Setup

1. You need the ENV vars found in `chat_is_dumb/config/runtime.exs`

2. If using `docker compose`, copy `cp .env.example .env` and set the twitch ENV vars.

## Starting

### With Docker

```
docker compose up
```

### Without Docker

in `./chat_is_dumb`

```
mix run --no-halt
```

Or starting with the Elixir interactive shell

```
iex -S mix
```
