local EnvClient = require("vim-with-me.envelope");
local TCP = require("vim-with-me.tcp");
local Enum = require("vim-with-me.enum");
local Pobo = require("vim-with-me.pobo");

local CommandTypes = Enum({
    VimCommand = 0,
    SystemCommand = 1,
    StatusLineUpdate = 3,
    GiveawayEnter = 4,
    VimInsert = 5,
    VimAfter = 6,
})

local M = {
    _giveaway_entries = {},
}

function M.reset_giveaway()
    M._giveaway_entries = {}
end

function M.pick_an_entry()
    --local length = tablelength(giveaway_entries)
    local pick = math.ceil(math.random() * #M._giveaway_entries)
    return M._giveaway_entries[pick]
end

function M.disconnect()
    if not M._tcp then
        return
    end

    M._tcp:disconnect()
    M._init = false
    M._tcp = nil
    M._env = nil
end

function M.init()
    if M._init then
        return
    end
    M._init = true

    local vwm_host = vim.env.VWM_HOST or "localhost"
    M._tcp = TCP:new(vwm_host, 42069)
    M._env = EnvClient:new(M._tcp);
    M._env:connect()
    M._env:on("connect", function()
        print("Connected to the server")
    end)

    M._env:on("data", function(line)
        local pobo = Pobo:new(line, 1)
        -- TODO: MAKOE THE STATUS LINE AGAIN YA GOOF
        -- local status = pobo:get_status()

        -- require("theprimeagen.statusline").set_status(status)

        if line[1] == CommandTypes.VimCommand or
            line[1] == CommandTypes.VimAfter or
            line[1] == CommandTypes.VimInsert then
            local cmd = pobo:get_data()
            vim.cmd(cmd)
        elseif line[1] == CommandTypes.SystemCommand then
            vim.cmd(string.format("silent! !%s", pobo:get_data()))
        elseif line[1] == CommandTypes.GiveawayEnter then
            local name = pobo:get_data()
            table.insert(M._giveaway_entries, name)
        end
    end)
end

return M
