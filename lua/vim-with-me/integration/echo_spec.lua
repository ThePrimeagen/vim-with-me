local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")

describe("vim with me", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)
    it("integartion testing", function()
        int_utils.create_test_server("echo_server", 42069)
        local tcp = int_utils.create_tcp_connection(42069)

        local next = int_utils.create_tcp_next(tcp)
        tcp:send({ command = 1, data = "world" })

        local hello_back = next()
        eq(hello_back ~= nil, true)
        -- not needed
        if hello_back == nil then
            return
        end

        eq(hello_back.command, 1)
        eq(hello_back.data, "world")
    end)
end)
