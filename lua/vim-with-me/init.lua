--- TERRIBLE NAME
local TD = require("vim-with-me.td")

-- luacheck: ignore 111
local uv = vim.loop

local function key(k)
    k = vim.api.nvim_replace_termcodes(k, true, false, true)
    vim.api.nvim_feedkeys(k, "t", false)
end

---@param client any
---@param td TowerOffense
local function read(client, td)
    uv.read_start(
        client,
        vim.schedule_wrap(function(_, chunk)
            --if chunk == "<dot>" then
            --    chunk = "."
            --end
            --key(chunk)

            -- split by :
            local parts = {}
            for part in string.gmatch(chunk, "([^:]+)") do
                table.insert(parts, part)
            end

            if parts[1] ~= "t" then
                error("unknown command")
            end

            local x = tonumber(parts[2])
            local y = tonumber(parts[3])

            if y == 0 then
                print("chat sucks or your protocol sucks")
                return
            end

            td:place(x, y)
        end)
    )
end

local client = nil
local td = nil

function START()
    assert(client == nil, "client already started")
    assert(td == nil, "td already started")

    td = TD:new()
    td:start()

    print(vim.inspect(td._window_details))

    client = uv.new_tcp()
    client:connect("127.0.0.1", 42069, function(_)
        read(client, td)
    end)

end

function CLOSE()
    assert(client ~= nil, "client hasn't been started")
    assert(td ~= nil, "td hasn't been started")

    client:shutdown()
    client:close()

    td:close()

    client = nil
    td = nil
end
