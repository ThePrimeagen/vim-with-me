local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")

local port = 42075

describe("vim with me :: Color test", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("do we color?", function()
        int_utils.create_test_server("color_server", port)
        local tcp = int_utils.create_tcp_connection(port)

        local next = int_utils.create_tcp_next(tcp)
        local commands = next()
        local open_window = next()
        local x = next()

        print(vim.inspect(commands))
        print(vim.inspect(open_window))
        print(vim.inspect(x))
    end)
end)

