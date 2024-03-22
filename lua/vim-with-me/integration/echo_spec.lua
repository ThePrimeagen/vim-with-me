local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")

describe("vim with me", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)
    it("integartion testing", function()
        local tcp = int_utils.create_test_conn("echo_server", 42069)
        local next = int_utils.create_tcp_next(tcp)
        tcp:send("hello", "world")

        local hello_back = next()
        eq(hello_back ~= nil, true)
        -- not needed
        if hello_back == nil then
            return
        end
        eq(hello_back.command, "world")
        eq(hello_back.data, "hello")
    end)
end)
