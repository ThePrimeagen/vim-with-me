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
        app = App:new(conn, handle_commands)
    end)

end

function CLOSE()
    assert(app ~= nil, "app not started")
    app:close()
end
