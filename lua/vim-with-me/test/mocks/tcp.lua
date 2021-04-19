local TcpMocks = {}

function TcpMocks:new()
    local obj = {
        callbacks = {},
    }

    setmetatable(obj, self)
    self.__index = self

    return obj
end

function TcpMocks:on(event, callback)
    if self.callbacks[event] == nil then
        self.callbacks[event] = {}
    end

    table.insert(self.callbacks[event], callback);
end

function TcpMocks:connect()
    self:emit("connect")
end

function TcpMocks:emit(event, ...)

    if self.callbacks[event] == nil then
        return;
    end

    for idx = 1, #self.callbacks[event] do
        self.callbacks[event][idx](...)
    end
end

return TcpMocks
