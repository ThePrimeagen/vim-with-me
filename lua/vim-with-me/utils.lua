local namespace = vim.api.nvim_create_namespace("VWM")

local function key(k)
    k = vim.api.nvim_replace_termcodes(k, true, false, true)
    vim.api.nvim_feedkeys(k, "t", false)
end

---@param path string
---@return string | nil
local function read_file(path)
    local ok, fh = pcall(vim.loop.fs_open, path, "r", 493)
    if not ok then
        return nil
    end

    local data = ""
    while true do
        local ok1, chunk = pcall(vim.loop.fs_read, fh, 1024)
        if not ok1 then
            break
        end

        if chunk == "" then
            break
        end

        data = data .. chunk
    end

    vim.loop.fs_close(fh)

    return data
end

return {
    key = key,
    read_file = read_file,
    namespace = namespace,
    optsify = function(opts, default)
        return vim.tbl_extend("force", {}, default, opts or {})
    end,
}
