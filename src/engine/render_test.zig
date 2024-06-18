const std = @import("std");
const assert = @import("assert").assert;
const r = @import("primitives.zig");
const Renderable = r.Renderable;
const Rendered = r.Rendered;
const Renderer = @import("render.zig").Renderer;

const TestRenderable = struct {
    _id: usize,
    _z: usize,

    const vtable = r.Renderable.VTable {
        .render = render,
        .id = id,
        .z = z,
    };

    const Self = @This();

    pub fn asRenderable(self: *Self) Renderable {
        return Renderable{
            .vtable = &vtable,
            .ptr = self,
        };
    }

    pub fn render(ptr: *anyopaque) Rendered {
        const self: *TestRenderable = @ptrCast(@alignCast(ptr));
        _ = self;

        return .{
            .loc = .{
                .row = 6,
                .col = 9,
            },
            .data = &.{},
            .cols = 2
        };
    }

    pub fn deinit(ptr: *anyopaque) void {
        const self: *TestRenderable = @ptrCast(@alignCast(ptr));
        _ = self;
    }

    pub fn id(ptr: *anyopaque) usize {
        const self: *TestRenderable = @ptrCast(@alignCast(ptr));
        return self._id;
    }

    pub fn z(ptr: *anyopaque) usize {
        const self: *TestRenderable = @ptrCast(@alignCast(ptr));
        return self._z;
    }

};

const testing = std.testing;
test "adding some renderables" {
    const alloc = std.testing.allocator;

    var t1 = TestRenderable{._id = 69, ._z = 1};
    var t2 = TestRenderable{._id = 70, ._z = 2};
    var t3 = TestRenderable{._id = 71, ._z = 1};
    var t4 = TestRenderable{._id = 72, ._z = 3};
    var t5 = TestRenderable{._id = 73, ._z = 1};

    const r1 = t1.asRenderable();
    const r2 = t2.asRenderable();
    const r3 = t3.asRenderable();
    const r4 = t4.asRenderable();
    const r5 = t5.asRenderable();

    var renderer = Renderer.init(alloc);
    defer renderer.deinit();

    try renderer.add(r1);
    try renderer.add(r2);
    try renderer.add(r3);
    try renderer.add(r4);
    try renderer.add(r5);

    try testing.expectEqualSlices(
        @TypeOf(r1),
        &.{r3, r5, r1, r2, r4},
        renderer.renderables.items,
    );

    renderer.remove(r1);

    try testing.expectEqualSlices(
        @TypeOf(r2),
        &.{r3, r5, r2, r4},
        renderer.renderables.items,
    );
}


