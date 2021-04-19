local EnvClient = require("vim-with-me.envelope")
local TcpMock = require("vim-with-me.test.mocks.tcp")
local ass = assert.are

local States = EnvClient.States

local function create()
    local tcpMock = TcpMock:new()
    local env = EnvClient:new(tcpMock)
    env:connect()

    return tcpMock, env
end

describe("", function()
    it("should be able to parse a single payload packet.", function()
        local tcpMock, env = create()
        local called = false
        local data = nil
        env:on("data", function(line)
            called = true
            data = line
        end)

        tcpMock:emit("data", {
            0,
            5,
            1,
            2,
            3,
            4,
            5,
        })

        ass.same(called, true)
        ass.same(data, {
            1,
            2,
            3,
            4,
            5,
        })

        -- TODO: Why did I mark these as private when I am testing them??
        -- THIS IS THE WORST WAY TO DO THIS.
        -- but its also really nice to make sure I didn't screw up the internal
        -- state of this TCP wrapper
        ass.same(env._current_length, 0)
        ass.same(env._current_env, nil)
        ass.same(env._state, States.WaitingForHeader)
    end)

    it("can parse back to back 93 and 94 packets", function()
        print("TEST 2")
        local tcpMock, env = create()
        local called = 0
        local data = {}
        env:on("data", function(line)
            called = called + 1
            table.insert(data, line)
        end)

        tcpMock:emit("data", {
            0,
            1,
            93,
            0,
            1,
            94,
        })

        ass.same(called, 2)
        ass.same(data, {{
            93,
        }, {
            94,
        }})

        -- TODO: Why did I mark these as private when I am testing them??
        -- THIS IS THE WORST WAY TO DO THIS.
        -- but its also really nice to make sure I didn't screw up the internal
        -- state of this TCP wrapper
        ass.same(env._current_length, 0)
        ass.same(env._current_env, nil)
        ass.same(env._state, States.WaitingForHeader)
    end)
end)

