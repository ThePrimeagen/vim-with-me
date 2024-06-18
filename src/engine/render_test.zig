const std = @import("std");
const assert = @import("assert").assert;
const r = @import("renderable.zig");
const Renderable = r.Renderable;
const Rendered = r.Rendered;
const Renderer = @import("render.zig").Renderer;

const TestRenderable = struct {

    _id: usize,
    _z: usize,
    data: []const u8,
    alloc: std.mem.Allocator,

    const Self = @This();
    pub fn render(self: *Self) Rendered {
        return .{
            .loc = .{
                .row = 6,
                .col = 9,
            },
            .data = self.data,
            .cols = 2
        };
    }

    pub fn id(self: *Self) usize {
        return self._id;
    }

    pub fn z(self: *Self) usize {
        return self._z;
    }

    pub fn deinit(self: *Self) void {
        self.alloc.destroy(self);
    }

};

fn createTestRenderable(alloc: std.mem.Allocator, id: usize, z: usize) !Renderable {
    const renderable = try alloc.create(TestRenderable);
    renderable.* = .{
        ._id = id,
        ._z = z,
        .data = &.{},
        .alloc = alloc,
    };

    return Renderable.init(renderable);
}

const testing = std.testing;
test "adding some renderables" {
    const alloc = std.testing.allocator;

    var r1 = try createTestRenderable(alloc, 69, 1);
    var r2 = try createTestRenderable(alloc, 70, 2);
    var r3 = try createTestRenderable(alloc, 71, 1);
    var r4 = try createTestRenderable(alloc, 72, 3);
    var r5 = try createTestRenderable(alloc, 73, 1);

    var renderer = Renderer.init(alloc);
    defer renderer.deinit();

    try renderer.add(&r1);
    try renderer.add(&r2);
    try renderer.add(&r3);
    try renderer.add(&r4);
    try renderer.add(&r5);

    try testing.expectEqualSlices(
        @TypeOf(&r1),
        &.{&r3, &r5, &r1, &r2, &r4},
        renderer.renderables.items,
    );

    renderer.remove(&r1);
    defer r1.deinit();

    try testing.expectEqualSlices(
        @TypeOf(&r2),
        &.{&r3, &r5, &r2, &r4},
        renderer.renderables.items,
    );
}


