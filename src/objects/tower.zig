const math = @import("math");
const colors = @import("colors.zig");

const INITIAL_AMMO = 50;

pub const Tower = struct {
    id: usize,
    team: u8,

    // position
    pos: math.Vec2 = math.ZERO_POS,
    maxAmmo: u16 = INITIAL_AMMO,
    ammo: u16 = INITIAL_AMMO,
    alive: bool = true,
    level: u8 = 1,
    radius: u8 = 1,
    damage: u8 = 1,

    // rendered
    rSized: math.Sized = math.ZERO_SIZED,
    rAmmo: u16 = INITIAL_AMMO,
    rCells: [3]colors.Cell,
};


