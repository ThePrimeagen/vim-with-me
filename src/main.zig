const std = @import("std");
const print = std.debug.print;

const render = @import("engine/render.zig");
const engine = @import("engine/engine.zig");
const encoding = @import("encoding/encoding.zig");

pub fn main() !void {
    print("goodbye world", .{});
}

test { _ = encoding; }
test { _ = render; }
test { _ = engine; }

