local VERSION = 1

local function parse(chunk, start, len)
    assert(type(chunk) == "string", "chunk must be a string")
    assert(start > 0, "start must be greater than 0")
    assert(len > 0, "len must be greater than 0")

    local remaining = string.sub(chunk, start + len)
    local idx = string.find(chunk, ":", start)

    local command = string.sub(chunk, start, idx - 1)
    local data = string.sub(chunk, idx + 1, start + len - 1)

    return remaining, command, data
end

local function process_packets()
    local previous_chunk = ""

    ---@param chunk string | nil
    ---@return string | nil, string | nil
    return function(chunk)
        chunk = chunk or ""

        -- split by :
        chunk = previous_chunk .. chunk
        previous_chunk = ""

        local idx = string.find(chunk, ":")
        -- what about keeping previous chunks?
        if idx == nil then
            previous_chunk = chunk
            return nil, nil
        end

        local version = tonumber(string.sub(chunk, 1, idx - 1))
        assert(version == VERSION, "version mismatch")

        local prev_idx = idx
        idx = string.find(chunk, ":", prev_idx + 1)
        -- what about keeping previous chunks?
        if idx == nil then
            previous_chunk = chunk
            return nil, nil
        end

        local len = tonumber(string.sub(chunk, prev_idx + 1, idx - 1))
        assert(type(len) == "number", "len must be a number")

        if len + idx > string.len(chunk) then
            previous_chunk = chunk
            return nil, nil
        end

        local next_chunk, command, data = parse(chunk, idx + 1, len)

        previous_chunk = next_chunk
        return command, data
    end
end

---@param command string
---@param data string
---@return string
local function create_tcp_command(command, data)
    assert(type(command) == "string", "command must be a string")
    assert(type(data) == "string", "data must be a string")

    local tcp_data = string.format("%s:%s", command, data)
    local len = string.len(tcp_data)
    local header = string.format("%d:%d:", VERSION, len)
    return string.format("%s%s", header, tcp_data)
end

return {
    parse = parse,
    process_packets = process_packets,
    create_tcp_command = create_tcp_command,
    VERSION = VERSION,
}
