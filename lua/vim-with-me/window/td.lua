local window = require("vim-with-me.window.window")

---@class TowerOffense
---@field _window_details WindowDetails | nil
---@field _closing boolean
---@field _offset WindowPosition
local TowerOffense = {}
TowerOffense.__index = TowerOffense

function TowerOffense.new()
    local self = setmetatable({
        _window_details = nil,
        _closing = false,
        _offset = window.create_window_offset(2, 2),
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
    if self._window_details ~= nil and not self._closing then
        self.closing = true
        window.close_window(self._window_details)
        self.closing = false
        self._window_details = nil
    end
end

function TowerOffense:start()
    if self._window_details ~= nil then
        return
    end
    self._window_details = window.create_window(self._offset)

    local buf = {}
    for _ = 1, self._window_details.dim.height do
        table.insert(buf, self:create_row(0))
    end

    vim.api.nvim_buf_set_lines(self._window_details.buffer, 0, -1, false, buf)
    vim.api.nvim_set_current_win(self._window_details.win_id)
end

---@param x number
---@return string
function TowerOffense:create_row(x)
    local str = ""
    for i = 1, self._window_details.dim.width do
        if i == x then
            str = str .. "X"
        else
            str = str .. " "
        end
    end

    return str
end


function TowerOffense:place(x, y)
    if type(x) ~= "number" or type(y) ~= "number" then
        print("invalid point")
        return
    end

    assert(self._window_details ~= nil, "please call #open first before placing a tower")
    if self._window_details.dim.width < x or self._window_details.dim.height < y then
        print("point ignored since it is too girthy")
        return
    end

    if x == 0 or y == 0 then
        print("point ignored since it is too small, growl")
        return
    end

    local str = self:create_row(x)
    vim.api.nvim_buf_set_lines(self._window_details.buffer, y, y + 1, false, {str})
end

return TowerOffense
