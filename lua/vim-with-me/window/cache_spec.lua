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

    it("the above unit test is aptly named", function()
        local w = window.create_window(window.create_window_dimensions(3, 3))
        local cache = DisplayCache:new(w.dim)

        cache:from_string("XXXXXXXXX")
        eq({
            "XXX",
            "XXX",
            "XXX",
        }, cache:to_string_rows())
        vim.api.nvim_buf_set_lines(
            w.buffer,
            0,
            -1,
            false,
            cache:to_string_rows()
        )

        local expect = {
            { "YXX", "XXX", "XXX" },
            { "YYX", "XXX", "XXX" },
            { "YYY", "XXX", "XXX" },
            { "YYY", "YXX", "XXX" },
            { "YYY", "YYX", "XXX" },
            { "YYY", "YYY", "XXX" },
            { "YYY", "YYY", "YXX" },
            { "YYY", "YYY", "YYX" },
            { "YYY", "YYY", "YYY" },
        }

        local updates = {
            { 1, 1 },
            { 1, 2 },
            { 1, 3 },
            { 2, 1 },
            { 2, 2 },
            { 2, 3 },
            { 3, 1 },
            { 3, 2 },
            { 3, 3 },
        }

        for i, item in pairs(updates) do
            cache:place(item[1], item[2], "Y")
            cache:render_into(w)
            eq(expect[i], vim.api.nvim_buf_get_lines(w.buffer, 0, -1, false))
        end
    end)
end)
