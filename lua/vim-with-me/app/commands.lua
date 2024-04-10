---@class TCPCommands
---@field _commands table<string, number>
local Commands = {}
Commands.__index = Commands

---@return TCPCommands
function Commands:new()
    return setmetatable({
        _commands = { commands = 0 },
    }, self)
end

---@param name "render" | "partial" | "close" | "error" | "openWindow" | "commands" | string
---@return number
function Commands:get(name)
    local value = self._commands[name]
    assert(value ~= nil, "command not found " .. name)
    return value
end

---@param str string
function Commands:parse(str)
    local idx = 1
    while #str > idx do
        local newline_idx = string.find(str, "\n", idx)
        assert(newline_idx ~= nil, "newline_idx should never be nil")

        local command_name = str.sub(str, idx, newline_idx - 1)
        local command = str.byte(str, newline_idx + 1)

        self._commands[command_name] = command

        idx = newline_idx + 2
    end
end

return {
    Commands = Commands,
}
