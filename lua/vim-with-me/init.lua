local uv = vim.loop

local function read(client)
    uv.read_start(client, vim.schedule_wrap(function(e, chunk)

iiffuucckkkyPmriimmeeaageennngeggsgfhf<esc><esc>?xx:q!<cr>flocal client = uv.new_tcp()
client:connect("127.0.0.1", 42069, function(err)
    read(client)
end)

function CLOSE()
    client:shutdown()
    client:close()
end
