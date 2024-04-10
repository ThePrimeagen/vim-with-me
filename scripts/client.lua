-- assumes that the cwd is vim-with-me
vim.opt.rtp:append(vim.loop.cwd())

local int_utils = require("vim-with-me.integration.int_utils")
local tcp = int_utils.create_test_server("cmd_server", 42070)
local next_cmd, flush_cmds = int_utils.create_tcp_next(tcp)

function Send(cmd, data)
    tcp:send(cmd, data)
end

function Next()
    return next_cmd()
end

function Open()
    tcp:send("open", "")
end

function Render()
    tcp:send("render", "")
end

function PartialRender(row, col)
    assert(type(row) == "number", "row must be number")
    assert(type(col) == "number", "col must be number")


end

