const assert = @import("assert");
const std = @import("std");
const Allocator = std.mem.Allocator;

const primitives = @import("primitives.zig");
const Renderer = @import("render.zig").Renderer;
const UpdateableList = std.ArrayList(*primitives.Updateable);

pub const Engine = struct {
    renderer: Renderer,
    updater: UpdateableList,
    alloc: Allocator,

    pub fn init(alloc: Allocator) Engine {
        return .{
            .alloc = alloc,
            .renderer = Renderer.init(alloc),
            .updater = UpdateableList.init(alloc),
        };
    }

    pub fn add(self: *Engine, item: anytype) !void {
        const Struct = @TypeOf(item);
        const StructInfo = @typeInfo(Struct);

        comptime assert(StructInfo == .Struct, "expected a struct");

        switch (Struct) {
            inline primitives.Renderable => |renderable| {
                try self.renderer.add(renderable);
            },
            inline primitives.Updateable => |updateable| {
                try self.updater.append(updateable);
            },
            inline else => unreachable,
        }
    }

    pub fn remove(self: *Engine, item: anytype) void {
        const Struct = @TypeOf(item);
        const StructInfo = @typeInfo(Struct);

        comptime assert(StructInfo == .Struct, "expected a struct");

        switch (Struct) {
            inline primitives.Renderable => |renderable| {
                self.renderer.remove(renderable);
            },
            inline primitives.Updateable => |updateable| {
                const id = updateable.id();
                for (self.updater.items, 0..) |i, idx| {
                    if (i.id() == id) {
                        self.updater.orderedRemove(idx);
                        return;
                    }
                }
                assert(false, "you have tried to remove something that doesn't exist");
            },
            inline else => unreachable,
        }
    }

    pub fn gameLoop(self: *Engine, timePassed: i64) void {
        _ = timePassed;
        _ = self;
    }
};

const milliTimestamp = std.time.milliTimestamp();
pub const RealTime = struct {
    lastTime: i64,

    pub fn reset(self: *RealTime) void {
        const time = milliTimestamp();
        self.lastTime = time;
    }

    pub fn tick(self: *RealTime) i64 {
        const time = milliTimestamp();
        const delta = time - self.lastTime;
        self.lastTime = time;

        return delta;
    }

};

pub const FauxTime = struct {
    returnDelta: i64,

    pub fn tick(self: *FauxTime) i64 {
        return self.returnDelta;
    }

};
