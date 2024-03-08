local system = vim.system or require("vim-with-me.system")
local TCP = require("vim-with-me.tcp").TCP

---@param name string
---@param port number
---@return TCP
local function test_app(name, port)
    local done_building = false
    system.run({"go", "build", "-o", name, string.format("./test/%s/main.go", name)}, {
    }, function()
        done_building = true
    end)
    vim.wait(1000, function()
        return done_building
    end)

    system.run({string.format("./%s", name), "--port", tostring(port)}, {
        stdout = function(_, data)
            print(data)
        end
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

return {
    test_app = test_app,
}
