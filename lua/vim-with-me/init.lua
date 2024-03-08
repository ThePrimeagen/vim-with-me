--- TERRIBLE NAME
-- local TD = require("vim-with-me.td")
local window = require("vim-with-me.window")
local TCP = require("vim-with-me.tcp").TCP
local App = require("vim-with-me.app")

---@type VWMApp | nil
local app = nil

function START()
    assert(app == nil, "client already started")

    local conn = TCP:new()
    conn:start(function()
        local function handle_commands(cmd, data)
            print("unhandled command", cmd, data)
        end
        app = App:new(conn, handle_commands):with_window(80, 24)

        assert(app.window)
        assert(app.cache)

        local function run()
            if app == nil then
                return
            end

            app.conn:send("render", "")

            vim.defer_fn(run, 500)
        end
        run()
    end)

end

function CLOSE()
    assert(app ~= nil, "app not started")
    app:close()
end
