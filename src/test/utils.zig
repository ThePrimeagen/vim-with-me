const std = @import("std");
const objects = @import("objects");

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
