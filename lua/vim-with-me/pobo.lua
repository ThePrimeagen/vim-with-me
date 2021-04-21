local Enum = require("vim-with-me.enum");

local States = Enum({
    Start = 0,
    MiddleEarth = 1,
    Finisher = 2,
})

local InverseStates = {
    States.Start,
    States.MiddleEarth,
    States.Finisher,
}


local Pobo = {}

function Pobo:new(data, offset)
end

function Pobo:getCommand()
    return InverseStates[self.data[self.offset]]
end

return Pobo

