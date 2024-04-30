local ZERO = string.byte("$")
local BASE = 90

local function to_string(...)
    local out = {}
    for i = 1, select("#", ...) do
        local v = select(i, ...)
        if type(v) == "number" then
            table.insert(out, string.char(v))
        elseif type(v) == "string" then
            table.insert(out, v)
        else
            assert(
                false,
                "should never provide anything other than numbers, strings, or tables of strings and numbers"
            )
        end
    end
    return table.concat(out, "")
end

---@param cmd TCPCommand | nil
---@return TCPCommand | nil
local function pretty_print(cmd)
    if cmd == nil then
        print("command is nil")
        return cmd
    end

    print("command", cmd.command)

    local to_print = {}
    for i = 1, #cmd.data do
        table.insert(to_print, string.byte(cmd.data, i, i))
    end

    print(table.concat(to_print, ", "))
    return cmd
end

return {
    to_string = to_string,
    pretty_print = pretty_print,
}
