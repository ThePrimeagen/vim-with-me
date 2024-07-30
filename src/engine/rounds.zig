const a = @import("../assert/assert.zig");
const assert = a.assert;
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
    return 1 + @divFloor(round, 3);
}

pub fn placeRandomCreep(gs: *GS, range: math.Range, team: u8) !void {
    const row = gs.values.randRange(usize, range.startRow, range.endRow);
    _ = try gamestate.placeCreep(gs, .{
        .col = 0,
        .row = row,
    }, team);
}

pub const CreepSpawner = struct {
    gs: *GS,
    spawnStartUS: i64 = 0,
    spawnDurationUS: i64 = 10_000_000,
    spawnCount: usize = 0,
    spawned: usize = 0,

    pub fn init(gs: *GS) CreepSpawner {
        return .{
            .gs = gs,
        };
    }

    pub fn startRound(self: *CreepSpawner) void {
        const rr: f64 = @floatFromInt(self.gs.round * self.gs.round);
        self.spawnCount = @as(usize,
            @intFromFloat(rr * self.gs.values.creep.scaleSpawnRate)) + self.gs.round;
        self.spawned = 0;
        self.spawnStartUS = self.gs.time;
    }

    pub fn tick(self: *CreepSpawner) !void {
        if (self.spawned == self.spawnCount) {
            return;
        }

        const dur = self.gs.time - self.spawnStartUS;
        const norm = engine.utils.normalize(i64, dur, self.spawnDurationUS);

        const count: usize = @intFromFloat(@as(f64, @floatFromInt(self.spawnCount)) * norm + 1);
        const diff = @min(self.spawnCount - self.spawned, count - self.spawned);

        for (0..diff) |_| {
            try self.placeCreep(self.gs.oneCreepRange, Values.TEAM_ONE);
            try self.placeCreep(self.gs.twoCreepRange, Values.TEAM_TWO);
            self.spawned += 1;
        }

        assert(self.spawned <= self.spawnCount, "spawned too many creeps");
    }

    pub fn placeCreep(self: *CreepSpawner, range: math.Range, team: u8) !void {
        const row = self.gs.values.randRange(usize, range.startRow, range.endRow);
        _ = try gamestate.placeCreep(self.gs, .{
            .col = 0,
            .row = row,
        }, team);
    }
};
