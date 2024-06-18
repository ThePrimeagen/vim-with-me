const std = @import("std");
const print = std.debug.print;

const encoding = @import("encoding/encoding.zig");
//const render = @import("render/render.zig");

pub fn main() !void {
    print("goodbye world", .{});
}

test { _ = encoding; }
//test { _ = render; }

