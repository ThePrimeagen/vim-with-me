-- local token = require("vim-with-me.app.token")
local DisplayCache = require("vim-with-me.window.cache")
local window = require("vim-with-me.window")
local Commands = require("vim-with-me.app.commands")
local parse = require("vim-with-me.tcp.parse")

---@class VWMApp
---@field public window WindowDetails | nil
---@field public cache DisplayCache | nil
---@field public conn TCP
---@field public commands TCPCommands
---@field public _auth_cb (fun(): nil) | nil
---@field public _on_render (fun(): nil) | nil
---@field public _on_command (fun(cmd: TCPCommand): nil) | nil
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
        commands = Commands.Commands:new(),
        _render_cb = nil,
        window = nil,
        cache = nil,
    }, self)

    conn:listen(function(command)
        app:_process(command)
    end)
    return app
end

---@param cb fun(): nil
---@return VWMApp
function App:on_render(cb)
    self._render_cb = cb
    return self
end

---@param cb fun(cmd: TCPCommand): nil
---@return VWMApp
function App:on_cmd_received(cb)
    self._render_cb = cb
    return self
end

---@param partials CellWithLocation
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

---@param command TCPCommand
function App:_process(command)
    local cmd = command.command
    local data = command.data
    local cmds = self.commands

    if cmd == cmds:get("commands") then
        self.commands:parse(data)
    elseif cmd == cmds:get("partial") and self.window and self.cache then
        self:partial_render(parse.parse_partial_renders(data))
    elseif cmd == cmds:get("render") and self.window and self.cache then
        self:render(data)
    elseif cmd == cmds:get("close") then
        self:close()
    elseif cmd == cmds:get("openWindow") then
        local dim = window.parse_command_data(data)
        self:with_window(dim, true)
    elseif cmd == cmds:get("error") then
        -- TODO: error and then close
        self:close()
        --[[
    elseif cmd == "auth" then
        assert(self._auth_cb, "no auth callback")
        self._auth_cb()
        self._auth_cb = nil
    end
    -]]
    end

    if self._on_command then
        self._on_command(command)
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
    self.window = window.create_window(dim, center)
    window.focus(self.window)

    self.cache = DisplayCache:new(dim)
    return self
end

---@param cb function
function App:authenticate(cb)
    assert(cb, self)
    --[[
    local t = token.get_token()
    self.conn:send("auth", t)
    self._auth_cb = cb
    --]]
end

return App
