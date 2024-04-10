local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")
local utils = require("vim-with-me.tcp.utils")
local App = require("vim-with-me.app")

describe("vim with me :: app_spec", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("app integration test", function()
        int_utils.create_test_server("cmd_server", 42071)
        local tcp = int_utils.create_tcp_connection(42071)
        local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)
        local count = 0
        local app = App:new(tcp):on_render(function()
            count = count + 1
        end)

        next_cmd()

        tcp:send({ command = app.commands:get("open"), data = "" })
        next_cmd()

        eq(false, app.window == nil)
        eq(app.window.dim.width, 80)
        eq(app.window.dim.height, 24)
        eq(int_utils.empty, table.concat(app.cache:to_string_rows()))

        tcp:send({ command = app.commands:get("render"), data = "" })
        flush_cmds()
        eq(int_utils.theprimeagen, table.concat(app.cache:to_string_rows()))

        tcp:send({
            command = app.commands:get("partial"),
            data = utils.to_string(1, 1),
        })
        local cmds = flush_cmds()
        eq(#cmds > 0, true)
        eq(
            int_utils.theprimeagen_partial,
            table.concat(app.cache:to_string_rows())
        )
    end)
end)
