local eq = assert.are.same
local Commands = require("vim-with-me.app.commands")

describe("app#commands", function()
    it("should decode the commands and return their values", function()
        local decode = {
            open = 8,
            close = 9,
            foo = 10,
        }

        local str = ""
        for k, v in pairs(decode) do
            str = str .. k .. "\n" .. string.char(v)
        end

        local cmd = Commands.Commands:new()
        cmd:parse(str)

        eq(8, cmd:get("open"))
        eq(9, cmd:get("close"))
        eq(10, cmd:get("foo"))
        eq(Commands.DefinedCommands.MISSING, cmd:get("bar"))
    end)
end)





