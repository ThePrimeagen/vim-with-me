const std = @import("std");
const assert = @import("../assert/assert.zig").assert;

const objects = @import("../objects/objects.zig");
const Projectile = objects.projectile.Projectile;
const GS = objects.gamestate.GameState;

fn kill(self: *Projectile, gs: *GS) void {
    self.alive = false;
    self.deadUS = gs.time;
    gs.fns.?.strike(gs, self);
}

pub fn update(self: *Projectile, gs: *GS) !void {
    // TODO: Change the fns definition of gs and make only creatable via
    // some init function
    assert(gs.fns != null, "fns must be defined for this function");

    if (!self.alive) {
        return;
    }

    const alive = switch (self.target) {
        .creep => |c| gs.creeps.items[c].alive,
        .tower => |t| gs.towers.items[t].alive,
    };

    if (!alive) {
        kill(self, gs);
        return;
    }

    assert(gs.time - self.createdAt < self.maxTimeAlive, "bullet has lived longer than expected");

    const target = switch (self.target) {
        .creep => |c| gs.creeps.items[c].pos,
        .tower => |t| gs.towers.items[t].pos,
    };

    const to = self.pos.sub(target);
    const len = to.len();
    const lenUS = (len / self.speed) * 1_000_000;
    const delta = @as(f64, @floatFromInt(gs.loopDeltaUS));
    const deltaP = delta / 1_000_000;

    if (delta > lenUS) {
        kill(self, gs);
        return;
    }

    self.pos = self.pos.add(to.norm().scale(-deltaP * self.speed));
}

pub fn render(self: *Projectile, gs: *GS) void {
    self.rSized.pos = self.pos.position();
    _ = gs;
}
