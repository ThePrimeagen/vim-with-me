local ZERO = string.byte("$")
local BASE = 90

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


---@param length number
---@return string
local function to_tcp_int(length)
    assert(length >= 0, "negative numbers are not allowed")

    local n = length
    local out = {}
    while n > 0 do
        local v = n % BASE
        table.insert(out, 1, string.char(ZERO + v))
        n = math.floor(n / BASE)
    end

    return table.concat(out, "")
end

--- This produces a base10 number from the string vim-with-me compressed number (base 90)
---@param str string
---@return number
local function from_tcp_int(str)
    local int = 0
    for i = 1, #str do
        local base = math.pow(BASE, #str - i)
        local value = string.byte(str, i) - ZERO

        int = int + value * base
    end

    return int
end

--- @class ColorCompression
--- @field table table<string>
--- @field size number

local ColorCompression = {}
ColorCompression.__index = ColorCompression

function ColorCompression:new()
    local compression = setmetatable({
        table = {},
        size = 0
    }, self)
    return compression
end

---@param color string
---@return string
function ColorCompression:decompress(color)
    if string.sub(color, 1, 1) == "#" then
        table.insert(self.table, color)
        self.size = self.size + 1

        return color
    end

    -- we plus one because we are using 1 based indexing and 1 based indexing is a gift that keeps giving skill issues
    local index = from_tcp_int(color) + 1
    local value = self.table[index]

    assert(value ~= nil, string.format("index: %d does not exist in table(%d)", index, #self.table))
    return value
end

function ColorCompression:clear()
    self.table = {}
    self.size = 0
end

return {
    to_tcp_int = to_tcp_int,
    from_tcp_int = from_tcp_int,

    ColorCompression = ColorCompression,
    to_string = to_string,
}


