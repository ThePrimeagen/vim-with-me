local DATA = "./data/partial"
local plenary = require("plenary.reload")
plenary.reload_module("vim-with-me")

local App = require("vim-with-me.app")
local ColorSet = require("vim-with-me.app.colors")
local parse = require("vim-with-me.tcp.parse")
local window = require("vim-with-me.window")
local TestUtils = require("vim-with-me.test-utils")

local function manual_run()
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

    ---[[
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
end

local function run_app()
    local app = App:new(TestUtils.fake_tcp_from_file(DATA))
end

run_app()

