const never = @import("../assert/assert.zig").never;
const std = @import("std");
const objects = @import("../objects/objects.zig");
const engine = @import("../engine/engine.zig");
const math = @import("../math/math.zig");

const GS = objects.gamestate.GameState;
const gamestate = engine.gamestate;
const Values = objects.Values;

pub fn towerCount(gs: *GS) usize {
    return 1 + @divFloor(gs.round, 4);
}

pub fn creepCount(values: *const Values, round: usize) usize {
    _ = values;
    return 1 + @divFloor(round, 2);
}

pub fn placeRandomCreep(gs: *GS, range: math.Range, team: u8) !void {
    const row = gs.values.randRange(usize, range.startRow, range.endRow);
    _ = try gamestate.placeCreep(gs, .{
        .col = 0,
        .row = row,
    }, team);
}

//pub const CreepSpawner = struct {
//    gs: *GS,
//    lastSpawn: usize = 0,
//    spawner: *const fn (usize, usize) usize,
//
//    pub fn init(gs: *GS, spawner: *const fn (usize, usize) usize) CreepSpawner {
//        return .{
//            .gs = gs,
//            .spawner = spawner,
//        };
//    }
//
//    pub fn tick(self: *CreepSpawner) !void {
//        if (self.gs.noBuildZone) {
//            return;
//        }
//
//        const count = self.spawner(self.gs.round, self.lastSpawn);
//
//        for (0..count) |_| {
//            try self.placeCreep(self.gs.oneCreepRange, Values.TEAM_ONE);
//            try self.placeCreep(self.gs.twoCreepRange, Values.TEAM_TWO);
//        }
//
//        self.lastSpawn += count;
//    }
//
//    pub fn placeCreep(self: *CreepSpawner, range: math.Range, team: u8) !void {
//        const row = self.gs.values.randRange(usize, range.startRow, range.endRow);
//        _ = try gamestate.placeCreep(self.gs, .{
//            .col = 0,
//            .row = row,
//        }, team);
//    }
//};
