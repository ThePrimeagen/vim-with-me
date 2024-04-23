local plenary = require("plenary.reload")
plenary.reload_module("vim-with-me")

local App = require("vim-with-me.app")
local ColorSet = require("vim-with-me.app.colors")
local window = require("vim-with-me.window")
local TestUtils = require("vim-with-me.test-utils")
local IntUtils = require("vim-with-me.integration.int_utils")

local DATA_PATH = "./data/partial"
local TEST_SERVER = "color_server"
local PORT = 42069

---@type TCP
local tcp = nil
local function server_run()
    local server_info = IntUtils.build_go_test_server(TEST_SERVER)
    if not server_info.success then
        error(string.format("unable to start server: %d", server_info.exit_code))
        return
    end

    IntUtils.run_test_server(server_info, PORT)
    vim.wait(100)

    tcp = IntUtils.create_tcp_connection(PORT)
end

local function file_run()
    tcp = TestUtils.fake_tcp_from_file(DATA_PATH)
end

file_run()
assert(tcp ~= nil, "please call file_run or server_run")
App:new(tcp)

