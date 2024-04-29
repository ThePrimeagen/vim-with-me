local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")

local PORT = 42073

---@return TestTCP
local function create_tcp()
    local tcp = int_utils.create_tcp_connection(PORT)
    local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)
    return {
        tcp = tcp,
        next = next_cmd,
        flush = flush_cmds,
    }
end

---@param cmds TCPCommand[]
---@param expected TCPCommand
local function assert_cmds(cmds, expected)
    assert(#cmds > 0, "there must be commands to assert on")
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
    end)
end)
