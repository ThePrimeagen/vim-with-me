local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")
local parse = require("vim-with-me.tcp.parse")
local utils = require("vim-with-me.tcp.utils")

local port = 42075

describe("vim with me :: Color test", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("do we color?", function()
        int_utils.create_test_server("color_server", port)
        local tcp = int_utils.create_tcp_connection(port)

        local next = int_utils.create_tcp_next(tcp)
        next()
        next()
        local x = next()

        utils.pretty_print(x)

        eq({
            cell = {
                background = {
                    blue = 0,
                    foreground = false,
                    green = 0,
                    red = 0
                },
                foreground = {
                    blue = 69,
                    foreground = true,
                    green = 169,
                    red = 100
                },
                value = "X"
            },
            loc = {
                col = 8,
                row = 8,
            }
        }, parse.parse_cell_with_location(x.data))

    end)
end)

