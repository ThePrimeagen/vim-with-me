local utils = require("vim-with-me.tcp.utils")

local VERSION = 1
local HEADER_LENGTH = 4

---@param str string
---@param index number?
---@return number
local function parse_big_endian_16(str, index)
    index = index or 1
    assert(#str >= index + 1, "string is not long enough")

    local msb = str.byte(str, index, index)
    local lsb = str.byte(str, index + 1, index + 1)

    return bit.lshift(msb, 8) + lsb
end

---@param number number
---@return string
local function to_big_endian_16(number)
    local msb = bit.rshift(number, 8)
    local lsb = number % 256

    return table.concat({
        string.char(msb),
        string.char(lsb),
    }, "")
end

local function debug_packet(chunk)
    for i = 1, #chunk do
        print("chunk:", i, string.byte(chunk, i))
    end

    if #chunk < HEADER_LENGTH then
        print("chunk is smaller than header length")
        return
    end

    local version = string.byte(chunk, 1)
    local command = string.byte(chunk, 2)
    local length = parse_big_endian_16(chunk, 3)

    print("  version", version)
    print("  command", command)
    print("  length", length)
    if #chunk >= HEADER_LENGTH + length then
        print(
            "  data",
            string.sub(chunk, HEADER_LENGTH, HEADER_LENGTH + length)
        )
        print("  extra", #chunk - (HEADER_LENGTH + length))
    else
        print("  packet doesn't have all the data yet")
    end
end

-- Version : Command : Length : Data
-- 1 byte  : 1 byte  : 2 bytes: Length bytes
---@param str string
---@param index number?
---@return TCPCommand, number
local function parse_tcp_command(str, index)
    index = index or 1

    assert(string.byte(str, index) == VERSION, "version mismatch")

    local length = parse_big_endian_16(str, 3)
    local data_start = index + HEADER_LENGTH

    assert(
        #str >= HEADER_LENGTH + length,
        "string is too short to parse tcp command from"
    )

    return {
        command = string.byte(str, index + 1),
        data = string.sub(str, data_start, data_start + length - 1),
    },
        length
end

local function process_packets()
    local previous_chunk = ""

    ---@param chunk string | nil
    ---@return TCPCommand?
    return function(chunk)
        chunk = chunk or ""

        -- split by :
        chunk = previous_chunk .. chunk

        -- Version : Command : Length : Data
        -- 1 byte  : 1 byte  : 2 bytes: Length bytes
        if #chunk < HEADER_LENGTH then
            previous_chunk = chunk
            return nil
        end

        local len = parse_big_endian_16(chunk, 3)
        if #chunk < HEADER_LENGTH + len then
            previous_chunk = chunk
            return nil
        end

        local command, n = parse_tcp_command(chunk)
        previous_chunk = string.sub(chunk, HEADER_LENGTH + n + 1)

        return command
    end
end

---@param command TCPCommand
---@return string
local function encode_tcp_command(command)
    return utils.to_string(
        VERSION,
        command.command,
        to_big_endian_16(#command.data),
        command.data
    )
end

return {
    process_packets = process_packets,
    parse_big_endian_16 = parse_big_endian_16,
    to_big_endian_16 = to_big_endian_16,
    encode_tcp_command = encode_tcp_command,
    VERSION = VERSION,
}
