local eq = assert.are.same
local int_utils = require("vim-with-me.integration.int_utils")
local utils = require("vim-with-me.utils")
local App = require("vim-with-me.app")

local function replace(str, old, new)
    return string.gsub(str, old, new)
end

local function load(name)
    local file_contents = utils.read_file(name)
    file_contents = replace(replace(file_contents, "*", " "), "\n", "")
    return file_contents
end

local theprimeagen = load("lua/vim-with-me/integration/theprimeagen")
local stdout = {}
local function stdout_cb(_, data)
    table.insert(stdout, data)
end

describe("vim with me", function()
    it("full command set", function()
        local tcp = int_utils.create_test_conn("cmd_server", 42070, stdout_cb)
        local next_cmd = int_utils.create_tcp_next(tcp)

        tcp:send("open", "")
        eq({
            command = "open-window",
            data = "24:80",
        }, next_cmd())

        tcp:send("render", "")
        local cmd = next_cmd()

        eq({
            command = "r",
            data = theprimeagen,
        }, cmd)
        tcp:send("partial", "1:1")

        local cmds = {}
        while true do
            cmd = next_cmd()
            if cmd == nil then
                break
            end
            table.insert(cmds, cmd)
        end

        local expected = {}
        local theprimeagen_str = "theprimeagen"
        for i = 1, #theprimeagen_str do
            table.insert(expected, {
                command = "p",
                data = string.format("1:%d:%s", i, theprimeagen_str:sub(i, i)),
            })
        end

        eq(expected, cmds)
    end)

    after_each(function()
        for _, v in ipairs(stdout) do
            print("stdout: ", v)
        end
    end)
end)

