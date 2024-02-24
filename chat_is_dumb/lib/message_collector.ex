defmodule ChatWrites.MessageCollector do
  @moduledoc """
  Collects messages and does something with the most common message for every
  interval `tick`.
  """
  use GenServer

  require Logger

  @default_interval_ms 5000

  @doc """
  Starts the message collector.
  """
  def start_link(opts) do
    GenServer.start_link(__MODULE__, opts, name: __MODULE__)
  end

  @doc """
  Add a message to the collector.
  """
  def add(message) do
    GenServer.cast(__MODULE__, {:add, message})
  end

  # ----------------------------------------------------------------------------
  # GenServer callbacks
  # ----------------------------------------------------------------------------

  @impl GenServer
  def init(opts) do
    interval_ms = Keyword.get(opts, :interval_ms, @default_interval_ms)
    timer_ref = schedule_next(interval_ms)

    state = %{
      interval_ms: interval_ms,
      timer_ref: timer_ref,
      messages: []
    }

    {:ok, state}
  end

  @impl GenServer
  def handle_cast({:add, message}, state) do
    {:noreply, %{state | messages: [message | state.messages]}}
  end

  @impl GenServer
  def handle_info(:tick, state) do
    with [{most_common, _freq} | _] <- get_most_common(state.messages) do
      count = Enum.count(state.messages)
      message = outgoing_message(count, most_common)
      Logger.info(message)
      ChatWrites.TCPServer.send(message)
    end

    timer_ref = schedule_next(state.interval_ms)

    {:noreply, %{state | messages: [], timer_ref: timer_ref}}
  end

  # ----------------------------------------------------------------------------
  # Helpers
  # ----------------------------------------------------------------------------

  defp get_most_common(messages) do
    messages
    |> Enum.frequencies()
    |> Enum.sort_by(fn {_msg, freq} -> freq end, :desc)
  end

  defp outgoing_message(count, message) do
    "Out of #{count}, the most common message: #{most_common}"
  end

  defp schedule_next(interval_ms) do
    Process.send_after(self(), :tick, interval_ms)
  end
end
