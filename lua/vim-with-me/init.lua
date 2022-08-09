VWMClient = nil;

local M = {}

local function rtl()
    local function on_expire()
        vim.cmd("set norightleft")
    end

    vim.cmd("set rightleft")
    vim.defer_fn(function()
        on_expire()
    end, 5000)
end

local function parse_message(chunk)
    print("chunk", chunk)

    local size = string.byte(chunk, 1)
    local type = string.byte(chunk, 2)
    local cmd = ""

    if #chunk > 2 then
        cmd = string.sub(chunk, 3, 3 + size - 1)
    end

    return type, cmd
end

local function handle_message(type, cmd)
    if type == 0 then
        print("cmd", cmd)
        vim.cmd(string.format("norm! %s", cmd))
    elseif type == 1 then
        rtl();
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

            vim.schedule(function()
                handle_message(parse_message(chunk))
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

function M.StatusVimWithMe()
    if VWMClient ~= nil then
        return "Client Running"
    end
    return "Client Stopped"
end

return M
