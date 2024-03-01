local uv = vim.loop

local function key(k)
    k = vim.api.nvim_replace_termcodes(k, true, false, true)
    vim.api.nvim_feedkeys(k, "t", false)
end

local function read(client)
    uv.read_start(client, vim.schedule_wrap(function(e, chunk)
        if chunk == "<dot>" then
            chunk = "."
        end
        key(chunk)
    end))
end

local client = nil
function START()
    client = uv.new_tcp()
    client:connect("127.0.0.1", 42069, function(err)
        read(client)
    end)
end

function CLOSE()
    if client == nil then
        return
    end

    client:shutdown()
    client:close()
end
