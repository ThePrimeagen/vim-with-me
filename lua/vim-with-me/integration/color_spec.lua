local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")
local parse = require("vim-with-me.tcp.parse")
local utils = require("vim-with-me.tcp.utils")

local port = 42075

describe("vim with me :: Color test", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("do we color?", function()
        int_utils.create_test_server("color_server", port, { debug = true })
        local tcp = int_utils.create_tcp_connection(port)

        local next = int_utils.create_tcp_next(tcp)
        next()
        next()

        local x = next()

        eq({
            cell = {
                background = {
                    blue = 69,
                    foreground = false,
                    green = 0,
                    red = 255,
                },
                foreground = {
                    blue = 42,
                    foreground = true,
                    green = 255,
                    red = 69,
                },
                value = "X",
            },
            loc = {
                col = 9,
                row = 6,
            },
        }, parse.parse_cell_with_location(x.data))
    end)
end)
