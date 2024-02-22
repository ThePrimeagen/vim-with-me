local uv = vim.loop

local function read(client)
    print("read", vim.inspect(client))
    uv.read_start(client, function(e, chunk)
        print("read start", e, chunk)
    end)
end

local client = uv.new_tcp()
client:connect("127.0.0.1", 42069, function(err)
    print("connect", err)
    read(client)
end)

function CLOSE()
    client:shutdown()
    client:close()
end
