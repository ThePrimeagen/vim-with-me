--- TERRIBLE NAME
-- local TD = require("vim-with-me.td")
local window = require("vim-with-me.window")
local tcp = require("vim-with-me.tcp")
local DisplayCache = require("vim-with-me.window.cache")

---@param w WindowDetails
local function on_read(w)
    local cache = DisplayCache.new(80, 24)
    local count = 0
    local start = vim.loop.now()

    local previous_chunk = ""

    ---@param chunk string
    return function(chunk)
        -- split by :
        chunk = previous_chunk .. chunk
        previous_chunk = ""

        while #chunk > 0 do
            local idx = string.find(chunk, ":")

            -- what about keeping previous chunks?
            if idx == nil then
                previous_chunk = chunk
                break
            end

            local len = tonumber(string.sub(chunk, 1, idx - 1))
            if len + idx > #chunk then
                previous_chunk = chunk
                break
            end

            local next_chunk, command, data = tcp.parse(chunk, idx + 1, len)

            if command == "t" then
                -- aint no one is placing towers right now
            elseif command == "r" then
                count = count + 1
                if count % 500 == 0 then
                    print("count: ", count, vim.loop.now() - start)
                    start = vim.loop.now()
                end
                -- check to see if last character is a new line
                if string.sub(data, -1) == "\n" then
                    data = string.sub(data, 1, -2)
                end
                cache:from_string(data)

                local rows = cache:to_string_rows()
                vim.api.nvim_buf_set_lines(w.buffer, 0, -1, false, rows)
            end

            chunk = next_chunk
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




