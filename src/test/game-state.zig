const std = @import("std");
const objects = @import("../objects/objects.zig");
const engine = @import("../engine/engine.zig");
const Params = @import("params.zig");

const GS = objects.gamestate.GameState;
const gamestate = engine.gamestate;

pub const Spawner = struct {
    gs: *GS,
    lastSpawn: isize = 0,
    currentTime: isize = 100_000_000,
    spawnRate: isize,
    params: *Params,

    pub fn init(params: *Params, gs: *GS) Spawner {
        return .{
            .spawnRate = @intCast(params.creepRate),
            .gs = gs,
            .params = params,
        };
    }

    pub fn tick(self: *Spawner, deltaUS: isize) !void {
        if (self.gs.noBuildZone) {
            return;
        }

        self.currentTime += deltaUS;
        if (self.currentTime - self.lastSpawn > self.spawnRate) {
            self.lastSpawn = self.currentTime;

            const range = self.gs.oneCreepRange;
            const row = self.params.randRange(usize, range.startRows, range.endRow);

            _ = try gamestate.placeCreep(self.gs, .{
                .col = 0,
                .row = row,
            }, objects.Values.TEAM_ONE);
        }
    }
};
