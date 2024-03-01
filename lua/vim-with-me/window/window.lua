---@class WindowPosition
---@field width number
---@field height number
---@field row number
---@field col number

---@class WindowDetails
---@field dim WindowPosition
---@field buffer number
---@field win_id number

local M = {}

---@param width number
---@param height number
---@return WindowPosition
function M.create_window_offset(width, height)
    return {
        width = width,
        height = height,
        row = math.floor(height / 2),
        col = math.floor(width / 2),
    }
end

---@param details WindowDetails
function M.close_window(details)
    local win_id = details.win_id
    local buffer = details.buffer

    if win_id ~= nil and vim.api.nvim_win_is_valid(win_id) then
        vim.api.nvim_win_close(win_id, true)
    end

    if buffer ~= nil and vim.api.nvim_buf_is_valid(buffer) then
        vim.api.nvim_buf_delete(buffer, { force = true })
    end
end

---@param offset WindowPosition
---@return WindowPosition
function M.get_window_dim(offset)
    offset = offset or M.create_window_offset(2, 2)
    local ui = vim.api.nvim_list_uis()[1]
    local width = 40
    local height = 20
    if ui ~= nil then
        width = math.max(ui.width, 0)
        height = math.max(ui.height, 0)
    end

    return {
        width = math.max(width - offset.width, 0),
        height = math.max(height - offset.height, 0),
        row = offset.row,
        col = offset.col,
    }
end

---@param pos WindowPosition
function M.create_window_config(pos)
    return {
        relative = "editor",
        anchor = "NW",
        row = pos.row,
        col = pos.col,
        width = pos.width,
        height = pos.height,
        border = "none",
        title = "",
        style = "minimal",
    }
end

---@param pos WindowPosition | nil
---@return WindowDetails
function M.create_window(pos)
    pos = pos or M.create_window_offset(2, 2)
    local buffer = vim.api.nvim_create_buf(false, true)
    local dim = M.get_window_dim(pos)

    local config = M.create_window_config(dim)
    local win_id = vim.api.nvim_open_win(buffer, false, config)

    return {
        dim = dim,
        buffer = buffer,
        win_id = win_id,
    }
end

---@param details WindowDetails
---@param cb function
function M.on_close(details, cb)
    vim.api.nvim_create_autocmd("BufUnload", {
        group = M.vim_apm_group_id(),
        buffer = details.buffer,
        callback = function()
            cb()
        end,
    })
end

return M
