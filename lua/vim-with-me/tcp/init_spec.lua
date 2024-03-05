-- luacheck: globals describe it assert
local eq = assert.are.same
local tcp = require("vim-with-me.tcp")

describe("vim with me :: tcp", function()
    it("parse", function()
        local str = "5:r:hel"
        local next_chunk, command, data = tcp.parse(str, 3, 5)

        eq(next_chunk, "")
        eq(command, "r")
        eq(data, "hel")
    end)
end)

