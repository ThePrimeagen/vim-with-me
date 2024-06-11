const assert = @import("assert");
const std = @import("std");
const Allocator = std.mem.Allocator;

const primitives = @import("primitives.zig");

pub const Engine = struct {
    alloc: Allocator,

    pub fn init(alloc: Allocator) Engine {
        return .{
            .alloc = alloc,
        };
    }

    pub fn add(self: *Engine, item: anytype) !void {
        _ = self;
        _ = item;
    }

    pub fn remove(self: *Engine, item: anytype) void {
        _ = self;
        _ = item;
    }

    pub fn gameLoop(self: *Engine, timePassed: i64) void {
        _ = timePassed;
        _ = self;
    }
};

const milliTimestamp = std.time.milliTimestamp;
pub const RealTime = struct {
    lastTime: i64,

    pub fn init() RealTime {
        return .{
            .lastTime = 0,
        };
    }

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
