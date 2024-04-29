local Locations = require("vim-with-me.location")

---@class DisplayCache
---@field data string[][]
---@field rows number
---@field cols number
---@field dirty table<table<boolean>>
local DisplayCache = {}
DisplayCache.__index = DisplayCache

---@param dim WindowPosition
---@return DisplayCache
function DisplayCache:new(dim)
    local data = {}
    local dirty = {}
    for _ = 1, dim.height do
        local row = {}
        local dirty_row = {}
        for _ = 1, dim.width do
            table.insert(row, " ")
            table.insert(dirty_row, false)
        end
        table.insert(data, row)
        table.insert(dirty, dirty_row)
    end

    return setmetatable({
        data = data,
        rows = dim.height,
        cols = dim.width,
        dirty = dirty,
    }, self)
end
---@param partial CellWithLocation
function DisplayCache:partial(partial)
    local row, col = Locations.to_cache(partial)
    self:place(row, col, partial.cell.value)
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
    assert(
        row <= self.rows,
        "row must be less than or equal to the rows: "
            .. row
            .. " "
            .. self.rows
    )
    assert(
        col <= self.cols,
        "col must be less than or equal to the cols: "
            .. col
            .. " "
            .. self.cols
    )
    self.data[row][col] = item
    self.dirty[row][col] = true
end

---@param window WindowDetails
function DisplayCache:render_into(window)
    --- TODO: Make this even more fantastic by doing dirty contiguous regions instead of dirty cells
    for r, dirty_row in ipairs(self.dirty) do
        for c, dirty in ipairs(dirty_row) do

            if dirty then
                local loc = Locations.from_cache(r, c)
                local ok, _ = pcall(
                    vim.api.nvim_buf_set_text,
                    window.buffer,
                    loc.row,
                    loc.col,
                    loc.row,
                    loc.col + 1,
                    { self.data[r][c] }
                )

                assert(
                    ok,
                    string.format(
                        "unable to nvim_buf_set_text: %d %d %d %d",
                        self.rows,
                        self.cols,
                        loc.row,
                        loc.col
                    )
                )

                dirty_row[c] = false
            end
        end
    end
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
    while y <= self.rows do
        local base = ((y - 1) * self.cols) + 1
        local line = string.sub(str, base, base + self.cols - 1)

        for x = 1, #line do
            local char = string.sub(line, x, x)
            self.data[y][x] = char
        end

        y = y + 1
    end
end

return DisplayCache
