-- luacheck: globals describe it assert
local eq = assert.are.same
local utils = require("vim-with-me.tcp.utils")

describe("vim with me :: tcp.utils", function()
    it("from", function()
        eq(97201, utils.from_tcp_int("0$%"))
    end)

    it("to", function()
        eq("0$%", utils.to_tcp_int(97201))
    end)

    it("from(to)", function()
        eq(69420, utils.from_tcp_int(utils.to_tcp_int(69420)))
    end)

    it("color compression", function()
        local one = "#FF0000"
        local two = "#00FF00"
        local three = "#0000FF"

        local compression = utils.ColorCompression:new()
        eq(one, compression:decompress(one))
        eq(one, compression:decompress("$"))

        local ok, _ = pcall(compression.decompress, compression, "%")
        eq(ok, false)

        eq(two, compression:decompress(two))
        eq(two, compression:decompress("%"))

        ok, _ = pcall(compression.decompress, compression, "&")
        eq(ok, false)

        eq(three, compression:decompress(three))
        eq(three, compression:decompress("&"))

        ok, _ = pcall(compression.decompress, compression, "'")
        eq(ok, false)

        eq(compression.size, 3)
    end)
end)
