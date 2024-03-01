local window = require("vim-with-me.window.window")
local group = vim.api.nvim_create_augroup(false, "vim-with-me.window")

---@class TowerOffenseWindow
---@field _window_details WindowDetails | nil
---@field _closing boolean
local TowerOffenseWindow = {}
TowerOffenseWindow.__index = TowerOffenseWindow

function TowerOffenseWindow.new()
    local self = setmetatable({
        _window_details = nil,
        _closing = false,
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

    local pos = window.create_window_set_position(2, 2)
    self._window_details.dim = window.get_window_dim()
    local config = window.create_window_config(self._window_details.dim)
    vim.api.nvim_win_set_config(self._window_details.win_id, config)
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
        self._window_details = window.create_window()

        window.on_close(self._window_details, function()
            self._window_details = nil
            self._closing = false
        end)

    else
        self:close()
    end
end

return TowerOffenseWindow


