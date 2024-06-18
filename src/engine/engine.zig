const std = @import("std");
const Allocator = std.mem.Allocator;

const r = @import("primitives.zig");
const u = @import("update.zig");

pub const Engine = struct {
    renderer: r.Renderer,
    updater: u.Updater,

    pub fn init(alloc: Allocator) Engine {
    }
};

