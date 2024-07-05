const std = @import("std");
const testing = @import("testing");
const objects = @import("objects");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const alloc = gpa.allocator();
    const args = try testing.params.readFromArgs(alloc);
    const values = args.values();

    _ = values;
}
