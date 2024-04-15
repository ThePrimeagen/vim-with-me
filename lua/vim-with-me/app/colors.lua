local utils = require("vim-with-me.utils")
local ns = utils.namespace

---@param color VWMColor
---@return string
local function color_to_string(color)
    return string.format("%02x%02x%02x", color.red, color.green, color.blue)
end

---@class VWMColors
---@field foreground string
---@field background string

---@class VWMColorSet
---@field seen table<string, string>
---@field count number
---@field buffer number
---@field rows number
---@field cols number
local ColorSet = {}
ColorSet.__index = ColorSet

---@param buffer number
---@param rows number
---@param cols number
---@param base_foreground VWMColor
---@param base_background VWMColor
function ColorSet:new(buffer, rows, cols, base_foreground, base_background)
    local color = setmetatable({
        seen = {},
        count = 0,
        rows = rows,
        cols = cols,
        buffer = buffer,
    }, self)

    color:_get_name(base_foreground)
    color:_get_name(base_background)

    color:color_range(base_foreground, 1, 0, rows + 1, cols)
    color:color_range(base_background, 1, 0, rows + 1, cols)

    return color
end

---@param vwm_color VWMColor
---@return string
function ColorSet:_get_name(vwm_color)
    local color = color_to_string(vwm_color)

    if self.seen[color] == nil then
        local name = string.format("VWM_%d", self.count)
        self.seen[color] = name
        self.count = self.count + 1

        local opts = {}
        opts[vwm_color.foreground and "foreground" or "background"] = string.format("#%s", color)

        vim.api.nvim_set_hl(ns, name, opts)
    end
    return self.seen[color]
end

--- I am very worried about this interface, i think that its pretty bad as far
--- as performance is going to go
---@param cell CellWithLocation
function ColorSet:color_cell(cell)
    local fg_name = self:_get_name(cell.cell.foreground)
    local bg_name = self:_get_name(cell.cell.background)
    local loc = cell.loc
    local col = loc.col
    local row = loc.row

    vim.api.nvim_buf_add_highlight(self.buffer, ns, fg_name, row, col, col + 1)
    vim.api.nvim_buf_add_highlight(self.buffer, ns, bg_name, row, col, col + 1)
end

---@param color VWMColor
---@param row_start number
---@param col_start number
---@param row_end number
---@param col_end number
function ColorSet:color_range(color, row_start, col_start, row_end, col_end)
    local name = self:_get_name(color)
    local function add(row, s, e)
        vim.api.nvim_buf_add_highlight(
            self.buffer,
            ns,
            name,
            row,
            s,
            e
        )
    end

    if row_start == row_end then
        add(row_start, col_start, col_end)
    end

    add(row_start, col_start, self.cols)

    for i = row_start + 1, row_end - 1 do
        add(i, 0, self.cols)
    end

    add(row_end, 0, col_end)
end

return ColorSet
