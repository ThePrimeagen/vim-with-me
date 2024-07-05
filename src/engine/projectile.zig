const objects = @import("../objects/objects.zig");
const Projectile = objects.projectile.Projectile;
const GS = objects.gamestate.GameState;

pub fn update(self: *Projectile, gs: *GS) void {
    _ = self;
    _ = gs;
}
