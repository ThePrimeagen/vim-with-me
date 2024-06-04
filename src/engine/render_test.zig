const std = @import("std");
const assert = @import("assert").assert;
const r = @import("primitives.zig");
const Renderable = r.Renderable;
const Rendered = r.Rendered;
const Renderer = @import("render.zig").Renderer;

const TestRenderable = struct {
    renderable: r.Renderable,

    const vtable = r.Renderable.VTable {
        .render = render,
        .deinit = deinit,
    };

    pub fn init(id: usize, z: usize) TestRenderable {
        return .{
            .renderable = Renderable.init(&vtable, id, z),
        };
    }

    const Self = @This();
    pub fn render(renderable: *r.Renderable) Rendered {
        const self: *TestRenderable = @fieldParentPtr("renderable", renderable);
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

    pub fn deinit(rend: *r.Renderable) void {
        const self: *TestRenderable = @fieldParentPtr("renderable", rend);
        _ = self;
    }

};

const testing = std.testing;
test "adding some renderables" {
    const alloc = std.testing.allocator;

    var r1 = TestRenderable.init(69, 1);
    var r2 = TestRenderable.init(70, 2);
    var r3 = TestRenderable.init(71, 1);
    var r4 = TestRenderable.init(72, 3);
    var r5 = TestRenderable.init(73, 1);

    var renderer = Renderer.init(alloc);
    defer renderer.deinit();

    try renderer.add(&r1.renderable);
    try renderer.add(&r2.renderable);
    try renderer.add(&r3.renderable);
    try renderer.add(&r4.renderable);
    try renderer.add(&r5.renderable);

    try testing.expectEqualSlices(
        @TypeOf(&r1.renderable),
        &.{&r3.renderable, &r5.renderable, &r1.renderable, &r2.renderable, &r4.renderable},
        renderer.renderables.items,
    );

    renderer.remove(&r1.renderable);
    defer r1.renderable.deinit();

    try testing.expectEqualSlices(
        @TypeOf(&r2.renderable),
        &.{&r3.renderable, &r5.renderable, &r2.renderable, &r4.renderable},
        renderer.renderables.items,
    );
}


