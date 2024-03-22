local system = vim.system or require("vim-with-me.system")
local TCP = require("vim-with-me.tcp").TCP

---@param tcp TCP
---@return fun(): {command: string, data: string} | nil
local function create_tcp_next(tcp)
    local received = {}
    tcp:listen(function(command, data)
        table.insert(received, {command = command, data = data})
    end)

    return function()
        vim.wait(1000, function()
            return #received > 0
        end)

        if #received == 0 then
            return nil
        end

        local out = table.remove(received, 1)
        return out
    end
end

---@param name string
---@param port number
---@param stdout (fun(_: any, data: string): nil) | nil
---@return TCP
local function create_test_conn(name, port, stdout)
    local done_building = false
    system.run({"go", "build", "-o", name, string.format("./test/%s/main.go", name)}, {
    }, function()
        done_building = true
    end)
    vim.wait(1000, function()
        return done_building
    end)

    system.run({string.format("./%s", name), "--port", tostring(port)}, {
        stdout = stdout,
    })
    vim.wait(100)

    local connected = false
    local tcp = TCP:new({
        host = "127.0.0.1",
        port = port,
        retry_count = 3,
    })
    tcp:start(function()
        connected = true
    end)

    vim.wait(1000, function()
        return connected == true
    end)

    assert(connected, "could not connect to server")
    return tcp
end

-- --- TODO: This is for me to think about before i create this... I don't want to create complex testing application
-- ---@class VWMTestApp
-- ---@field public app VWMApp
-- ---@field public conn TCP
-- local TestApp = {}
-- TestApp.__index = TestApp
--
-- ---@param app VWMApp
-- ---@param conn TCP
-- ---@return VWMTestApp
-- function TestApp:new(app, conn)
--     return setmetatable({
--         app = app,
--         conn = conn,
--     }, self)
-- end
--
-- ---@param tcp TCP
-- local function create_app(tcp)
--     local app = App:new(tcp)
--
--     local latest_command = nil
--     local function on_command(cmd, data)
--
--     end
--
--     return app, function()
--     end
-- end

return {
    create_tcp_next = create_tcp_next,
    create_test_conn = create_test_conn,
}
