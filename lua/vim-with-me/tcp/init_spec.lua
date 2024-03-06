-- luacheck: globals describe it assert
local eq = assert.are.same
local tcp = require("vim-with-me.tcp")

describe("vim with me :: tcp", function()
    it("parse", function()
        local str = "1:5:r:hel"
        local next_chunk, command, data = tcp.parse(str, 5, 5)

        eq(next_chunk, "")
        eq(command, "r")
        eq(data, "hel")
    end)


    it("process packets chunks", function()
        local chunks = {
            "1:",
            "5:r:he",
            "l1:10:r",
            ":lo wo",
            "rld",
        }

        local packets = tcp.process_packets()

        local command, data = packets(chunks[1])
        eq(command, nil)
        eq(data, nil)

        command, data = packets(chunks[2])
        eq(command, nil)
        eq(data, nil)

        command, data = packets(chunks[3])
        eq(command, "r")
        eq(data, "hel")

        command, data = packets(chunks[4])
        eq(command, nil)
        eq(data, nil)

        command, data = packets(chunks[5])
        eq(command, "r")
        eq(data, "lo world")
    end)

    it("process packets multiple in one chunk", function()
        local chunks = "1:5:r:hel1:10:r:lo world"
        local packets = tcp.process_packets()

        local command, data = packets(chunks)
        eq(command, "r")
        eq(data, "hel")
        command, data = packets(chunks)
        eq(command, "r")
        eq(data, "lo world")
    end)
end)
