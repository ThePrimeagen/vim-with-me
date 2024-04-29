local process = require("vim-with-me.tcp.process")

---@class TCPReplayer
---@field processor fun(chunk: string?): TCPCommand | nil
---@field data string
---@field opts {replay_speed: number}
local TCPReplayer = {}
TCPReplayer.__index = TCPReplayer

---@param tcp_data string
---@param opts {replay_speed: number}?
---@return TCP
function TCPReplayer:new(tcp_data, opts)
    opts = vim.tbl_extend("force", {}, opts or {
        replay_speed = 500,
    })

    local item = setmetatable({
        processor = process.process_packets(),
        data = tcp_data,
        opts = opts,
    }, self)

    --- TODO: Add the other methods
    return item
end

function TCPReplayer:connected()
    return true
end

function TCPReplayer:listen(cb)
    local packet = self.processor(self.data)
    local function read_one()
        if packet == nil then
            return
        end

        cb(packet)
        packet = self.processor()

        vim.defer_fn(function()
            read_one()
        end, self.opts.replay_speed)
    end

    read_one()
end

---@param file_path string
---@return TCP
local function fake_tcp_from_file(file_path)
    local ok, fh = pcall(vim.loop.fs_open, file_path, "r", 493)
    if not ok then
        error("cannot open file")
    end

    local ok, data = pcall(vim.loop.fs_read, fh, 2048)
    if not ok then
        error("cannot read data")
    end

    vim.loop.fs_close(fh)

    return TCPReplayer:new(data)
end

return {
    fake_tcp_from_file = fake_tcp_from_file,
    FakeTCP = TCPReplayer,
}
