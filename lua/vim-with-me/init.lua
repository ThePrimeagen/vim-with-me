--- TERRIBLE NAME
-- local TD = require("vim-with-me.td")
local window = require("vim-with-me.window")
local tcp = require("vim-with-me.tcp")
local DisplayCache = require("vim-with-me.window.cache")

---@param w WindowDetails
local function on_read(w)
    local cache = DisplayCache.new(80, 24)

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
            vim.api.nvim_buf_set_lines(w.buffer, 0, -1, false, rows)
        end
    end
end

function START()
    assert(not tcp.tcp_connected(), "client already started")

    tcp.tcp_start()
    local w = window.create_window()

    tcp.listen(on_read(w))
end

function CLOSE()
    assert(tcp.tcp_connected(), "client must be connected")
    tcp.tcp_stop()
end




