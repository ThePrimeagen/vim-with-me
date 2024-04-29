local eq = assert.are.same
local Colors = require("vim-with-me.app.colors")

describe("app#color", function()
    it("just try to color things without erroring", function()
        local color = Colors:new({
            dim = {
                width = 80,
                height = 24,
                row = 0,
                col = 0,
            },
            buffer = 0,
            win_id = 0,
        })

        color:color_cell({
            loc = {
                row = 13,
                col = 69,
            },
            cell = {
                foreground = {
                    red = 12,
                    green = 34,
                    blue = 56,
                    foreground = true,
                },
                background = {
                    red = 65,
                    green = 43,
                    blue = 21,
                    foreground = false,
                },
                value = "E",
            },
        })

        -- TODO: Teej help, how do i check the color of a namespace at a particular place?
    end)
end)
