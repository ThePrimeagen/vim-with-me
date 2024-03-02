local M = {}
local namespace = vim.api.nvim_create_namespace("vim-with-me.colors")

function M.color(color, x, y)
    vim.api.nvim_buf_add_highlight(0, namespace, color, x, y, -1)
end

return M
