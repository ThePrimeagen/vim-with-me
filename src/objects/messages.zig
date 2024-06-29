const assert = @import("assert").assert;
const std = @import("std");
const math = @import("math");

const Coord = math.Coord;
const Allocator = std.mem.Allocator;

var scratchBuf: [8196]u8 = undefined;

pub const NextRound = struct {
    pub fn init(msg: []const u8) ?NextRound {
        if (msg.len == 1 and msg.ptr[0] == 'n') {
            return .{ };
        }
        return null;
    }
};

pub const PossibleCoord = struct {
};

pub const Message = union(enum) {
    round: NextRound,
    coord: Coord,
    possibleCoord: PossibleCoord,

    pub fn init(msg: []const u8) ?Message {
        const coord = Coord.init(msg);
        if (coord) |c| {
            return .{.coord = c};
        }

        const next = NextRound.init(msg);
        if (next) |n| {
            return .{.round = n};
        }

        return null;
    }
};

