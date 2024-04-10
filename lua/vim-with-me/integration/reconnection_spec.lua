local eq = assert.are.same
local utils = require("vim-with-me.tcp.utils")
local int_utils = require("vim-with-me.integration.int_utils")
local Commands = require("vim-with-me.app.commands")
local theprimeagen = int_utils.theprimeagen

local PORT = 42073

---@return TestTCP
local function create_tcp()
    local tcp = int_utils.create_tcp_connection(PORT)
    local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)
    return {
        tcp = tcp,
        next = next_cmd,
        flush = flush_cmds
    }
end

describe("vim with me :: reconnecting test", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("multiuser test", function()
        int_utils.create_test_server("connection_server", PORT)

        for _ = 1, 5 do
            print("create connection")
            local tcp = create_tcp()
            vim.wait(500)
            print("close connection")
            tcp.tcp:close()
            vim.wait(500)
        end

        --[[

        local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)

        local server_commands = next_cmd()
        assert(server_commands ~= nil)
        local commands = Commands.Commands:new()
        commands:parse(server_commands.data)

        tcp:send({ command = commands:get("open"), data = "" })
        eq({
            command = commands:get("openWindow"),
            data = utils.to_string(24, 80),
        }, next_cmd())

        tcp:send({ command = commands:get("render"), data = "" })
        local cmd = next_cmd()

        eq({
            command = commands:get("render"),
            data = theprimeagen,
        }, cmd)
        tcp:send({
            command = commands:get("partial"),
            data = utils.to_string(1, 1),
        })
        local theprimeagen_str = "theprimeagen"
        local data = ""
        for i = 1, #theprimeagen_str do
            data = data
                .. string.format(
                    utils.to_string(1, i, theprimeagen_str:sub(i, i))
                )
        end

        eq({
            { command = commands:get("partial"), data = data },
        }, flush_cmds())
        --]]
    end)
end)


