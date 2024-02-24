defmodule ChatWrites.TCPServer do
  @moduledoc """
  The TCP server.
  """
  use GenServer

  require Logger

  @default_port 42069

  @doc """
  Start accepting connections.

  ## Options

   * `:port` - The port to start accepting connections on (defaults to `42069`).

  """
  def start_link(opts) do
    if Keyword.get(opts, :start?, true) do
      GenServer.start_link(__MODULE__, opts, name: __MODULE__)
    else
      :ignore
    end
  end

  @doc """
  Send a message from the server to the connected client.
  """
  def send(message) do
    GenServer.cast(__MODULE__, {:send, message})
  end

  # ----------------------------------------------------------------------------
  # GenServer callbacks
  # ----------------------------------------------------------------------------

  @impl GenServer
  def init(opts) do
    port = Keyword.get(opts, :port, @default_port)

    # The options below mean:
    #
    # 1. `:binary` - receives data as binaries (instead of lists).
    # 3. `active: false` - blocks on `:gen_tcp.recv/2` until data is available.
    # 4. `reuseaddr: true` - allows us to reuse the address if the listener crashes.
    #
    {:ok, socket} = :gen_tcp.listen(port, [:binary, active: true, reuseaddr: true])
    Logger.info("[ChatWrites.TCPServer] accepting connections on port #{port}...")

    state = %{socket: socket, client: nil}

    {:ok, state, {:continue, :accept}}
  end

  @impl GenServer
  def handle_continue(:accept, state) do
    send(self(), :accept)
    {:noreply, state}
  end

  @impl GenServer
  def handle_info(:accept, state) do
    {:ok, client} = :gen_tcp.accept(state.socket)
    Logger.info("[ChatWrites.TCPServer] Client connected")
    {:noreply, %{state | client: client}}
  end

  def handle_info({:tcp, socket, data}, state) do
    Logger.info("[ChatWrites.TCPServer] Received #{data}")
    Logger.info("[ChatWrites.TCPServer] Sending it back")

    :ok = :gen_tcp.send(socket, data)

    {:noreply, state}
  end

  def handle_info({:tcp_closed, _}, state), do: {:stop, :normal, state}
  def handle_info({:tcp_error, _}, state), do: {:stop, :normal, state}

  @impl GenServer
  def handle_cast({:send, message}, state) do
    :ok = :gen_tcp.send(state.client, message)
    {:noreply, state}
  end
end
