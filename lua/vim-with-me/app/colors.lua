local Locations = require("vim-with-me.location")

---@type VWMColor
local DEFAULT_BACKGROUND = {
    foreground = false,
    red = 0,
    green = 0,
    blue = 0,
}

---@type VWMColor
local DEFAULT_FOREGROUND = {
    foreground = true,
    red = 50,
    green = 75,
    blue = 120,
}

---@param color VWMColor
---@return string
local function color_to_string(color)
    return string.format("#%02x%02x%02x", color.red, color.green, color.blue)
end

---@class VWMColors
---@field foreground string
---@field background string

---@class VWMColorSet
---@field seen table<string, string>
---@field namespaces table<table<string>>
---@field count number
---@field buffer number
---@field rows number
---@field cols number
local ColorSet = {}
ColorSet.__index = ColorSet

---@param window WindowDetails
---@param base_foreground VWMColor | nil
---@param base_background VWMColor | nil
function ColorSet:new(window, base_foreground, base_background)
    base_background = base_background or DEFAULT_BACKGROUND
    base_foreground = base_foreground or DEFAULT_FOREGROUND

    local rows = window.dim.height
    local cols = window.dim.width
    local color = setmetatable({
        seen = {},
        count = 0,
        rows = rows,
        namespaces = {},
        cols = cols,
        buffer = window.buffer,
    }, self)

    color:_get_name(base_foreground)
    color:_get_name(base_background)

    for r = 1, rows do
        local namespace_row = {}
        table.insert(color.namespaces, namespace_row)

        for c = 1, cols do
            local namespace =
                vim.api.nvim_create_namespace(string.format("VWM%d_%d", r, c))
            table.insert(namespace_row, namespace)

            color:color_cell({
                loc = Locations.from_cache(r, c),
                cell = {
                    foreground = base_foreground,
                    background = base_background,
                    value = "E", -- nice
                },
            })
        end
    end

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
        opts[vwm_color.foreground and "foreground" or "background"] = color

        vim.api.nvim_set_hl(0, name, opts)
    end
    return self.seen[color]
end

--- I am very worried about this interface, i think that its pretty bad as far
--- as performance is going to go
---@param cell CellWithLocation
function ColorSet:color_cell(cell)
    local fg_name = self:_get_name(cell.cell.foreground)
    local bg_name = self:_get_name(cell.cell.background)

    local row, col = Locations.to_cache(cell)
    local ns = self.namespaces[row][col]

    row, col = Locations.to_line_api(cell)

    vim.api.nvim_buf_clear_namespace(self.buffer, ns, row, row + 1)

    vim.highlight.range(
        self.buffer,
        ns,
        fg_name,
        { row, col },
        { row, col + 1 }
    )
    vim.highlight.range(
        self.buffer,
        ns,
        bg_name,
        { row, col },
        { row, col + 1 }
    )

    --[[
    vim.api.nvim_buf_add_highlight(self.buffer, ns, fg_name, row, col - 1, col + 1)
    vim.api.nvim_buf_add_highlight(self.buffer, ns, fg_name, row - 1, col - 1, col + 1)
    vim.api.nvim_buf_add_highlight(self.buffer, ns, fg_name, row + 1, col - 1, col + 1)
    --]]
end

return ColorSet
