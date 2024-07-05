const math = @import("../math/math.zig");
const colors = @import("colors.zig");
const Target = @import("target.zig").Target;

pub const PROJECTILE_SPEED = 4;

pub const Projectile = struct {
    id: usize,
    target: Target,

    pos: math.Vec2,
    speed: f64 = PROJECTILE_SPEED,
    alive: bool = true,

    // rendered
    rColor: colors.Color = colors.Red,
    rText: u8 = '*',
};


