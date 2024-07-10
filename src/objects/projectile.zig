const math = @import("../math/math.zig");
const colors = @import("colors.zig");
const Target = @import("target.zig").Target;

const INITIAL_PROJECTILE_COLOR: colors.Color = .{.r = 0, .g = 255, .b = 0};

pub const ProjectileSize = 1;
pub const ProjectileCell: [1]colors.Cell = .{
    .{.text = 'x', .color = INITIAL_PROJECTILE_COLOR },
};

pub const Projectile = struct {
    id: usize,
    target: Target,

    createdAt: i64,
    maxTimeAlive: i64,

    pos: math.Vec2,
    speed: f64 = 0,
    alive: bool = true,
    deadUS: i64 = 0,
    damage: usize,

    // rendered
    rColor: colors.Color = colors.Red,
    rSized: math.Sized = .{
        .pos = math.ZERO_POS,
        .cols = ProjectileSize,
    },
    rCells: [1]colors.Cell = ProjectileCell,
};


