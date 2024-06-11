const assert = @import("assert");
const std = @import("std");
const Allocator = std.mem.Allocator;

const primitives = @import("primitives.zig");

pub const Engine = struct {
    alloc: Allocator,
    foo: ?*u8,

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
pub const Time = struct {
    delta: i64,
    createdAt: i64,

    pub fn init(delta: i64, createdAt: i64) Time {
        return .{
            .createdAt = createdAt,
            .delta = delta,
        };
    }

    pub fn since(self: Time) i64 {
        const now = milliTimestamp();
        return now - self.createdAt;
    }

    pub fn sleep(self: Time, total: i64) void {
        const delta = self.since();
        const remaining: u64 = @intCast(total - delta);
        if (remaining > 0) {
            std.time.sleep(remaining * 1_000_000);
        }
    }
};

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

    pub fn tick(self: *RealTime) Time {
        const time = milliTimestamp();
        const delta = time - self.lastTime;
        self.lastTime = time;

        return Time.init(delta, time);
    }

};

pub const FauxTime = struct {
    returnDelta: i64,

    pub fn tick(self: *FauxTime) i64 {
        return self.returnDelta;
    }

};
