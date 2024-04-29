local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")
local utils = require("vim-with-me.tcp.utils")
local optsify = require("vim-with-me.utils").optsify
local App = require("vim-with-me.app")

---@param app VWMApp
---@param received_cmds TCPCommand[]
---@param command TCPCommandName
---@param opts {timeout: number?, debug: boolean?}?
local function wait_for(app, received_cmds, command, opts)
    opts = optsify(opts, { timeout = 1000, debug = false })
    local found = false
    vim.wait(opts.timeout, function()
        for i, cmd in ipairs(received_cmds) do
            if app.commands:is(cmd, command) then
                if opts.debug then
                    utils.pretty_print(cmd)
                end
                table.remove(received_cmds, i)
                found = true
                return true
            end
        end

        return false
    end)
    if opts.debug and not found then
        print("wait_for: timer expired", command)
    end
end

describe("vim with me :: app_spec", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("app integration test", function()
        int_utils.create_test_server("cmd_server", 42071)
        local tcp = int_utils.create_tcp_connection(42071)
        local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)
        local count = 0
        local received_cmds = {}
        local app = App:new(tcp)
            :on_render(function()
                count = count + 1
            end)
            :on_command(function(cmd)
                table.insert(received_cmds, cmd)
            end)

        wait_for(app, received_cmds, "commands")

        tcp:send({ command = app.commands:get("open"), data = "" })
        wait_for(app, received_cmds, "openWindow")

        eq(false, app.window == nil)
        eq(80, app.window.dim.width)
        eq(24, app.window.dim.height)
        eq(int_utils.empty, table.concat(app.cache:to_string_rows()))

        tcp:send({ command = app.commands:get("render"), data = "" })
        wait_for(app, received_cmds, "partial")

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
