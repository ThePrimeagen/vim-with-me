const std = @import("std");
const utils = @import("../test/utils.zig");
const engine = @import("../engine/engine.zig");
const objects = @import("../objects/objects.zig");
const Sim = @import("sim.zig").Sim;

const GS = objects.gamestate.GameState;
const TEAM_ONE = objects.Values.TEAM_ONE;
const TEAM_TWO = objects.Values.TEAM_TWO;

//pub fn randPlacement(gs: *GS) (std.mem.Allocator.Error || std.fmt.BufPrintError)!void {
//
//    const cnt: usize = @intCast(@divFloor(engine.gamestate.getTotalTowerPlacement(gs), 2));
//    for (0..cnt) |_| {
//        const one = engine.utils.positionInRange(gs, TEAM_ONE);
//        const two = engine.utils.positionInRange(gs, TEAM_TWO);
//        try engine.gamestate.message(gs, .{ .coord = .{ .pos = one, .team = TEAM_ONE, }, });
//        try engine.gamestate.message(gs, .{ .coord = .{ .pos = two, .team = TEAM_TWO, }, });
//    }
//
//}

