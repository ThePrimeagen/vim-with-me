local eq = assert.are.same
local parse = require("vim-with-me.tcp.parse")
local utils = require("vim-with-me.tcp.utils")

describe("vim with me :: tcp.parse", function()
    it("parse out location", function()
        eq({
            row = 69,
            col = 42,
        }, parse.parse_location(utils.to_string(69, 42)))
    end)

    it("parse out color", function()
        eq({
            foreground = true,
            red = 69,
            green = 42,
            blue = 128,
        }, parse.parse_color(utils.to_string(1, 69, 42, 128)))

        eq({
            foreground = false,
            red = 69,
            green = 42,
            blue = 128,
        }, parse.parse_color(utils.to_string(0, 69, 42, 128)))
    end)

    it("parse out cell", function()
        eq(
            {
                foreground = {
                    foreground = true,
                    red = 69,
                    green = 42,
                    blue = 128,
                },
                background = {
                    foreground = false,
                    red = 68,
                    green = 41,
                    blue = 127,
                },
                value = "E",
            },
            parse.parse_cell(
                utils.to_string(
                    "E",
                    utils.to_string(1, 69, 42, 128),
                    utils.to_string(0, 68, 41, 127)
                )
            )
        )
    end)

    it("parse out cell with location", function()
        eq(
            {
                loc = {
                    row = 13,
                    col = 37,
                },
                cell = {
                    foreground = {
                        foreground = true,
                        red = 69,
                        green = 42,
                        blue = 128,
                    },
                    background = {
                        foreground = false,
                        red = 68,
                        green = 41,
                        blue = 127,
                    },
                    value = "E",
                },
            },
            parse.parse_cell_with_location(
                utils.to_string(
                    utils.to_string(13, 37),
                    "E",
                    utils.to_string(1, 69, 42, 128),
                    utils.to_string(0, 68, 41, 127)
                )
            )
        )
    end)
end)
