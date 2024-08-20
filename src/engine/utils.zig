// TODO(probably not): Test this file

const std = @import("std");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const objects = @import("../objects/objects.zig");
const a = @import("../assert/assert.zig");
const math = @import("../math/math.zig");

const never = a.never;
const assert = a.assert;
const GameState = objects.gamestate.GameState;
const Values = objects.Values;

pub const MILLI: isize = 1000;
pub const SECOND: isize = 1000 * 1000;
pub const SECOND_F: f64 = 1000 * 1000;
pub const MINUTE: isize = 60 * SECOND;

pub fn humanTime(timeUS: isize) ![]u8 {
    const minutes = @divFloor(timeUS, MINUTE);
    const seconds = @divFloor(@mod(timeUS, MINUTE), SECOND);

    return std.fmt.bufPrint(scratchBuf(50), "{}m {}s", .{minutes, seconds});
}

pub fn normalize(comptime t: type, current: t, total: t) f64 {
    const cF: f64 = @floatFromInt(current);
    const tF: f64 = @floatFromInt(total);

    return cF / tF;
}

pub fn scale(current: f64, total: f64, s: f64) f64 {
    return (current / total) * s;
}

pub fn getRangeByTeam(gs: *GameState, team: u8) math.Range {
    return switch (team) {
        Values.TEAM_ONE => gs.oneNoBuildTowerRange,
        Values.TEAM_TWO => gs.twoNoBuildTowerRange,
        else => {
            never("invalid team id");
            unreachable;
        }
    };
}

pub fn positionInRange(gs: *GameState, team: u8) math.Position {
    const range = getRangeByTeam(gs, team);
    // TODO: the ranges really are weird here...
    // the problem is that for chat AIs i adjust the range to begin with with
    // the tower row count and then again here leading to the end being lower
    // than the start
    const row = gs.values.randRange(usize, range.startRow, range.endRow);
    const col = gs.values.randRange(usize, objects.tower.TOWER_COL_COUNT, gs.values.cols - objects.tower.TOWER_COL_COUNT);

    return .{
        .row = row,
        .col = col,
    };
}

pub fn aabbInValidRange(gs: *GameState, aabb: math.AABB, team: u8) bool {
    return getRangeByTeam(gs, team).containsAABB(aabb);
}

test "displayTime" {
    const out = try humanTime(69 * 1000 * 1000);
    try std.testing.expectEqualStrings("1m 9s", out);
}
