# ChatWrites

1. Messages come in and are filtered/transformed in [`ChatWrites`](./lib/chat_writes.ex) and then sent to `ChatWrites.MessageCollector`.

2. `ChatWrites.MessageCollector` gets the most frequent message every `CHAT_INTERVAL` milliseconds (env var, defaults to `5000`) and sends them to the TCP client from `ChatWrites.TCPServer`.

3. Edit the outgoing TCP message in [`ChatWrites.MessageCollector.outgoing_message/2`](./lib/message_collector.ex#L73-L75):

```elixir
  defp outgoing_message(count, message) do
    "Out of #{count}, the most common message: #{most_common}"
  end
```
