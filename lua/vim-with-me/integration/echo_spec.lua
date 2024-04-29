local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")

describe("vim with me", function()
    before_each(int_utils.before_each)
    after_each(int_utils.after_each)

    it("integartion testing", function()
        int_utils.create_test_server("echo_server", 42069)
        local tcp = int_utils.create_tcp_connection(42069)

        local next = int_utils.create_tcp_next(tcp)
        next()

        tcp:send({ command = 1, data = "world" })

        local hello_back = next()

        eq({ command = 1, data = "world" }, hello_back)
    end)
end)
