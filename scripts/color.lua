local DATA = "./data/partial"
local plenary = require("plenary.reload")
plenary.reload_module("vim-with-me")

local App = require("vim-with-me.app")
local ColorSet = require("vim-with-me.app.colors")
local parse = require("vim-with-me.tcp.parse")
local process = require("vim-with-me.tcp.process")
local window = require("vim-with-me.window")

local processor = process.process_packets()
local FakeTCP = { }
FakeTCP.__index = FakeTCP

function FakeTCP:new(tcp_data)
    local item = setmetatable({
        process = processor,
        data = tcp_data
    }, self)

    return item
end

function FakeTCP:connected()
    return true
end

function FakeTCP:listen(cb)
    local packet = processor(self.data)
    local function read_one()
        if packet == nil then
            return
        end

        cb(packet)
        packet = processor()

        vim.defer_fn(function()
            read_one()
        end, 500)
    end

    read_one()
end

local ok, fh = pcall(vim.loop.fs_open, DATA, "r", 493)
if not ok then
    error("cannot open file")
end

local ok, data = pcall(vim.loop.fs_read, fh, 2048)
if not ok then
    error("cannot read data")
end

vim.loop.fs_close(fh)

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
    local app = App:new(FakeTCP:new(data))
end

run_app()

