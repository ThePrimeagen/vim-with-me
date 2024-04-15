---@class WindowPosition
---@field width number
---@field height number
---@field row number
---@field col number

---@class WindowDetails
---@field dim WindowPosition
---@field buffer number
---@field win_id number

local group = vim.api.nvim_create_augroup("vim-with-me.window", {
    clear = true,
})

local M = {}

---@param width number
---@param height number
---@return WindowPosition
function M.create_window_dimensions(width, height)
    return {
        width = width,
        height = height,
        row = 0,
        col = 0,
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

---@param dim WindowPosition
function M.center_window_dimension(dim)
    dim = dim or M.create_window_dimensions(80, 24)
    local ui = vim.api.nvim_list_uis()[1]

    local x_offset = 80
    local y_offset = 24

    if ui ~= nil then
        local w_diff = math.floor((ui.width - dim.width) / 2)
        local h_diff = math.floor((ui.height - dim.height) / 2)

        x_offset = math.max(w_diff, 0)
        y_offset = math.max(h_diff, 0)
    end

    dim.row = y_offset
    dim.col = x_offset
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

---@param dim WindowPosition | nil
---@param center boolean | nil defaults true
---@return WindowDetails
function M.create_window(dim, center)
    center = center == nil or center
    dim = dim or M.create_window_dimensions(80, 24)

    if center then
        M.center_window_dimension(dim)
    end

    local buffer = vim.api.nvim_create_buf(false, true)
    local config = M.create_window_config(dim)
    local win_id = vim.api.nvim_open_win(buffer, false, config)

    return {
        dim = dim,
        buffer = buffer,
        win_id = win_id,
    }
end

---@param details WindowDetails
---@return boolean
local function clear_if_invalid(details)
    if not vim.api.nvim_win_is_valid(details.win_id) then
        vim.api.nvim_clear_autocmds({
            group = group,
        })
        return true
    end
    return false
end

--- THIS IS SLIGHTLY INCORRECT
---@param details WindowDetails
function M.resize(details)
    if clear_if_invalid(details) then
        return
    end

    M.center_window_dimension(details.dim)
    local config = M.create_window_config(details.dim)
    vim.api.nvim_win_set_config(details.win_id, config)
end

---@param details WindowDetails
local function validate_details(details)
    assert(vim.api.nvim_win_is_valid(details.win_id), "window must be valid")
    assert(vim.api.nvim_buf_is_valid(details.buffer), "buffer must be valid")
end

---@param details WindowDetails
---@param cb function
function M.on_close(details, cb)
    assert(cb ~= nil, "callback must be provided")
    validate_details(details)

    vim.api.nvim_create_autocmd("BufUnload", {
        group = group,
        buffer = details.buffer,
        callback = function()
            if clear_if_invalid(details) then
                return
            end
            cb()
        end,
    })
end

---@param details WindowDetails
function M.focus(details)
    validate_details(details)
    vim.api.nvim_set_current_win(details.win_id)
end

---@param details WindowDetails
function M.refocus(details)
    validate_details(details)

    vim.api.nvim_create_autocmd("BufEnter", {
        group = group,
        callback = function()
            if clear_if_invalid(details) then
                return
            end
            vim.api.nvim_set_current_win(details.win_id)
        end,
    })
end

---@param data string
---@return WindowPosition
function M.parse_command_data(data)
    assert(#data == 2, "window position command should have length 2")
    local open = {
        row = 0,
        col = 0,
        height = string.byte(data, 1, 1),
        width = string.byte(data, 2, 2),
    }

    return open
end

---@param str string
---@param idx number | nil
---@return number | nil, number
local function next_number(str, idx)
    assert(type(str) == "string", "next requires str to be a string")
    if idx == nil then
        return nil, #str
    end
    local next_idx = string.find(str, ":", idx)
    if not next_idx then
        return nil, #str
    end
    local num_str = string.sub(str, idx, next_idx - 1)
    local num = tonumber(num_str)
    assert(num ~= nil, "parsed item was not a string: " .. num_str)
    return num, next_idx + 1
end

return M
