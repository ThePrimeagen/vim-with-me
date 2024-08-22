const assert = @import("../assert/assert.zig").assert;
const std = @import("std");
const math = @import("../math/math.zig");
const objects = @import("../objects/objects.zig");
const scratch = @import("../scratch/scratch.zig");

const Allocator = std.mem.Allocator;
const Values = objects.Values;

var scratchBuf: [8196]u8 = undefined;

pub const PlayRound = struct {
    pub fn init(msg: []const u8) ?PlayRound {
        if (msg.len == 1 and msg.ptr[0] == 'p') {
            return .{ };
        }
        return null;
    }
};

pub const Countdown = struct {
    countdownUS: i64,

    pub fn init(msg: []const u8) !?Countdown {
        if (msg.ptr[0] != 'c') {
            return null;
        }

        return .{
            .countdownUS = try std.fmt.parseInt(i64, msg.ptr[1..msg.len], 10),
        };
    }
};

pub const Message = union(enum) {
    round: PlayRound,
    positions: math.PossiblePositions,
    countdown: Countdown,

    pub fn init(msg: []const u8) !?Message {
        try std.io.getStdErr().writeAll(try std.fmt.bufPrint(scratch.scratchBuf(500), "message received: {s}\n", .{msg}));

        const set = math.PossiblePositions.init(msg);
        if (set) |s| {
            return .{.positions = s};
        }

        const next = PlayRound.init(msg);
        if (next) |n| {
            return .{.round = n};
        }

        const countdown = try Countdown.init(msg);
        if (countdown) |c| {
            return .{.countdown = c};
        }

        return null;
    }
};

