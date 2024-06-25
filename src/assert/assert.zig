const std = @import("std");

const print = std.debug.print;

pub fn unwrap(comptime T: type, val: anyerror!T) T {
    if (val) |v| {
        return v;
    } else |err| {
        std.debug.panic("unwrap error: {any}", .{err});
    }

}

pub fn assert(truthy: bool, msg: []const u8) void {
    if (!truthy) {
        @panic(msg);
    }
}
