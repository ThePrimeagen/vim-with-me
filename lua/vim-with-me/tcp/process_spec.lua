-- luacheck: globals describe it assert
local eq = assert.are.same
local tcp = require("vim-with-me.tcp.process")

local function to_string(...)
    local out = {}
    for i = 1, select("#", ...) do
        local v = select(i, ...)
        if type(v) == "number" then
            table.insert(out, string.char(v))
        elseif type(v) == "string" then
            table.insert(out, v)
        else
            assert(false, "should never provide anything other than numbers, strings, or tables of strings and numbers")
        end
    end
    return table.concat(out, "")
end

describe("vim with me :: tcp.process", function()
    it("parse_big_endian_16", function()
        eq(0x45, tcp.parse_big_endian_16(to_string(0, 69)))
        eq(0x4500, tcp.parse_big_endian_16(to_string(69, 0)))
        eq(0x4545, tcp.parse_big_endian_16(to_string(69, 69)))
    end)

    it("process packets chunks", function()
        local chunks = {
            to_string(1), -- version
            to_string(69), -- cmd
            tcp.to_big_endian_16(5),
            to_string("n", "o", "i", "c", "e", 1),
            to_string(72, tcp.to_big_endian_16(3), "f", "o", "o"), -- cmd, len, data
        }

        local packets = tcp.process_packets()

        local command = packets(chunks[1])
        eq(nil, command)

        command = packets(chunks[2])
        eq(nil, command)

        command = packets(chunks[3])
        eq(nil, command)

        command = packets(chunks[4])
        eq({command = 69, data = "noice"}, command)

        command = packets(chunks[5])
        eq({command = 72, data = "foo"}, command)

        eq(nil, packets(""))
    end)

    it("process packets multiple in one chunk", function()
        local chunk = table.concat({
            to_string(1), -- version
            to_string(69), -- cmd
            tcp.to_big_endian_16(5),
            to_string("n", "o", "i", "c", "e", 1),
            to_string(72, tcp.to_big_endian_16(3), "f", "o", "o"), -- cmd, len, data
        }, "")

        local packets = tcp.process_packets()
        eq({command = 69, data = "noice"}, packets(chunk))
        eq({command = 72, data = "foo"}, packets())
    end)

    it("version mismatch should cause an error", function()
        local chunk = table.concat({
            to_string(2), -- version
            to_string(69), -- cmd
            tcp.to_big_endian_16(5),
            to_string("n", "o", "i", "c", "e"),
        }, "")
        local packets = tcp.process_packets()
        local ok, command = pcall(packets, chunk)
        eq(ok, false)
    end)
end)
