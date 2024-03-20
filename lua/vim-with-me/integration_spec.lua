local eq = assert.are.same
local utils = require("vim-with-me.test.utils")

describe("vim with me", function()
    local function create_next(tcp)
        local received = nil
        tcp:listen(function(command, data)
            received = {command, data}
        end)

        return function()
            vim.wait(1000, function()
                return received ~= nil
            end)

            local out = received
            received = nil
            return out
        end
    end

    it("integartion testing", function()
        local tcp = utils.test_app("echo_server", 42069)
        local next = create_next(tcp)
        tcp:send("hello", "world")

        local hello_back = next()
        eq(hello_back ~= nil, true)
        eq(hello_back.command, "world")
        eq(hello_back.data, "hello")
    end)

    it("", function()
        local tcp = utils.test_app("partial_render_server", 42070)
        local next = create_next(tcp)

        eq({
            command = "open-window",
            data = "24:80",
        }, next())

        tcp:send("hello", "")

        eq({
            command = "p",
            data = "0:1:1",
        }, next())

    end)
end)
