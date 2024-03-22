local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")
local theprimeagen = int_utils.theprimeagen

describe("vim with me", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("full command set", function()
        local tcp = int_utils.create_test_conn("cmd_server", 42070)
        local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)

        tcp:send("open", "")
        eq({
            command = "open-window",
            data = "24:80",
        }, next_cmd())

        tcp:send("render", "")
        local cmd = next_cmd()

        eq({
            command = "r",
            data = theprimeagen,
        }, cmd)
        tcp:send("partial", "1:1")

        local theprimeagen_str = "theprimeagen"
        local data = ""
        for i = 1, #theprimeagen_str do
            data = data .. string.format("1:%d:%s", i, theprimeagen_str:sub(i, i))
        end

        eq({
            {command = "p", data = data}
        }, flush_cmds())
    end)
end)

