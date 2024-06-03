const std = @import("std");

const print = std.debug.print;

pub fn assert(truthy: bool, msg: []const u8) void {
    if (!truthy) {
        @panic(msg);
    }
}
