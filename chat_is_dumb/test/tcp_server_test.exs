defmodule ChatWrites.TCPServerTest do
  use ExUnit.Case

  alias ChatWrites.TCPServer

  @port 4040

  setup do
    _server_pid = start_supervised!({ChatWrites.TCPServer, [start?: true, port: @port]})
    :ok
  end

  test "accepts connections and sends messages" do
    assert {:ok, client} = ChatWrites.TestTCPClient.start(port: @port)
    :erlang.trace(client, true, [:receive])
    :ok = TCPServer.send("foo")
    assert_receive {:trace, ^client, :receive, {:tcp, _port, "foo"}}
  end
end
