---@alias TCPCommandName "render" | "partial" | "close" | "error" | "openWindow" | "commands" | string

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

---@param name TCPCommandName
---@return number
function Commands:get(name)
    local value = self._commands[name]
    assert(value ~= nil, "command not found " .. name)
    return value
end

---@param cmd TCPCommand
---@param name TCPCommandName
---@return boolean
function Commands:is(cmd, name)
    local value = self._commands[name]
    assert(value ~= nil, "command not found " .. name)
    return cmd.command == value
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

---@param cmd TCPCommand
function Commands:pretty_print(cmd)
    local name = ""
    for k, v in pairs(self._commands) do
        if v == cmd.command then
            name = k
            break
        end
    end

    print("cmd", name, #cmd.data)
end

return {
    Commands = Commands,
}
