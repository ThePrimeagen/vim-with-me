import Config

# Twitch chat bot config.
config :chat_writes,
  bot: [
    bot: ChatWrites,
    user: System.fetch_env!("TWITCH_USER"),
    pass: System.fetch_env!("TWITCH_OAUTH_TOKEN"),
    channels: [System.fetch_env!("TWITCH_USER")]
  ]

# TCP Server config.
config :chat_writes, ChatWrites.TCPServer,
  start?: System.get_env("TCP_SERVER", "false") == "true",
  port: System.get_env("PORT", "42069") |> String.to_integer()

# Chat message collector interval in milliseconds.
config :chat_writes, ChatWrites.MessageCollector,
  interval_ms: System.get_env("CHAT_INTERVAL", "5000") |> String.to_integer()

if config_env() == :test do
  # Don't show info logs in tests.
  config :logger, :console, level: :warning
else
  # Logger config. Valid values are: `debug, info, warning, error`.
  config :logger, :console,
    level: System.get_env("LOG_LEVEL", "info") |> String.to_existing_atom()
end
