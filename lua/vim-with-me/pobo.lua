local Enum = require("vim-with-me.enum");

local Pobo = {}

Pobo.type_idx = 0;
Pobo.statusline_idx = 1;
Pobo.statusline_length = 50;
Pobo.cost_idx = 51;
Pobo.data_idx = 53;
Pobo.data_length = 200;

local function slice(t, start, stop)
    local out = {}
    for idx = start, stop do
        table.insert(out, t[idx])
    end
    return out
end

-- stackoverflow, though dangerous behind the wheel can still serve a purpose
local function utf8_from(t, start, length)
    local bytearr = {}
    local idx = start
    local stop = start + length
    while idx <= stop and t[idx] ~= 0 and t[idx] ~= nil do
        local v = t[idx]
        local utf8byte = v < 0 and (0xff + v + 1) or v
        table.insert(bytearr, string.char(utf8byte))
        idx = idx + 1
    end
    return table.concat(bytearr)
end

local function parse_string(line, start, stop)
    local length = start
    while line[start + length] ~= 0 and (start + length) <= stop do
        length = length + 1
    end

    -- Why am i still doing this?  I feel like I should stop doing this...
    return utf8_from(line, start, length + 1)
end

function Pobo:new(data, offset)
    local obj = {
        _data = data,
        offset = offset or 1,
    }

    setmetatable(obj, self)
    self.__index = self

    return obj
end

function Pobo:get_data()
    return parse_string(self._data, self.offset + Pobo.data_idx, Pobo.data_length)
end

function Pobo:get_status()
    return parse_string(self._data, self.offset + Pobo.statusline_idx, Pobo.statusline_length)
end

function Pobo:get_command()
    return self._data[self.offset + Pobo.type_idx]
end

function Pobo:get_cost()
    return self._data[self.offset + Pobo.cost_idx] * 256 +
        self._data[self.offset + Pobo.cost_idx + 1]
end


return Pobo

