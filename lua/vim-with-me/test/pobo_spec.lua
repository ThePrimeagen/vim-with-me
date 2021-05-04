local Pobo = require("vim-with-me.pobo")
local ass = assert.are

describe("Pobo baybeee", function()
    -- Taken from dat
    -- type = 4 (giveaway)
    -- status: ThePrimeagen: Thanks for entering the giveaway
    -- cost = 65535
    -- data = ThePrimeagen
    local test_data = { 4, 84, 104, 101, 80, 114, 105, 109, 101, 97, 103, 101, 110, 58, 32, 84, 104, 97, 110, 107, 115, 32, 102,
 111, 114, 32, 101, 110, 116, 101, 114, 105, 110, 103, 32, 116, 104, 101, 32, 103, 105, 118, 101, 97, 119,
 97, 121, 0, 0, 0, 0, 255, 255, 84, 104, 101, 80, 114, 105, 109, 101, 97, 103, 101, 110, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
, 0, 0, 0, 0, 0 }

    it("should decode a pobo properly.", function()

        local type = 4
        local status = "ThePrimeagen: Thanks for entering the giveaway"
        local cost = 65535
        local data = "ThePrimeagen"
        local pobo = Pobo:new(test_data)

        ass.same(pobo:get_data(), data)

        ass.same(pobo:get_command(), type)
        ass.same(pobo:get_status(), status)
        ass.same(pobo:get_cost(), cost)

    end)

end)


