local EnvClient = require("vim-with-me.envelope");
local TCP = require("vim-with-me.tcp");

local tcp = TCP:new("localhost", 42069)
local env = EnvClient:new(tcp);

env:connect()

local Enum = require("vim-with-me.enum");

local CommandTypes = Enum({
    StatusLineUpdate = 3,
})

local function slice(t, start, stop)
    local out = {}
    for idx = start, stop do
        table.insert(out, t[idx])
    end
    return out
end

-- stackoverflow, though dangerous behind the wheel can still serve a purpose
local function utf8_from(t, start, stop)
    local bytearr = {}
    local idx = start
    while idx <= stop do
        local v = t[idx]
        local utf8byte = v < 0 and (0xff + v + 1) or v
        table.insert(bytearr, string.char(utf8byte))
        idx = idx + 1
    end
    return table.concat(bytearr)
end

env:on("data", function(line)
    print("Line", vim.inspect(line))
    if line[1] == CommandTypes.StatusLineUpdate then
        local length = 1
        while line[1 + length] ~= 0 and length < 50 do
            length = length + 1
        end

        local status = utf8_from(line, 1, length)
        -- TODO: why... for future prime.. why the sub 2?
        require("theprimeagen.statusline").set_status(status:sub(2))
    end
end)

