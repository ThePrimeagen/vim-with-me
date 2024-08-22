const std = @import("std");
const objects = @import("objects/objects.zig");
const math = @import("math/math.zig");
const engine = @import("engine/engine.zig");

const Values = objects.Values;

var testValues = blk: {
    var values = objects.Values{
        .rows = 6,
        .cols = 7,
    };
    values.init();
    break :blk values;
};

const testBoard = [42]bool{
    true, true, true, false, true, true, true,
    true, false, true, false, true, false, true,
    true, false, true, false, true, false, true,
    true, false, true, true, true, false, true,
    true, false, false, false, false, false, true,
    true, true, true, true, true, true, true,
};

const testBoard2 = [42]bool{
    true, true, true, true, true, true, true,
    true, true, true, true, true, true, true,
    true, true, true, true, true, true, true,
    true, true, true, true, true, true, true,
    true, true, true, true, true, false, true,
    true, true, true, true, true, true, true,
};

const testBoard3 = [42]bool{
    true, true, true, true, true, true, true,
    true, true, true, true, true, false, true,
    true, true, true, true, true, false, true,
    true, true, true, true, true, false, true,
    true, true, true, true, true, false, true,
    true, true, true, true, true, true, true,
};

pub fn viz() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var gs = try objects.gamestate.GameState.init(allocator, &testValues);
    defer gs.deinit();

    var creep = try engine.creep.create(allocator, 0, Values.TEAM_ONE, &testValues, .{.x = 0, .y = 4});
    defer creep.deinit();
    _ = try engine.creep.calculatePath(&creep, &testBoard3);
}

