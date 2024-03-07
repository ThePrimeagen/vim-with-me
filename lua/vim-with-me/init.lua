--- TERRIBLE NAME
-- local TD = require("vim-with-me.td")
local window = require("vim-with-me.window")
local TCP = require("vim-with-me.tcp")
local DisplayCache = require("vim-with-me.window.cache")

---@type WindowDetails | nil
local w = nil

---@type TCP | nil
local conn = nil

---@type DisplayCache | nil
local cache = nil

---@param win WindowDetails
---@param cache DisplayCache
local function on_read(win, cache)

    ---@param command string
    ---@param data string
    return function(command, data)
        if command == "t" then
            -- aint no one is placing towers right now
        elseif command == "r" then
            -- check to see if last character is a new line
            if string.sub(data, -1) == "\n" then
                data = string.sub(data, 1, -2)
            end
            cache:from_string(data)

            local rows = cache:to_string_rows()
            vim.api.nvim_buf_set_lines(win.buffer, 0, -1, false, rows)
        end
    end
end

function START()
    assert(conn == nil, "client already started")
    assert(w == nil, "window already created")
    assert(cache == nil, "cache already created")

    conn = TCP:new()
    conn:start(function()
        print("connected")
    end)

    w = window.create_window(
        window.create_window_dimensions(80, 24),
        true
    )
    cache = DisplayCache:new(80, 24)
    conn:listen(on_read(w, cache))
end

function CLOSE()
    assert(conn ~= nil, "client not started")
    assert(w ~= nil, "window not created")

    conn:close()
    window.close_window(w)

    conn = nil
    w = nil
    cache = nil
end

function SEND()
    assert(conn ~= nil, "client not started")
    conn:send("render", "")
end
