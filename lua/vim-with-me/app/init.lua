local token = require("vim-with-me.app.token")
local DisplayCache = require("vim-with-me.window.cache")
local window = require("vim-with-me.window")

--- PROBABLY NOT A THING TO DO--- yet
---@class VWMApp
---@field public window WindowDetails | nil
---@field public cache DisplayCache | nil
---@field public conn TCP
---@field public _auth_cb (fun(): nil) | nil
---@field public _unhandled_commands fun(cmd: string, data: string): nil
local App = {}
App.__index = App

---@param conn TCP
---@param unhandled_commands fun(cmd: string, data: string): nil
---@return VWMApp
function App:new(conn, unhandled_commands)
    assert(conn:connected(), "connection not established")
    assert(unhandled_commands, "no unhandled commands function")

    local app = setmetatable({
        conn = conn,
        _unhandled_commands = unhandled_commands,
        _auth_cb = nil,
        window = nil,
        cache = nil,
    }, self)

    conn:listen(function(command, data) app:_process(command, data) end)
    return app
end

local function print_command(command, data)
    print("command", command, "data", data)
end

function App:_process(command, data)

    print_command(command, data)
    if command == "p" and self.window and self.cache then
        local partial = window.parse_partial_render(data)
        print("partial", vim.inspect(partial))
        self.cache:partial(partial)

        --- TODO: Create it so that i only get back partial row updates
        --- Consider some sort of debounce here too
        local rows = self.cache:to_string_rows()
        vim.api.nvim_buf_set_lines(self.window.buffer, 0, -1, false, rows)
    elseif command == "r" and self.window and self.cache then
        -- check to see if last character is a new line
        if string.sub(data, -1) == "\n" then
            data = string.sub(data, 1, -2)
        end
        self.cache:from_string(data)

        local rows = self.cache:to_string_rows()
        vim.api.nvim_buf_set_lines(self.window.buffer, 0, -1, false, rows)
    elseif command == "c" then
        self:close()
    elseif command == "open-window" then
        local dim = window.parse_command_data(data)
        print("open-window", vim.inspect(data), vim.inspect(dim))
        self:with_window(dim.width, dim.height, true)
    elseif command == "e" then
        -- TODO: error and then close
        self:close()
    elseif command == "auth" then
        assert(self._auth_cb, "no auth callback")
        self._auth_cb()
        self._auth_cb = nil
    else
        assert(self._unhandled_commands, "no unhandled commands function")
        self._unhandled_commands(command, data)
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

---@param width number
---@param height number
---@param center boolean | nil
---@return VWMApp
function App:with_window(width, height, center)

    print("with_window", vim.inspect(self.window), width, height, center)
    assert(self.window == nil, "window already open")

    center = center or center == nil
    self.window = window.create_window(window.create_window_dimensions(80, 24), center);
    window.focus(self.window)

    self.cache = DisplayCache:new(width, height)
    return self
end

---@param cb function
function App:authenticate(cb)
    local t = token.get_token()
    self.conn:send("auth", t)
    self._auth_cb = cb
end

return App
