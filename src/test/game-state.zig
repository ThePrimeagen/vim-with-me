const never = @import("../assert/assert.zig").never;
const std = @import("std");
const objects = @import("../objects/objects.zig");
const engine = @import("../engine/engine.zig");
const Params = @import("params.zig");
const math = @import("../math/math.zig");

const GS = objects.gamestate.GameState;
const gamestate = engine.gamestate;
const Values = objects.Values;

pub const Spawner = struct {
    gs: *GS,
    lastSpawn: usize = 0,
    spawner: *const fn (usize, usize) usize,

    pub fn init(gs: *GS, spawner: *const fn (usize, usize) usize) Spawner {
        return .{
            .gs = gs,
            .spawner = spawner,
        };
    }

    pub fn tick(self: *Spawner) !void {
        if (self.gs.noBuildZone) {
            return;
        }

        const count = self.spawner(self.gs.round, self.lastSpawn);

        for (0..count) |_| {
            try self.placeCreep(self.gs.oneCreepRange, Values.TEAM_ONE);
            try self.placeCreep(self.gs.twoCreepRange, Values.TEAM_TWO);
        }

        self.lastSpawn += count;
    }

    pub fn placeCreep(self: *Spawner, range: math.Range, team: u8) !void {
        const row = self.gs.values.randRange(usize, range.startRow, range.endRow);
        _ = try gamestate.placeCreep(self.gs, .{
            .col = 0,
            .row = row,
        }, team);
    }
};
