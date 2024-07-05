const std = @import("std");

const objects = @import("../objects/objects.zig");
const Projectile = objects.projectile.Projectile;
const GS = objects.gamestate.GameState;

pub fn update(self: *Projectile, gs: *GS) !void {
    if (!self.alive) {
        return;
    }

    const target = switch (self.target) {
        .creep => |c| gs.creeps.items[c].pos,
        .tower => |t| gs.towers.items[t].pos,
    };

    const to = self.pos.sub(target);
    const len = to.len();
    const lenUS = len / self.speed * 1_000_000;
    const delta = @as(f64, @floatFromInt(gs.loopDeltaUS));
    const deltaP = delta / 1_000_000;

    if (delta > lenUS) {
        self.alive = false;
        self.deadUS = gs.time;
        return;
    }

    self.pos = self.pos.add(to.scale(-deltaP * self.speed));
}

pub fn render(self: *Projectile, gs: *GS) void {
    self.rSized.pos = self.pos.position();
    _ = gs;
}
