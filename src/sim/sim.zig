const std = @import("std");
const assert = @import("../assert/assert.zig");
const objects = @import("../objects/objects.zig");
const Params = @import("../params.zig");
const SimFile = @import("file.zig");
const RandSim = @import("rand.zig");

const GS = objects.gamestate.GameState;
const never = assert.never;

pub const Sim = union(enum) {
    rand: *const fn(gs: *GS) (std.mem.Allocator.Error || std.fmt.BufPrintError)!void,
    file: SimFile,
};

pub fn fromParams(params: *const Params) !Sim {
    if (params.simulationType) |sim| {
        if (std.mem.startsWith(u8, sim, "rand")) {
            return .{
                .rand = RandSim.randPlacement,
            };
        } else if (std.mem.startsWith(u8, sim, "file:") or std.mem.eql(u8, sim, "stdin")) {
            return .{
                .file = try SimFile.fromParams(params),
            };
        }

        never("invalid simulation type");
    }

    return .{
        .rand = RandSim.randPlacement,
    };
}

pub fn simulate(s: *Sim, gs: *GS) !void {
    switch (s.*) {
        .rand => |r| try r(gs),
        .file => |*f| {
            try f.next(gs);
        },
    }
}

