const std = @import("std");

const assert = @import("assert").assert;
const r = @import("renderable.zig");
const Renderable = r.Renderable;

const RenderList = std.ArrayList(*Renderable);

pub const Renderer = struct {
    renderables: RenderList,

    pub fn init(alloc: std.mem.Allocator) Renderer {
        return .{ .renderables = RenderList.init(alloc) };
    }

    pub fn deinit(self: *Renderer) void {
        self.renderables.deinit();
    }

    pub fn add(self: *Renderer, renderable: *Renderable) void {
        var lo: usize = 0;
        var hi: usize = self.renderables.items.len;

        const needle = renderable.z();
        var idx: isize = -1;

        while (lo < hi) {
            const midpoint = lo + (hi - lo) / 2;
            const value = self.renderables.items[midpoint].z();

            if (value == needle) {
                idx = midpoint;
                break;
            } else if (value > needle) {
                hi = midpoint;
            } else if (value < needle) {
                lo = midpoint + 1;
            }
        }

        if (idx == -1) {
            self.renderables.append(renderable);
        } else {
            self.renderables.insert(idx, renderable);
        }
    }

    pub fn remove(self: *Renderer, renderable: *Renderable) void {
        const id = renderable.id();
        for (self.renderables.items, 0..) |*item, idx| {
            if (item.id() == id) {
                self.renderables.orderedRemove(idx);
                break;
            }
        }
    }

    pub fn render(self: *Renderer) void {
        _ = self;
    }
};

test "adding some renderables" {
}
