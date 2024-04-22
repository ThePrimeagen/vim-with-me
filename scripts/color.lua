local plenary = require("plenary.reload")
plenary.reload_module("vim-with-me")

local ColorSet = require("vim-with-me.app.colors")
local parse = require("vim-with-me.tcp.parse")
local process = require("vim-with-me.tcp.process")
local window = require("vim-with-me.window")

local processor = process.process_packets()

local ok, fh = pcall(vim.loop.fs_open, "/tmp/partials", "r", 493)
if not ok then
    error("cannot open file")
end

local ok, data = pcall(vim.loop.fs_read, fh, 2048)
if not ok then
    error("cannot read data")
end

vim.loop.fs_close(fh)

local parsed_partials = parse.parse_partial_renders(processor(data).data)
local win = window.create_window({
    width = 5,
    height = 5,
    row = 0,
    col = 0,
}, true)
vim.api.nvim_buf_set_lines(win.buffer, 0, -1, false, {
    "XXXXX",
    "XXXXX",
    "XXXXX",
    "XXXXX",
    "XXXXX",
})
local colorer = ColorSet:new(win)

--- @param partials CellWithLocation[]
--- @param i number
local function color(partials, i)
    if i > #partials then
        print("leaving", i, #partials)
        return
    end

    colorer:color_cell(partials[i])
    vim.defer_fn(function()
        color(partials, i + 1)
    end, 50)
end

color(parsed_partials, 1)
