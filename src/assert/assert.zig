const std = @import("std");

const print = std.debug.print;

pub fn unwrap(comptime T: type, val: anyerror!T) T {
    if (val) |v| {
        return v;
    } else |err| {
        std.debug.panic("unwrap error: {any}", .{err});
    }

}

pub fn u(v: anyerror![]u8) []u8 {
    return unwrap([]u8, v);
}

pub fn assert(truthy: bool, msg: []const u8) void {
    if (!truthy) {
        @panic(msg);
    }
}

// TODO: DO SOMETHING WITH THIS...
pub fn printZZZ(toPrint: anytype) ![]u8 {
    const MyType = @TypeOf(toPrint);
    const hasStr = @hasDecl(MyType, "string");

    if (hasStr) {
        return toPrint.string();
    }
    return .{};
}

