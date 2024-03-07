local window = require("vim-with-me.window")
local cache = require("vim-with-me.window.cache")

---@class TowerOffense
---@field _window_details WindowDetails | nil
---@field _offset WindowPosition
---@field _display_cache DisplayCache
---@field _window_cache DisplayCache
local TowerOffense = {}
TowerOffense.__index = TowerOffense

---@param width number
---@param height number
---@return TowerOffense
function TowerOffense.new(width, height)
    width = width or 80
    height = height or 24

    local self = setmetatable({
        _window_details = nil,
        _display_cache = cache.new(width, height),
        _offset = window.create_window_dimensions(2, 2),
    }, TowerOffense)

    return self
end

function TowerOffense:resize()
    if self._window_details == nil then
        return
    end

    if not vim.api.nvim_win_is_valid(self._window_details.win_id) then
        self:close()
        return
    end

    window.resize(self._window_details)
end

function TowerOffense:close()
    assert(self._window_details ~= nil, "window already closed")

    self.closing = true
    window.close_window(self._window_details)
    self.closing = false
    self._window_details = nil
end

function TowerOffense:_render()
    self._window_cache:map(self._display_cache)
    vim.api.nvim_buf_set_lines(
        self._window_details.buffer,
        0,
        -1,
        false,
        self._window_cache:to_string_rows()
    )
end

function TowerOffense:start()
    assert(self._window_details == nil, "window already started")

    self._window_details = window.create_window(self._offset)
    self._window_cache = cache.new(
        self._window_details.dim.width,
        self._window_details.dim.height
    )

    self._display_cache:clear()
    self:_render()
    vim.api.nvim_set_current_win(self._window_details.win_id)
end

---@param mark string
---@param x number
---@param y number
function TowerOffense:place(mark, x, y)
    assert(type(mark) == "string", "mark must be a string")
    assert(#mark == 1, "mark must be a single character")
    assert(type(x) == "number", "x must be a number")
    assert(type(y) == "number", "y must be a number")
    assert(
        self._window_details ~= nil,
        "please call #open first before placing a tower"
    )
    assert(
        self._window_details.dim.width >= x,
        "x must be less than or equal to the width"
    )
    assert(
        self._window_details.dim.height >= y,
        "y must be less than or equal to the height"
    )
    assert(x > 0, "x must be greater than 0")
    assert(y > 0, "y must be greater than 0")

    self._display_cache:place(x, y, "T")
    self:_render()
end

return TowerOffense
