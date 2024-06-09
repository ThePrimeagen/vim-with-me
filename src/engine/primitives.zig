const assert = @import("assert");

pub const Color = struct {
    r: u8,
    g: u8,
    b: u8,
};

pub const Tower = struct {
    id: usize,
    team: u8,

    // position
    row: u8,
    col: u8,
    ammo: u16,
    dead: bool,

    // rendered
    rRow: u8,
    rCol: u8,
    rAmmo: u16,
    rColor: [9]Color,
    rText: [9]u8,
};

pub const Projectile = struct {
    id: usize,
    team: u8,

    row: u8,
    col: u8,
    life: u16,
    speed: f32,
    dead: bool,

    // rendered
    rRow: u8,
    rCol: u8,
    rLife: u16,
    rColor: Color,
    rText: u8,
};

pub const Creep = struct {
    id: usize,
    team: u8,

    row: u8,
    col: u8,
    life: u16,
    speed: f32,
    dead: bool,

    // rendered
    rRow: u8,
    rCol: u8,
    rLife: u16,
    rColor: Color,
    rText: u8,
};
