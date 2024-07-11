const std = @import("std");
const objects = @import("../objects/objects.zig");
const math = @import("../math/math.zig");
const Params = @import("params.zig");
const never = @import("../assert/assert.zig").never;

const GS = objects.gamestate.GameState;
const Values = objects.Values;

var testValues = blk: {
    var out = Values{
        .rows = 10,
        .cols = 10,
    };
    Values.init(&out);
    break :blk out;
};

pub fn values() *const Values {
    return &testValues;
}

pub fn positionInRange(gs: *GS, team: u8) math.Position {
    const range = switch (team) {
        Values.TEAM_ONE => gs.oneCreepRange,
        Values.TEAM_TWO => gs.twoCreepRange,
        else => {
            never("invalid team id");
            unreachable;
        }
    };

    const row = gs.values.randRange(usize, range.startRow, range.endRow);
    const col = gs.values.randRange(usize, 2, gs.values.cols - objects.tower.TowerSize);

    return .{
        .row = row,
        .col = col,
    };
}
