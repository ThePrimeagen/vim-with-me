local plenary = require("plenary.reload")
plenary.reload_module("vim-with-me")

local App = require("vim-with-me.app")
local TestUtils = require("vim-with-me.test-utils")
local IntUtils = require("vim-with-me.integration.int_utils")

local DATA_PATH = "./data/partial"
local TEST_SERVER = "particle_server"
local PORT = 42069
local LAUNCH_SERVER = false

---@type TCP
local tcp = nil
local function server_run()
    local server_info
    if LAUNCH_SERVER then
        server_info = IntUtils.build_go_test_server(TEST_SERVER, {debug = true})
        if not server_info.success then
            print(string.format("unable to start server: %d", server_info.exit_code))
            error(vim.inspect(server_info.stderr))
            return
        end
        IntUtils.run_test_server(server_info, PORT)
    else
        server_info = IntUtils.server_info_from_name(TEST_SERVER)
    end

    vim.wait(100)
    tcp = IntUtils.create_tcp_connection(PORT)
end

local function file_run()
    tcp = TestUtils.fake_tcp_from_file(DATA_PATH)
end

server_run()
assert(tcp ~= nil, "please call file_run or server_run")
App:new(tcp)

















