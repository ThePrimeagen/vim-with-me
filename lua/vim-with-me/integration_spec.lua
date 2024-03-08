local system = vim.system or require("vim-with-me.system")
local eq = assert.are.same
local TCP = require("vim-with-me.tcp").TCP

describe("vim with me", function()
    it("integartion testing", function()
        local port = 42075
        local done_building = false
        system.run({"go", "build", "-o", "test_server", "./cmd/test_server/main.go"}, {
        }, function()
            done_building = true
        end)
        vim.wait(1000, function()
            return done_building
        end)

        system.run({"./test_server", "--port", tostring(port)}, {
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

        local hello_back = nil
        tcp:listen(function(command, data)
            hello_back = {
                command = command,
                data = data,
            }
        end)
        tcp:send("hello", "world")

        vim.wait(1000, function()
            return hello_back ~= nil
        end)

        eq(hello_back ~= nil, true)
        eq(hello_back.command, "world")
        eq(hello_back.data, "hello")
    end)
end)
