-- luacheck: globals describe it assert
local eq = assert.are.same
local DisplayCache = require("vim-with-me.window.cache")
local window = require("vim-with-me.window")

describe("vim with me :: cache", function()
    it("i hate chat", function()
        local cache = DisplayCache:new(window.create_window_dimensions(10, 3))

        cache:from_string("1234567890abcdefghijXXXXXYYYYY")
        eq({
            "1234567890",
            "abcdefghij",
            "XXXXXYYYYY",
        }, cache:to_string_rows())
    end)
end)
