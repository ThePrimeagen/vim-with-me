---@class DisplayCache
---@field data string[][]
---@field width number
---@field height number
local DisplayCache = {}
DisplayCache.__index = DisplayCache

---@param width number
---@param height number
---@return DisplayCache
function DisplayCache.new(width, height)
    local data = {}
    for _ = 1, height do
        local row = {}
        for _ = 1, width do
            table.insert(row, " ")
        end
        table.insert(data, row)
    end

    local self = setmetatable({
        data = data,
        width = width,
        height = height,
    }, DisplayCache)
    return self
end

---@param x number
---@param y number
---@param item string
function DisplayCache:place(x, y, item)
    assert(type(x) == "number", "x must be a number")
    assert(type(y) == "number", "y must be a number")
    assert(type(item) == "string", "item must be a string")
    assert(#item == 1, "item must be a single character")
    assert(x >= 1, "x must be greater than or equal to 1")
    assert(y >= 1, "y must be greater than or equal to 1")
    assert(x <= self.width, "x must be less than or equal to the width")
    assert(y <= self.height, "y must be less than or equal to the height")
    self.data[y][x] = item
end

---@return string[]
function DisplayCache:to_string_rows()
    local out = {}
    for _, row in ipairs(self.data) do
        table.insert(out, table.concat(row))
    end

    return out
end

---@param other DisplayCache
---@param x_start number | nil
---@param y_start number | nil
function DisplayCache:map(other, x_start, y_start)
    x_start = x_start or 1
    y_start = y_start or 1
    for y, row in ipairs(other.data) do
        for x, item in ipairs(row) do
            self:place(x + x_start, y + y_start, item)
        end
    end
end

function DisplayCache:clear()
    for y = 1, self.height do
        for x = 1, self.width do
            self:place(x, y, " ")
        end
    end
end

---@param str string
function DisplayCache:from_string(str)
    assert(
        #str == self.width * self.height,
        "string must be the same length as the cache"
    )

    local y = 1
    while y <= self.height do
        local line = string.sub(str, y, y + self.width - 1)

        for x = 1, #line do
            local char = string.sub(
                str,
                (y - 1) * self.width + x,
                (y - 1) * self.width + x
            )
            self.data[y][x] = char
        end

        y = y + 1
    end
end

return DisplayCache
