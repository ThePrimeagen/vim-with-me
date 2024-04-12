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

---@param opts TCPOptions | nil
local function get_opts(opts)
    return vim.tbl_extend("force", {}, copy_opts(default_opts), opts or {})
end

---@class TCPCommand
---@field command number
---@field data string

---@alias TCPListener (fun(command: TCPCommand): nil)

---@class TCP
---@field _connection any | nil
---@field _listeners TCPListener[]
---@field opts TCPOptions
local TCP = {}
TCP.__index = TCP

---@param opts TCPOptions | nil
function TCP:new(opts)
    opts = get_opts(opts)
    return setmetatable({
        _connection = nil,
        _listeners = {},
        opts = opts,
    }, self)
end

---@param cb fun(err: string | nil): nil
---@param opts TCPOptions | nil
function TCP:start(cb, opts)
    opts = get_opts(opts or self.opts)

    assert(self._connection == nil, "client already started")
    assert(opts.retry_count > 0, "could not connect to server")

    cb = cb or function() end

    self._connection = uv.new_tcp()
    self._connection:connect(self.opts.host, self.opts.port, function(err)
        if err then
            self._connection = nil
            vim.defer_fn(function()
                local next_opts = copy_opts(opts)
                next_opts.retry_count = next_opts.retry_count - 1
                self:start(cb, opts)
            end, opts.retry_wait_ms)
            return
        end

        self:_read()
        vim.schedule(cb)
    end)
end

function TCP:_read()
    assert(self._connection, "client not started")

    local process = TcpProcess.process_packets()
    uv.read_start(
        self._connection,
        vim.schedule_wrap(function(_, chunk)
            while true do
                local command = process(chunk)
                chunk = ""
                if command == nil then
                    break
                end

                for _, listener in ipairs(self._listeners) do
                    listener(command)
                end
            end
        end)
    )
end

function TCP:close()
    self._listeners = {}

    --- no assert here as i am going to call this a lot when devving and
    --- i may or may not have a running server
    if self._connection == nil then
        return
    end

    self._connection:shutdown()
    self._connection:close()
    self._connection = nil
end

---@return boolean
function TCP:connected()
    return self._connection ~= nil
end

---@param command TCPCommand
function TCP:send(command)
    assert(self._connection, "client not started")
    local ok, _ = pcall(
        self._connection.write,
        self._connection,
        TcpProcess.encode_tcp_command(command)
    )

    -- TODO: Probably turn this into something different?
    assert(ok, "could not send data")
end

---@param cb TCPListener
function TCP:listen(cb)
    table.insert(self._listeners, cb)
end

return TCP
