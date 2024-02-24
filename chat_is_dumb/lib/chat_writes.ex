defmodule ChatWrites do
  @moduledoc """
  The twitch chat bot event handler.
  We handle the message event here and add the message text to the message collector.
  """
  use TwitchChat.Bot

  alias ChatWrites.MessageCollector
  alias TwitchChat.Events.Message

  @transform_allowed %{
    "<backspace>" => "\b",
    "<space>" => " ",
    "<cr>" => "\n",
    "<dot>" => ".",
    "<esc>" => "\e",
    "<tab>" => "\t"
  }

  @allowed_strings Map.keys(@transform_allowed)

  @printable_ascii 32..127

  @doc """
  Handle events from Twitch Chat.
  Takes allowed messages and adds them to the `MessageCollector`.
  """
  @impl TwitchChat.Bot
  def handle_event(%Message{message: special}) when special in @allowed_strings do
    special
    |> transform_special()
    |> MessageCollector.add()
  end

  def handle_event(%Message{message: <<c::8>>}) when c in @printable_ascii do
    MessageCollector.add(c)
  end

  @doc """
  Pattern-match on the special string to return the transformed version used in
  the editor.

  ## Examples

      iex> transform_special("<cr>")
      "\\n"

      iex> transform_special("<dot>")
      "."

  """
  for {special, transform} <- @transform_allowed do
    def transform_special(unquote(special)), do: unquote(transform)
  end
end
