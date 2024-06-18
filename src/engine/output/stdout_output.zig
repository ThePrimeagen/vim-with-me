const assert = @import("assert").assert;

const std = @import("std");
const io = std.io;

const ColorOutput = @import("framer.zig").Cell;

pub fn output(out: []const u8) !void {
    _ = try io.getStdOut().write(out);
}

