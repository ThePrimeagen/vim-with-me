const math = @import("../math/math.zig");
const colors = @import("colors.zig");

const INITIAL_AMMO = 50;
const INITIAL_FIRERATE_MS = 1000;

// TODO: Bring this into some sort of global values
const FIRING_DURATION = 200;

pub const Tower = struct {
    id: usize,
    team: u8,

    pos: math.Vec2 = math.ZERO_VEC2,
    aabb: math.AABB = math.ZERO_AABB,

    maxAmmo: u16 = INITIAL_AMMO,
    ammo: u16 = INITIAL_AMMO,

    alive: bool = true,
    deadTimeUS: u64 = 0,

    level: u8 = 1,
    radius: u8 = 1,
    damage: u8 = 1,

    fireRateMS: u64 = INITIAL_FIRERATE_MS,
    lastFiringUS: u64 = 0,
    firingDurationMS: u64 = FIRING_DURATION,
    firing: bool = false,
    fired: bool = false,

    rows: usize = 1,
    cols: usize = 3,

    // rendered
    rSized: math.Sized = math.ZERO_SIZED,
    rAmmo: u16 = INITIAL_AMMO,
    rCells: [3]colors.Cell,
};


