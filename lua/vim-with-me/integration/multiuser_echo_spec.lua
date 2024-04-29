local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")

---@alias TestTCP {tcp: TCP, next: (fun(): TCPCommand | nil), flush: fun(): TCPCommand[]}

local PORT = 42072

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

describe("vim with me :: multiuser", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("multiuser test", function()
        int_utils.create_test_server("echo_server", PORT)
        local tcps = {
            create_tcp(),
            create_tcp(),
            create_tcp(),
        }

        -- discard the command replies
        int_utils.read_all(tcps)

        for idx = 1, 3 do
            local expected = { command = idx, data = "hello" .. idx }

            tcps[1].tcp:send(expected)

            local outs = int_utils.read_all(tcps)

            for _, v in ipairs(outs) do
                eq(expected, v)
            end

            tcps[1].tcp:close()
            table.remove(tcps, 1)

            vim.wait(100)

            tcps[1].tcp:send(expected)
            outs = int_utils.read_all(tcps)

            for _, v in ipairs(outs) do
                eq(expected, v)
            end

            local tcp = create_tcp()
            table.insert(tcps, tcp)

            -- read command
            tcp.next()
        end

        for idx = 1, 3 do
            local expected = { command = idx, data = "hello" .. idx }

            tcps[3].tcp:send(expected)

            local outs = int_utils.read_all(tcps)

            for _, v in ipairs(outs) do
                eq(expected, v)
            end

            tcps[3].tcp:close()
            table.remove(tcps, 3)

            tcps[1].tcp:send(expected)
            outs = int_utils.read_all(tcps)

            for _, v in ipairs(outs) do
                eq(expected, v)
            end

            local tcp = create_tcp()
            table.insert(tcps, 1, tcp)

            -- read command
            tcp.next()
        end
    end)
end)
