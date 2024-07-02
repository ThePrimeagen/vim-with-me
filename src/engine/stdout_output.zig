const assert = @import("../assert/assert.zig").assert;

const std = @import("std");
const io = std.io;

const ColorOutput = @import("framer.zig").Cell;

pub fn output(out: []const u8) !void {
    _ = try io.getStdOut().write(out);
}

pub fn resetColor() void {
    std.debug.print("\x1b[0m", .{});
    _ = std.io.getStdOut().write("\x1b[0m") catch |e| {
        std.debug.print("error while clearing stdout: {}\n", .{e});
    };
}

