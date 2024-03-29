local utils = require("vim-with-me.utils")
local system = vim.system or require("vim-with-me.system")
local TCP = require("vim-with-me.tcp").TCP

local stdout = {}
local function get_stdout()
    local out = stdout
    stdout = {}
    return out
end

---@param tcp TCP
---@return (fun(): TCPCommand | nil), fun(): TCPCommand[]
local function create_tcp_next(tcp)
    local received = {}
    tcp:listen(function(command, data)
        table.insert(received, {command = command, data = data})
    end)

    local function next_cmd()
        vim.wait(1000, function()
            return #received > 0
        end)

        if #received == 0 then
            return nil
        end

        local out = table.remove(received, 1)
        return out
    end

    local function flush()
        local cmds = {}
        while true do
            local cmd = next_cmd()
            if cmd == nil then
                break
            end
            table.insert(cmds, cmd)
        end
        return cmds
    end

    return next_cmd, flush
end

---@param name string
---@param port number
---@return TCP
local function create_test_conn(name, port)
    local done_building = false
    system.run({"go", "build", "-o", name, string.format("./test/%s/main.go", name)}, {
    }, function(exit_info)
        if exit_info.code ~= 0 then
            print(exit_info.stderr or "no standard error")
            os.exit(exit_info.code, true)
        end

        done_building = true
    end)
    vim.wait(1000, function()
        return done_building
    end)

    system.run({string.format("./%s", name), "--port", tostring(port)}, {
        stdout = function(_, data)
            table.insert(stdout, data)
        end,
    })
    vim.wait(100)

    local connected = false
    local tcp = TCP:new({
        host = "127.0.0.1",
        port = port,
        retry_count = 3,
    })
    tcp:start(function()
        connected = true
    end)

    vim.wait(1000, function()
        return connected == true
    end)

    assert(connected, "could not connect to server")
    return tcp
end

-- --- TODO: This is for me to think about before i create this... I don't want to create complex testing application
-- ---@class VWMTestApp
-- ---@field public app VWMApp
-- ---@field public conn TCP
-- local TestApp = {}
-- TestApp.__index = TestApp
--
-- ---@param app VWMApp
-- ---@param conn TCP
-- ---@return VWMTestApp
-- function TestApp:new(app, conn)
--     return setmetatable({
--         app = app,
--         conn = conn,
--     }, self)
-- end
--
-- ---@param tcp TCP
-- local function create_app(tcp)
--     local app = App:new(tcp)
--
--     local latest_command = nil
--     local function on_command(cmd, data)
--
--     end
--
--     return app, function()
--     end
-- end


local function load(name)
    local file_contents = utils.read_file(name)
    if file_contents == nil then
        return nil
    end
    file_contents = string.gsub(string.gsub(file_contents, "*", " "), "\n", "")
    return file_contents
end

local theprimeagen = load("lua/vim-with-me/integration/theprimeagen")
local theprimeagen_partial = load("lua/vim-with-me/integration/theprimeagen.partial")
local empty = load("lua/vim-with-me/integration/empty")

local function before_each()
    stdout = {}
end
local function after_each()
    for _, v in ipairs(stdout) do
        print("stdout: ", v)
    end
end

return {
    create_tcp_next = create_tcp_next,
    create_test_conn = create_test_conn,
    theprimeagen = theprimeagen,
    theprimeagen_partial = theprimeagen_partial,
    empty = empty,
    get_stdout = get_stdout,
    before_each = before_each,
    after_each = after_each,
}
