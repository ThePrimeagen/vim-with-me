local EnvClient = require("vim-with-me.envelope");
local TCP = require("vim-with-me.tcp");
local Enum = require("vim-with-me.enum");
local Pobo = require("vim-with-me.pobo");

local CommandTypes = Enum({
    VimCommand = 0,
    ASDF = 1,
    Xrandr = 2,
    StatusLineUpdate = 3,
    GiveawayEnter = 4,
    VimInsert = 5,
    VimAfter = 6,
})

-- local tcp = TCP:new("vwm.theprimeagen.tv", 42069)
local tcp = TCP:new("localhost", 42069)
local env = EnvClient:new(tcp);

env:connect()
env:on("connect", function()
    print("Connected to the server")
end)

giveaway_entries = {}

local function tablelength(T)
    local count = 0
    for _ in pairs(T) do count = count + 1 end
    return count
end

function reset_giveaway()
    giveaway_entries = {}
end

function pick_an_entry()
    --local length = tablelength(giveaway_entries)
    local pick = math.ceil(math.random() * #giveaway_entries)
    return giveaway_entries[pick]

    --[[
    local count = 0
    for k in pairs(giveaway_entries) do
        if pick == count then
            return k
        end
        count = count + 1
    end
    ]]

    -- return "I somehow didn't find one.  I am very sorry.  Give it to the person you love the most";
end

env:on("data", function(line)
    local pobo = Pobo:new(line, 1)
    local status = pobo:get_status()

    print("STATUS", status);
    require("theprimeagen.statusline").set_status(status)

    if line[1] == CommandTypes.VimCommand or
       line[1] == CommandTypes.VimAfter or
       line[1] == CommandTypes.VimInsert then
        local cmd = pobo:get_data()
        vim.cmd(cmd)
    elseif line[1] == CommandTypes.GiveawayEnter then
        local name = pobo:get_data()
        table.insert(giveaway_entries, name)
    end
end)

