local Enum = require("vim-with-me.enum");

local EnvelopeClient = {}
local States = Enum({
    WaitingForHeader = "WaitingForHeader",
    ParsingHeader = "ParsingHeader",
    ParsingBody = "ParsingBody",
})

function EnvelopeClient:new(tcpClient)

    local obj = {
        _current_length = 0,
        _current_env = nil,
        _client = tcpClient,
        _state = States.WaitingForHeader,
        callbacks = {},
    }

    setmetatable(obj, self)
    self.__index = self

    return obj
end

function EnvelopeClient:connect()
    self._client:on("data", function(data)
        self:_parse(data)
    end)

    self._client:on("connect", function()
        self:_callback("connect")
    end)

    self._client:connect()
end

local function merge(t1, t2, t2_offset, t2_max)
    -- I am offended by how much waste
    local out = {}

    if t1 then
        for idx = 1, #t1 do
            table.insert(out, t1[idx])
        end
    end

    if t2 then
        t2_max = t2_offset + (t2_max or #t2) - 1
        for idx = t2_offset, t2_max do
            table.insert(out, t2[idx])
        end
    end

    return out
end

function EnvelopeClient:_callback(event, ...)

    if self.callbacks[event] == nil then
        return;
    end
    for idx = 1, #self.callbacks[event] do
        self.callbacks[event][idx](...)
    end
end

function EnvelopeClient:_reset()
    self._current_length = 0
    self._current_env = nil
    self._state = States.WaitingForHeader
end

local function readUint16BE(bytes, offset)
    local high, low

    if type(bytes) == "table" then
        high = bytes[offset]
        low = bytes[offset + 1]
    else
        high = bytes
        low = offset
    end

    return bit.lshift(high, 8) + low
end

function EnvelopeClient:_parse(bytes)
    local idx = 1

    repeat
        print("EnvelopeClient:_parse", self._state, idx)
        local bytes_remaining = #bytes - (idx - 1)

        if self._state == States.WaitingForHeader then

            if bytes_remaining < 2 then
                self._current_length = bytes[idx]
                idx = idx + 1
                self._state = States.ParsingHeader
            else
                self._state = States.ParsingBody
                self._current_length = readUint16BE(bytes, idx)
                idx = idx + 2
            end

        elseif self._state == States.ParsingHeader then
            self._current_length = readUint16BE(self._current_length, bytes[idx])
            idx = idx + 1
            self._state = States.ParsingBody

        elseif self._state == States.ParsingBody then
            local env_remaining = self._current_length
            if self._current_env then
                env_remaining = env_remaining - #self._current_env
            end

            local env = merge(self._current_env, bytes, idx, env_remaining)

            if bytes_remaining >= env_remaining then
                self:_callback("data", env)
                self:_reset()
            else
                self._current_env = env
            end

            idx = idx + env_remaining
        end

    until idx > #bytes
end

function EnvelopeClient:on(key, cb)
    if self.callbacks[key] == nil then
        self.callbacks[key] = {}
    end

    table.insert(self.callbacks[key], cb)
end

EnvelopeClient.States = States

return EnvelopeClient
