local DisplayCache = require("vim-with-me.window.cache")
local ColorSet = require("vim-with-me.app.colors")
local window = require("vim-with-me.window")
local Commands = require("vim-with-me.app.commands")
local parse = require("vim-with-me.tcp.parse")
local tcp_utils = require("vim-with-me.tcp.utils")

---@class VWMApp
---@field public window WindowDetails | nil
---@field public cache DisplayCache | nil
---@field public color_set VWMColorSet | nil
---@field public conn TCP
---@field public commands TCPCommands
---@field public _on_render (fun(): nil) | nil
---@field public _on_command (fun(cmd: TCPCommand): nil)[]
local App = {}
App.__index = App

---@param conn TCP
---@return VWMApp
function App:new(conn)
    assert(conn:connected(), "connection not established")

    local app = setmetatable({
        conn = conn,
        _on_command = {},
        commands = Commands.Commands:new(),
        _render_cb = nil,
        window = nil,
        cache = nil,
        color_set = nil,
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
function App:on_command(cb)
    table.insert(self._on_command, cb)
    return self
end

---@param partials CellWithLocation
function App:partial_render(partials)
    for _, partial in ipairs(partials) do
        self.cache:partial(partial)
    end

    --- TODO: Create it so that i only get back partial row updates
    --- Consider some sort of debounce here too
    self.cache:render_into(self.window)

    vim.schedule(function()
        for _, partial in ipairs(partials) do
            self.color_set:color_cell(partial)
        end

        if self._on_render then
            self._on_render()
        end
    end)
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

    if DEBUG ~= nil then
        self.commands:pretty_print(command)
    end

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
    end

    for _, cmd_cb in ipairs(self._on_command) do
        cmd_cb(command)
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
    self.color_set = ColorSet:new(self.window)

    vim.api.nvim_buf_set_lines(
        self.window.buffer,
        0,
        -1,
        false,
        self.cache:to_string_rows()
    )

    return self
end

return App
