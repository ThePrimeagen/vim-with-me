const assert = @import("../assert/assert.zig").assert;
const print = @import("std").debug.print;

pub fn xor(a: []const u8, b: []const u8, out: []u8) void {
    assert(a.len == b.len, "a and b must be the same length");
    assert(out.len >= a.len, "out must be at least as big as a");

    for (0..a.len) |i| {
        out[i] = a[i] ^ b[i];
    }
}

