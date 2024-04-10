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

---@param cmds TCPCommand[]
---@param expected TCPCommand
local function assert_cmds(cmds, expected)
    for _, v in ipairs(cmds) do
        eq(expected, v)
    end
end

describe("vim with me :: reconnecting test", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("reconnection test", function()
        int_utils.create_test_server("connection_server", PORT)

        local tcps = {
            create_tcp(),
            create_tcp(),
            create_tcp(),
        }

        -- read app commands
        int_utils.read_all(tcps)

        for i = 1, 3 do
            local expected = { command = 69, data = "your mom: " .. i }
            tcps[1].tcp:send(expected)

            local cmds = int_utils.read_all(tcps)
            assert_cmds(cmds, expected)

            tcps[1].tcp:close()
            vim.wait(100)

            table.remove(tcps, 1)

            local tcp = create_tcp()
            table.insert(tcps, tcp)
            tcp.next()
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


