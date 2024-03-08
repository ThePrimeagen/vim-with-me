local eq = assert.are.same
local utils = require("vim-with-me.test.utils")

describe("vim with me", function()
    it("integartion testing", function()
        local tcp = utils.test_app("echo_server", 42069)
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
