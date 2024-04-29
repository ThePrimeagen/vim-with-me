local LOCATION_ENCODING_LENGTH = 2
local FOREGROUND = 1
local COLOR = 3
local COLOR_ENCODING_LENGTH = FOREGROUND + COLOR
local CELL_ENCODING_LENGTH = COLOR_ENCODING_LENGTH * 2 + 1
local CELL_AND_LOC_ENCODING_LENGTH = CELL_ENCODING_LENGTH
    + LOCATION_ENCODING_LENGTH

---@class VWMLocation
---@field row number
---@field col number

---@class Cell
---@field foreground VWMColor
---@field background VWMColor
---@field value      string

-- {loc={row = 1, col = 1}, cell={ foreground = { red = 12, green = 34, blue = 56, foreground = true}, background = { red = 12, green = 34, blue = 56, foreground = true}},
---@class VWMColor
---@field red        number
---@field blue       number
---@field green      number
---@field foreground boolean

---@class CellWithLocation
---@field loc VWMLocation
---@field cell Cell

local M = {}

---@param data string
---@return VWMLocation
function M.parse_location(data)
    assert(
        #data >= LOCATION_ENCODING_LENGTH,
        "not enough data provided to location parse"
    )
    return {
        row = string.byte(data, 1, 1),
        col = string.byte(data, 2, 2),
    }
end

---@param data string
---@return VWMColor
function M.parse_color(data)
    assert(
        #data >= COLOR_ENCODING_LENGTH,
        "not enough data provided to color parse"
    )
    return {
        foreground = string.byte(data, 1, 1) == 1,
        red = string.byte(data, 2, 2),
        green = string.byte(data, 3, 3),
        blue = string.byte(data, 4, 4),
    }
end

---@param data string
---@return Cell
function M.parse_cell(data)
    assert(
        #data >= CELL_ENCODING_LENGTH,
        "not enough data provided to cell parse"
    )
    return {
        value = string.sub(data, 1, 1),
        foreground = M.parse_color(string.sub(data, 2)),
        background = M.parse_color(string.sub(data, 2 + COLOR_ENCODING_LENGTH)),
    }
end

---@param data string
---@return CellWithLocation
function M.parse_cell_with_location(data)
    assert(
        #data >= CELL_AND_LOC_ENCODING_LENGTH,
        "incomplete partial render string provided: " .. #data
    )

    local loc = M.parse_location(data)
    local str = string.sub(data, 1 + LOCATION_ENCODING_LENGTH)
    local cell = M.parse_cell(str)

    return {
        loc = loc,
        cell = cell,
    }
end

---@param data string
---@return CellWithLocation[]
function M.parse_partial_renders(data)
    assert(
        #data % CELL_AND_LOC_ENCODING_LENGTH == 0,
        "incomplete partial render string provided: " .. #data
    )

    local renders = {}
    for i = 1, #data, CELL_AND_LOC_ENCODING_LENGTH do
        local str = string.sub(data, i)
        table.insert(renders, M.parse_cell_with_location(str))
    end

    return renders
end

return M
