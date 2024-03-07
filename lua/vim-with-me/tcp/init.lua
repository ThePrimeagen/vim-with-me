local TcpProcess = require("vim-with-me.tcp.process")

---@class TCPOptions
---@field host string
---@field port number
---@field retry_wait_ms? number
---@field retry_count? number

---@param opts TCPOptions
---@return TCPOptions
local function copy_opts(opts)
    local out = {}
    for k, v in pairs(opts) do
        out[k] = v
    end
    return out
end

local default_opts = {
    host = "127.0.0.1",
    port = 42069,
    retry_wait_ms = 10000,
    retry_count = 3,
}

-- luacheck: ignore 111
local uv = vim.loop

Existing_TCP_Connection = Existing_TCP_Connection or nil

---@type (fun(command: string, data: string): nil)[]
Existing_TCP_Listeners = Existing_TCP_Listeners or {}

local function read(client)
    local process = TcpProcess.process_packets()
    uv.read_start(
        client,
        vim.schedule_wrap(function(_, chunk)
            while true do
                local command, data = process(chunk)
                chunk = ""
                if command == nil or data == nil then
                    break
                end

                for _, listener in ipairs(Existing_TCP_Listeners) do
                    listener(command, data)
                end
            end
        end)
    )
end

---@param opts TCPOptions | nil
---@param cb fun(): nil
local function tcp_start(opts, cb)
    cb = cb or function() end
    assert(Existing_TCP_Connection == nil, "client already started")

    opts = opts or copy_opts(default_opts)
    if opts.retry_count <= 0 then
        error("Failed to connect to server")
    end

    Existing_TCP_Connection = uv.new_tcp()
    Existing_TCP_Listeners = {}
    Existing_TCP_Connection:connect(opts.host, opts.port, function(err)
        if err then
            Existing_TCP_Connection = nil
            vim.defer_fn(function()
                local next_opts = copy_opts(opts)
                next_opts.retry_count = next_opts.retry_count - 1
                tcp_start(next_opts, cb)
            end, 10000)
            return
        end

        read(Existing_TCP_Connection)
        cb()
    end)
end

local function tcp_stop()
    Existing_TCP_Listeners = {}

    --- no assert here as i am going to call this a lot when devving and
    --- i may or may not have a running server
    if Existing_TCP_Connection == nil then
        return
    end

    Existing_TCP_Connection:shutdown()
    Existing_TCP_Connection:close()
    Existing_TCP_Connection = nil
end

return {
    tcp_start = tcp_start,
    tcp_stop = tcp_stop,
    tcp_connected = function()
        return Existing_TCP_Connection ~= nil
    end,

    ---@param command string
    ---@param data string
    tcp_send = function(command, data)
        assert(Existing_TCP_Connection, "client not started")
        Existing_TCP_Connection:write(TcpProcess.create_tcp_command(command, data))
    end,

    listen = function(listener)
        table.insert(Existing_TCP_Listeners, listener)
    end,
}
