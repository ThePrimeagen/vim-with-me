local function from_cache(row, col)
    return {
        row = row - 1,
        col = col - 1,
    }
end

---@param cell CellWithLocation
---@return number, number
local function to_cache(cell)
    return cell.loc.row + 1, cell.loc.col + 1
end

---@param cell CellWithLocation
---@return number, number
local function to_line_api(cell)
    return cell.loc.row, cell.loc.col
end

return {
    to_cache = to_cache,
    to_line_api = to_line_api,
    from_cache = from_cache,
}
