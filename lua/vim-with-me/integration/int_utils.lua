local utils = require("vim-with-me.utils")
local system = vim.system or require("vim-with-me.system")
local TCP = require("vim-with-me.tcp").TCP

---@type vim.SystemObj[]
local running = {}

local LEVEL = vim.fn.environ().LEVEL

---@param tcps TestTCP[]
---@return (TCPCommand | nil)[]
local function read_all(tcps)
    local out = {}
    for _, tcp in ipairs(tcps) do
        table.insert(out, tcp.next())
    end
    return out
end

---@param tcp TCP
---@return (fun(): TCPCommand | nil), fun(): TCPCommand[]
local function create_tcp_next(tcp)
    local received = {}
    tcp:listen(function(command)
        table.insert(received, command)
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

local function create_tcp_connection(port)
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

---@class ServerInfo
---@field success boolean
---@field exit_code number
---@field exec string
---@field cmd string[]
---@field server_path string
---@field stdout string[]
---@field stderr string[]

---@param server_name string
---@return ServerInfo
local function server_info_from_name(server_name)
    local exec = string.format("/tmp/%s", server_name)
    local server_path = string.format("./test/%s/main.go", server_name)
    return {
        cmd = { "go", "build", "-o", exec, server_path },
        success = nil,
        exit_code = 0,
        exec = exec,
        server_path = server_path,
        stdout = {},
        stderr = {},
    }
end

---@param server_name string
---@param opts {timeout: number}?
---@return ServerInfo
local function build_go_test_server(server_name, opts)
    opts = vim.tbl_extend("force", {}, {
        debug = false,
        timeout = 1000,
    }, opts or {})

    local server_info = server_info_from_name(server_name)

    if opts.debug then
        print("executing: ", table.concat(server_info.cmd, " "))
    end

    system.run(server_info.cmd, {
        stdout = function(_, data)
            table.insert(server_info.stdout, data)
        end,
        stderr = function(_, data)
            table.insert(server_info.stderr, data)
        end,
    }, function(exit_info)
        server_info.exit_code = exit_info.code
        server_info.success = exit_info.code == 0
    end)

    vim.wait(opts.timeout, function()
        return server_info.success ~= nil
    end)

    return server_info
end

---@param server_info ServerInfo
---@param port number | string
local function run_test_server(server_info, port)
    local run = system.run({ server_info.exec, "--port", tostring(port) }, {
        stdout = function(_, data)
            print("stdout:", data)
        end,
        stderr = function(_, data)
            print("stderr:", data)
        end,

        env = {
            LEVEL = LEVEL,
        },
    })
    table.insert(running, run)
end

---@param exec_name string
---@param port number | string
---@param opts table<string, any>?
local function create_test_server(exec_name, port, opts)
    local server_info = build_go_test_server(exec_name, opts)
    if not server_info.success then
        print("failed to launch server", table.concat(server_info.stderr, "\n"))
        os.exit(server_info.exit_code)
    end

    run_test_server(server_info, port)
    vim.wait(100)
end

local function load(name)
    local file_contents = utils.read_file(name)
    if file_contents == nil then
        return nil
    end
    file_contents = string.gsub(string.gsub(file_contents, "*", " "), "\n", "")
    return file_contents
end

local theprimeagen = load("lua/vim-with-me/integration/theprimeagen")
local theprimeagen_partial =
    load("lua/vim-with-me/integration/theprimeagen.partial")
local empty = load("lua/vim-with-me/integration/empty")

local function before_each() end

local function after_each()
    for _, proc in ipairs(running) do
        proc:kill(9)
    end
    running = {}
end

return {
    create_tcp_connection = create_tcp_connection,
    create_tcp_next = create_tcp_next,
    create_test_server = create_test_server,
    theprimeagen = theprimeagen,
    theprimeagen_partial = theprimeagen_partial,
    empty = empty,
    before_each = before_each,
    after_each = after_each,
    read_all = read_all,
    build_go_test_server = build_go_test_server,
    run_test_server = run_test_server,
    server_info_from_name = server_info_from_name,
}
