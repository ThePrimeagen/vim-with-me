const std = @import("std");

const a = @import("../assert/assert.zig");
const gamestate = @import("game-state.zig");
const towers = @import("tower.zig");
const creeps = @import("creep.zig");
const objects = @import("../objects/objects.zig");

const testing = std.testing;
test "find nearest creep" {
    var values = objects.Values{.rows = 6, .cols = 6};
    values.tower.fireRateUS = 10_000_000; // ensures we don't fire
    objects.Values.init(&values);
    var gs = try objects.gamestate.GameState.init(testing.allocator, &values);
    defer gs.deinit();

    gamestate.init(&gs);
    gs.noBuildZone = false;
    gs.playing = true;

    var gsDump = gs.dumper();
    a.addDump(&gsDump);
    defer a.removeDump(&gsDump);

    const tId = a.option(usize, try gamestate.placeTower(&gs, .{.row = 2, .col = 2}, objects.Values.TEAM_ONE));
    const tower = gamestate.towerById(&gs, tId);

    const one = try gamestate.placeCreep(&gs, .{.row = 1, .col = 0}, objects.Values.TEAM_ONE);
    const two = try gamestate.placeCreep(&gs, .{.row = 1, .col = 1}, objects.Values.TEAM_ONE);
    const three = try gamestate.placeCreep(&gs, .{.row = 1, .col = 2}, objects.Values.TEAM_ONE);

    try testing.expect(towers.withinRange(tower, gs.creeps.items[one].pos) == false);
    try testing.expect(towers.withinRange(tower, gs.creeps.items[two].pos) == true);
    try testing.expect(towers.withinRange(tower, gs.creeps.items[three].pos) == true);

    var creep = towers.creepWithinRange(tower, &gs);

    try testing.expect(creep != null);
    try testing.expect(creep.?.id == three);

    try gamestate.update(&gs, 2_990_000);

    creep = towers.creepWithinRange(tower, &gs);
    a.assert(creep != null, "expected to find the creep");
    try testing.expect(creep.?.id == three);

    try gamestate.update(&gs, 16_000);

    creep = towers.creepWithinRange(tower, &gs);
    try testing.expect(creep != null);
    std.debug.print("{s}\n", .{try creep.?.string()});
    try testing.expect(creep.?.id == two);
}

test "creep distance" {
    var values = objects.Values{.rows = 6, .cols = 7};
    objects.Values.init(&values);
    var gs = try objects.gamestate.GameState.init(testing.allocator, &values);
    defer gs.deinit();
    gamestate.init(&gs);
    gs.noBuildZone = false;

    var gsDump = gs.dumper();
    a.addDump(&gsDump);
    defer a.removeDump(&gsDump);

    const one = try gamestate.placeCreep(&gs, .{.row = 1, .col = 0}, objects.Values.TEAM_ONE);
    const two = try gamestate.placeCreep(&gs, .{.row = 1, .col = 1}, objects.Values.TEAM_ONE);
    const three = try gamestate.placeCreep(&gs, .{.row = 1, .col = 2}, objects.Values.TEAM_ONE);
    const creepList = gs.creeps.items;

    try testing.expect(creeps.distanceToExit(&creepList[one], &gs) == 7.0);
    try testing.expect(creeps.distanceToExit(&creepList[two], &gs) == 6.0);
    try testing.expect(creeps.distanceToExit(&creepList[three], &gs) == 5.0);
}
