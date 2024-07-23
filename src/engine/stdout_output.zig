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


pub fn hideCursor() void {
    _ = std.io.getStdOut().write("\x1b[?25l") catch |e| {
        std.debug.print("error while clearing stdout: {}\n", .{e});
    };
}

pub fn showCursor() void {
    _ = std.io.getStdOut().write("\x1b[?25h") catch |e| {
        std.debug.print("error while clearing stdout: {}\n", .{e});
    };
}

const os = std.posix;
pub fn showCursorOnSigInt() !void {
    const internal_handler = struct {
        fn internal_handler(sig: c_int) callconv(.C) void {
            _ = sig;
            showCursor();
            std.process.exit(0);
        }
    }.internal_handler;
    const act = os.Sigaction{
        .handler = .{ .handler = internal_handler },
        .mask = os.empty_sigset,
        .flags = 0,
    };
    try os.sigaction(os.SIG.INT, &act, null);
}
