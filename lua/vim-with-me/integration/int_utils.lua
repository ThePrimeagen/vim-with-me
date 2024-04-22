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

---@param exec_name string
---@param port number
local function create_test_server(exec_name, port)
    local name = string.format("/tmp/%s", exec_name)
    local done_building = false
    system.run(
        { "go", "build", "-o", name, string.format("./test/%s/main.go", exec_name) },
        {},
        function(exit_info)
            if exit_info.code ~= 0 then
                print(exit_info.stderr or "no standard error")
                os.exit(exit_info.code, true)
            end

            done_building = true
        end
    )
    vim.wait(1000, function()
        return done_building
    end)

    local run = system.run({ name, "--port", tostring(port) }, {
        stdout = function(_, data)
            print("stdout:", data)
        end,
        stderr = function(_, data)
            print("stderr:", data)
        end,

        env = {
            LEVEL = LEVEL
        }
    })
    table.insert(running, run)
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

local function before_each()
end

local function after_each()
    print("after_each?")
    vim.wait(5000)
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
}
