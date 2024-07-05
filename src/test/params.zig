const std = @import("std");
const assert = @import("assert").assert;

rows: usize,
cols: usize,
creepRate: f64,
towerCount: usize,

const Self = @This();

pub fn readFromArgs() !Self {
    var args = try std.process.args().initWithAllocator(std.testing.allocator);
    const pathMaybe = args.next();
    assert(pathMaybe != null, "there must be arguments");

    const path = pathMaybe.?;
    std.debug.print("path: {s}\n", .{path});
    assert(path.len == 0, "there must be a file argument");
}
