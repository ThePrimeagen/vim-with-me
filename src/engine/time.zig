const assert = @import("../assert/assert.zig").assert;
const std = @import("std");
const Allocator = std.mem.Allocator;

const microTimestamp = std.time.microTimestamp;
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
        const now = microTimestamp();
        return now - self.createdAt;
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
        const time = microTimestamp();
        self.lastTime = time;
    }

    pub fn tick(self: *RealTime) Time {
        const time = microTimestamp();
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

pub const FPS = struct {
    fps: i64,
    time: ?Time,
    timer: RealTime, // change this later to an enum

    pub fn init(fps: i64) FPS {
        return .{
            .fps = fps,
            .time = null,
            .timer = RealTime.init(),
        };
    }

    pub fn delta(self: *FPS) i64 {
        const d = self.timer.tick();
        self.time = d;

        return d.delta;
    }

    pub fn sleep(self: *FPS) void {
        assert(self.time != null, "called sleep before delta");
        const since = self.time.?.since();
        const remaining: u64 = @intCast(@max(0, self.fps - since));

        if (remaining > 0) {
            std.time.sleep(remaining * 1_000);
        }
    }
};
