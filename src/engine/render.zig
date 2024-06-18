const std = @import("std");
const print = std.debug.print;

const assert = @import("assert").assert;
const p = @import("primitives.zig");
const Renderable = p.Renderable;
const Rendered = p.Rendered;

comptime {
    _ = @import("render_test.zig");
}

const RenderList = std.ArrayList(*Renderable);

pub const Renderer = struct {
    renderables: RenderList,

    pub fn debug(self: *Renderer) void {
        print("Renderer({})\n", .{self.renderables.items.len});
        for (self.renderables.items) |renderable| {
            print("    renderable({}, {}): __add debug method__\n", .{renderable.id, renderable.z});
        }
    }

    pub fn init(alloc: std.mem.Allocator) Renderer {
        return .{ .renderables = RenderList.init(alloc) };
    }

    pub fn deinit(self: *Renderer) void {
        for (self.renderables.items) |renderable| {
            renderable.deinit();
        }

        self.renderables.deinit();
    }

    pub fn add(self: *Renderer, renderable: *Renderable) !void {
        var lo: usize = 0;
        var hi: usize = self.renderables.items.len;

        const needle = renderable.z;
        var idx: i32 = -1;

        while (lo < hi) {
            const midpoint = lo + (hi - lo) / 2;
            const value = self.renderables.items[midpoint].z;

            if (value == needle) {
                idx = @intCast(midpoint);
                break;
            } else if (value > needle) {
                hi = midpoint;
            } else if (value < needle) {
                lo = midpoint + 1;
            }
        }

        if (idx == -1) {
            try self.renderables.append(renderable);
        } else {
            try self.renderables.insert(@intCast(idx), renderable);
        }
    }

    pub fn remove(self: *Renderer, renderable: *Renderable) void {
        const id = renderable.id;
        for (self.renderables.items, 0..) |item, idx| {
            if (item.id == id) {
                _ = self.renderables.orderedRemove(idx);
                break;
            }
        }
    }

    pub fn render(self: *Renderer) void {
        _ = self;
    }

};


