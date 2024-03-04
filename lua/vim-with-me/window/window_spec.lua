-- luacheck: globals describe it assert
local eq = assert.are.same
local window = require("vim-with-me.window")

describe("vim with me :: window", function()
    it("should be able to open and close a float", function()
        local details = window.create_window()

        eq(true, vim.api.nvim_win_is_valid(details.win_id))
        eq(true, vim.api.nvim_buf_is_valid(details.buffer))

        window.close_window(details)

        eq(false, vim.api.nvim_win_is_valid(details.win_id))
        eq(false, vim.api.nvim_buf_is_valid(details.buffer))
    end)
end)
