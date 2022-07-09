
local client = vim.loop.new_tcp()
vim.loop.tcp_connect(client, "127.0.0.1", 6969, function (err)
    assert(not err, err)
    vim.loop.read_start(client, function (err, chunk)
        assert(not err, err)
        print("HELP ME TOM CRUISE", chunk)
    end)
end)


