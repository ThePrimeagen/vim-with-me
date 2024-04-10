--- @diagnostic disable: redefined-local
local path = vim.fs.normalize(vim.fn.stdpath("data") .. "/vim-with-me/token")

---@class AuthDetails
---@field token string
---@field twitch_name string

---@param read_path string
---@return AuthDetails | nil
local function read_file(read_path)
    local ok, fh = pcall(vim.loop.fs_open, read_path, "r", 493)
    if not ok then
        print("cannot open file")
        return nil
    end

    local ok, data = pcall(vim.loop.fs_read, fh, 1024)
    if not ok then
        return nil
    end

    vim.loop.fs_close(fh)

    local ok, json = pcall(vim.fn.json_decode, data)
    if not ok then
        return nil
    end

    return json
end

---@param write_path string
---@param data AuthDetails
local function write_file(write_path, data)
    local ok, json = pcall(vim.fn.json_encode, data)
    assert(ok, "failed to encode json")

    local dirname = vim.fs.dirname(write_path)
    ok, _ = pcall(vim.loop.fs_stat, dirname, 493)
    if not ok then
        ok, _ = pcall(vim.loop.fs_mkdir, dirname, 493)
        assert(ok, "failed to create directory")
    end

    local ok, fh = pcall(vim.loop.fs_open, write_path, "w", 493)
    assert(ok, "failed to open file")

    ok, _ = pcall(vim.loop.fs_write, fh, json)
    assert(ok, "failed to write to file")

    vim.loop.fs_close(fh)
end

---@return AuthDetails
local function get_token()
    return read_file(path)
        or {
            token = "",
            twitch_name = "",
        }
end

---@param token AuthDetails
local function set_token(token)
    write_file(path, token)
end

return {
    get_token = get_token,
    set_token = set_token,
}
