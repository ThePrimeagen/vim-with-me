local token = require("vim-with-me.app.token")
local DisplayCache = require("vim-with-me.window.cache")
local window = require("vim-with-me.window")

---@class VWMApp
---@field public window WindowDetails | nil
---@field public cache DisplayCache | nil
---@field public conn TCP
---@field public _auth_cb (fun(): nil) | nil
---@field public _on_render (fun(): nil) | nil
---@field public _on_command (fun(cmd: string, data: string): nil) | nil
local App = {}
App.__index = App

---@param conn TCP
---@return VWMApp
function App:new(conn)
    assert(conn:connected(), "connection not established")

    local app = setmetatable({
        conn = conn,
        _on_command = nil,
        _auth_cb = nil,
        _render_cb = nil,
        window = nil,
        cache = nil,
    }, self)

    conn:listen(function(command, data) app:_process(command, data) end)
    return app
end

---@param cb fun(): nil
---@return VWMApp
function App:on_render(cb)
    self._render_cb = cb
    return self
end

---@param cb fun(cmd: string, data: string): nil
---@return VWMApp
function App:on_cmd_received(cb)
    self._render_cb = cb
    return self
end

---@param partials PartialRender
function App:partial_render(partials)
    for _, partial in ipairs(partials) do
        self.cache:partial(partial)
    end

    --- TODO: Create it so that i only get back partial row updates
    --- Consider some sort of debounce here too
    local rows = self.cache:to_string_rows()
    vim.api.nvim_buf_set_lines(self.window.buffer, 0, -1, false, rows)

    if self._on_render then
        self._on_render()
    end
end

---@param str string
function App:render(str)
    -- check to see if last character is a new line
    if string.sub(str, -1) == "\n" then
        str = string.sub(str, 1, -2)
    end
    self.cache:from_string(str)

    local rows = self.cache:to_string_rows()
    vim.api.nvim_buf_set_lines(self.window.buffer, 0, -1, false, rows)

    if self._on_render then
        self._on_render()
    end
end

function App:_process(command, data)
    if command == "p" and self.window and self.cache then
        self:partial_render(window.parse_partial_render(data))
    elseif command == "r" and self.window and self.cache then
        self:render(data)
    elseif command == "c" then
        self:close()
    elseif command == "open-window" then
        local dim = window.parse_command_data(data)
        self:with_window(dim, true)
    elseif command == "e" then
        -- TODO: error and then close
        self:close()
    elseif command == "auth" then
        assert(self._auth_cb, "no auth callback")
        self._auth_cb()
        self._auth_cb = nil
    end

    if self._on_command then
        self._on_command(command, data)
    end

end

function App:close()
    self.conn:close()
    self.conn = nil

    if self.window then
        window.close_window(self.window)
        self.window = nil
        self.cache = nil
    end
end

---@param dim WindowPosition
---@param center boolean | nil
---@return VWMApp
function App:with_window(dim, center)

    assert(self.window == nil, "window already open")

    center = center or center == nil
    self.window = window.create_window(dim, center);
    window.focus(self.window)

    self.cache = DisplayCache:new(dim)
    return self
end

---@param cb function
function App:authenticate(cb)
    local t = token.get_token()
    self.conn:send("auth", t)
    self._auth_cb = cb
end

return App
