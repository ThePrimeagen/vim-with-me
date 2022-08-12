VWMClient = nil;

local M = {}
function mysplit(inputstr, sep)
    if sep == nil then
        sep = "%s"
    end
    local t = {}
    for str in string.gmatch(inputstr, "([^"..sep.."]+)") do
        table.insert(t, str)
    end
    return t
end

local function rtl()
    local function on_expire()
        vim.cmd("set norightleft")
    end

    vim.cmd("set rightleft")
    vim.defer_fn(function()
        on_expire()
    end, 5000)
end

local point_count = 0;
local point_expected = 750;
local function chat_yes_or_no(cmd_and_name)
    if point_count >= point_expected then
        return
    end

    local cmd, name = unpack(mysplit(cmd_and_name, ":"));

    if cmd == "yes" then
        point_count = point_count + 10;
    else
        point_count = point_count - 10;
    end

    if point_count < 0 then
        point_count = 0
    end

    if point_count >= point_expected then
        vim.loop.write(VWMClient, name);
        vim.defer_fn(function()
            vim.cmd(":qa!")
        end, 1000)
    end
end

function M.get_points()
    return point_count, point_expected
end

local function parse_message(chunk, offset)
    print("chunk", chunk)

    local size = string.byte(chunk, offset + 1)
    local type = string.byte(chunk, offset + 2)
    local cmd = ""

    if #chunk > 2 then
        cmd = string.sub(chunk, 3, offset + 3 + size - 1)
    end

    -- TODO: didn't do the math, just assuming i am awesome
    offset = offset + 2 + size

    return type, cmd, offset
end

local function handle_message(type, cmd)
    if type == 0 then
        print("cmd", cmd)
        vim.cmd(string.format("norm! %s", cmd))
    elseif type == 1 then
        rtl();
    elseif type == 2 then
        chat_yes_or_no(cmd)
    end
end

function M.StartVimWithMe()
    if VWMClient ~= nil then
        return
    end

    VWMClient = vim.loop.new_tcp()
    vim.loop.tcp_connect(VWMClient, "127.0.0.1", 6969, function (err)
        assert(not err, err)

        vim.loop.read_start(VWMClient, function (inner_err, chunk)
            assert(not inner_err, inner_err)

            -- disconnect
            if chunk == nil then
                M.StopVimWithMe()
                return
            end

            vim.schedule(function()
                -- do
                local type, cmd, _ = parse_message(chunk, 0)
                handle_message(type, cmd)
            end)
        end)
    end)
end

function M.StopVimWithMe()
    if VWMClient ~= nil then
        vim.loop.read_stop(VWMClient)
    end
    VWMClient = nil
end

M.ClientStopped = "Client Stopped"
M.ClientRunning = "Client Running"
function M.StatusVimWithMe()
    if VWMClient ~= nil then
        return "Client Running"
    end
    return "Client Stopped"
end

return M
