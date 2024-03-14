---@class DisplayCache
---@field data string[][]
---@field rows number
---@field cols number
local DisplayCache = {}
DisplayCache.__index = DisplayCache

---@param rows number
---@param cols number
---@return DisplayCache
function DisplayCache:new(rows, cols)
    local data = {}
    for _ = 1, cols do
        local row = {}
        for _ = 1, rows do
            table.insert(row, " ")
        end
        table.insert(data, row)
    end

    return setmetatable({
        data = data,
        rows = rows,
        cols = cols,
    }, self)
end
---@param partial PartialRender
function DisplayCache:partial(partial)
    self:place(partial.col, partial.row, partial.value)
end


---@param row number
---@param col number
---@param item string
function DisplayCache:place(row, col, item)
    assert(type(row) == "number", "x must be a number")
    assert(type(col) == "number", "y must be a number")
    assert(type(item) == "string", "item must be a string")
    assert(#item == 1, "item must be a single character")
    assert(row >= 1, "x must be greater than or equal to 1")
    assert(col >= 1, "y must be greater than or equal to 1")
    assert(row <= self.rows, "row must be less than or equal to the rows")
    assert(col <= self.cols, "col must be less than or equal to the cols")
    self.data[col][row] = item
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
    for y = 1, self.cols do
        for x = 1, self.rows do
            self:place(x, y, " ")
        end
    end
end

---@param str string
function DisplayCache:from_string(str)
    assert(
        #str == self.rows * self.cols,
        "string must be the same length as the cache"
    )

    local y = 1
    while y <= self.cols do
        local line = string.sub(str, y, y + self.rows - 1)

        for x = 1, #line do
            local char = string.sub(
                str,
                (y - 1) * self.rows + x,
                (y - 1) * self.rows + x
            )
            self.data[y][x] = char
        end

        y = y + 1
    end
end

return DisplayCache
