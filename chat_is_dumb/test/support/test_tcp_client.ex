defmodule ChatWrites.TestTCPClient do
  use GenServer

  require Logger

  @ip {127, 0, 0, 1}
  @default_port 4040

  def start(opts \\ []) do
    GenServer.start(__MODULE__, opts)
  end

  def send_message(pid, message) do
    GenServer.cast(pid, {:message, message})
  end

  def init(opts) do
    send(self(), :connect)
    state = %{port: Keyword.get(opts, :port, @default_port), socket: nil}
    {:ok, state}
  end

  def handle_info(:connect, state) do
    Logger.info("Connecting to #{:inet.ntoa(@ip)}:#{state.port}")

    case :gen_tcp.connect(@ip, state.port, [:binary, active: true]) do
      {:ok, socket} ->
        {:noreply, %{state | socket: socket}}

      {:error, reason} ->
        disconnect(state, reason)
    end
  end

  def handle_info({:tcp, _, data}, state) do
    Logger.info("Received #{data}")

    {:noreply, state}
  end

  def handle_info({:tcp_closed, _}, state), do: {:stop, :normal, state}
  def handle_info({:tcp_error, _}, state), do: {:stop, :normal, state}

  def handle_cast({:message, message}, %{socket: socket} = state) do
    Logger.info("Sending #{message}")

    :ok = :gen_tcp.send(socket, message)
    {:noreply, state}
  end

  def disconnect(state, reason) do
    Logger.info("Disconnected: #{reason}")
    {:stop, :normal, state}
  end
end
