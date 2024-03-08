-- luacheck: globals describe it assert
local eq = assert.are.same
local DisplayCache = require("vim-with-me.window.cache")

describe("vim with me :: cache", function()
    it("i hate chat", function()
        local cache = DisplayCache:new(10, 2)

        cache:from_string("1234567890abcdefghij")
        eq({ "1234567890", "abcdefghij" }, cache:to_string_rows())
    end)
end)
