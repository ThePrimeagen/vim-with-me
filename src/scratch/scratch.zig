const std = @import("std");
const assert = @import("../assert/assert.zig").assert;

var scratch: [2 * 65536]u8 = [_]u8{0} ** (65536 * 2);
var idx: usize = 0;

pub fn toNumber(size: usize) ![]u8 {
    return try std.fmt.bufPrint(scratchBuf(12), "{}", .{size});
}

pub fn scratchBuf(size: usize) []u8 {
    assert(size < scratch.len, "you cannot require more than the size of the scratch buffer");
    if (idx + size > scratch.len) {
        idx = 0;
    }

    const out = scratch[idx .. idx + size];
    idx += out.len;
    return out;
}
