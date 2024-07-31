const std = @import("std");
const objects = @import("../objects/objects.zig");
const math = @import("../math/math.zig");
const Params = @import("../params.zig");
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

