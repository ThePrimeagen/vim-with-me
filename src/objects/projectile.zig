const math = @import("math");
const colors = @import("colors.zig");

pub const Projectile = struct {
    id: usize,
    createdBy: usize,
    target: usize,
    targetType: u8, // i have to look up creep or tower
    team: u8,

    pos: math.Vec2,
    life: u16,
    speed: f32,
    alive: bool = true,

    // rendered
    rLife: u16,
    rColor: colors.Color,
    rText: u8,
};


