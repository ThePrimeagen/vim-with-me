--- PROBABLY NOT A THING TO DO--- yet
---@class VWMApp
---@field public window WindowDetails | nil
---@field public conn TCP
local App = {}
App.__index = App

---@param conn TCP
---@return VWMApp
function App:new(conn)
    return setmetatable({
        conn = conn
    }, self)
end

---@param window WindowDetails
---@return VWMApp
function App:withWindow(window)
    self.window = window
    return self
end

---@param cb function
function App:authenticate(cb)
end

return App
