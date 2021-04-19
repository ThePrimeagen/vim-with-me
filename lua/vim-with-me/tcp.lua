local states = {
    disconnected = 1,
    connecting = 2,
    connected = 3,
    error = 4,
    timedout = 5,
}

local TcpClient = {}

function TcpClient:new(host, port)
    local obj = {
        host = host,
        port = port,
        client = nil,
        state = states.disconnected,
        callbacks = {},
    }

    setmetatable(obj, self)
    self.__index = self

    return obj
end

function TcpClient:isConnected()
    return self.state == states.connected
end

function TcpClient:_callback(event, ...)
    if self.callbacks[event] == nil then
        return;
    end
    for idx = 1, #self.callbacks[event] do
        self.callbacks[event][idx](...)
    end
end

function TcpClient:_get_ip(host)
    local results = vim.loop.getaddrinfo(host)
    local actual_addr = nil

    for idx = 1, #results do
        local res = results[idx]
        if res.family == "inet" and res.socktype == "stream" then
            actual_addr = res.addr
        end
    end

    return actual_addr
end

function TcpClient:_connect_to_server()
    self.client = vim.loop.new_tcp()

    local ip = self:_get_ip(self.host)

    self.state = states.connecting

    self.client:connect(ip, tonumber(self.port), function (err)
        if self.state ~= states.connecting then
            return
        end

        if err ~= nil then
            self.state = states.error
            self:_callback("connect", err)
            return
        end

        self.state = states.connected

        self.client:read_start(vim.schedule_wrap(function(_, chunk)
            if chunk == nil then
                self:_callback("disconnected", nil)
                return
            end

            local bytes = { string.byte(chunk, 1,-1) }
            self:_callback("data", bytes)
        end))

        self:_callback("connect", nil);
    end)

    vim.fn.timer_start(10000, function()
        if self.state == states.connected then
            return
        end

        self.state = states.timedout
        self:_callback("connect",
            string.format("Unable to connect to %s:%d", self.host, self.port))
    end)
end

function TcpClient:isError()
    return self.state == states.error
end

function TcpClient:isTimeout()
    return self.state == states.timedout
end

function TcpClient:disconnect()
    if self.state == states.connecting then
        self.state = states.disconnected
    end

    if self.client == nil then
        return
    end

    self.client:shutdown()
    self.client:close()
    self.client = nil
end

function TcpClient:connect()
    if self.state ~= states.disconnected then
        return
    end

    self:_connect_to_server()
end

function TcpClient:on(event, callback)
    if self.callbacks[event] == nil then
        self.callbacks[event] = {}
    end

    table.insert(self.callbacks[event], callback);
end

return TcpClient
