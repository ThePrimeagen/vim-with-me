const math = @import("../math/math.zig");
const colors = @import("colors.zig");

pub const Tower = struct {
    id: usize,
    team: u8,

    pos: math.Vec2 = math.ZERO_VEC2,
    aabb: math.AABB = math.ZERO_AABB,

    maxAmmo: usize = 0,
    ammo: usize = 0,

    alive: bool = true,
    deadTimeUS: u64 = 0,

    level: u8 = 1,
    radius: u8 = 1,
    damage: u8 = 1,

    fireRateUS: i64 = 0,
    lastFiringUS: i64 = 0,
    firingDurationUS: i64 = 0,
    firing: bool = false,
    fired: bool = false,

    rows: usize = 1,
    cols: usize = 3,

    // rendered
    rSized: math.Sized = math.ZERO_SIZED,
    rAmmo: u16 = 0,
    rCells: [3]colors.Cell,
};


