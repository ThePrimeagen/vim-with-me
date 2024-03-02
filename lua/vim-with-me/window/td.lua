local window = require("vim-with-me.window.window")

---@class TowerOffenseWindow
---@field _window_details WindowDetails | nil
---@field _closing boolean
---@field _offset WindowPosition
local TowerOffenseWindow = {}
TowerOffenseWindow.__index = TowerOffenseWindow

function TowerOffenseWindow.new()
    local self = setmetatable({
        _window_details = nil,
        _closing = false,
        _offset = window.create_window_offset(2, 2),
    }, TowerOffenseWindow)
    return self
end

function TowerOffenseWindow:resize()
    if self._window_details == nil then
        return
    end

    if not vim.api.nvim_win_is_valid(self._window_details.win_id) then
        self:close()
        return
    end

    window.resize(self._window_details)
end

function TowerOffenseWindow:close()
    if self._window_details ~= nil and not self._closing then
        self.closing = true
        window.close_window(self._window_details)
        self.closing = false
        self._window_details = nil
    end
end

function TowerOffenseWindow:toggle()
    if self._window_details == nil then
        self._window_details = window.create_window(self._offset)

        window.on_close(self._window_details, function()
            self._window_details = nil
            self._closing = false
        end)
        window.refocus(self._window_details)
    else
        self:close()
    end
end

return TowerOffenseWindow
