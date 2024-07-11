const std = @import("std");

const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const math = @import("../math/math.zig");
const colors = @import("colors.zig");

pub const TowerSize = 3;
pub const TOWER_AABB: math.AABB = .{
    .min = math.ZERO_VEC2,
    .max = .{
        .y = 1,
        .x = TowerSize,
    },
};

pub const Tower = struct {
    id: usize,
    team: u8,

    pos: math.Vec2 = math.ZERO_VEC2,
    aabb: math.AABB = math.ZERO_AABB,
    firingRangeAABB: math.AABB = math.ZERO_AABB,

    maxAmmo: usize = 0,
    ammo: usize = 0,

    alive: bool = true,
    deadTimeUS: i64 = 0,

    level: u8 = 1,
    damage: usize = 1,

    fireRateUS: i64 = 0,
    lastFiringUS: i64 = 0,
    firingDurationUS: i64 = 0,
    firing: bool = false,
    fired: bool = false,

    // rendered
    rSized: math.Sized = math.ZERO_SIZED,
    rAmmo: u16 = 0,
    rCells: [3]colors.Cell,

    // TODO(render): THis is a bad plan
    rows: usize = 1,
    cols: usize = 3,

    pub fn string(self: *Tower) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(250),
            "Tower({}, {}): level={}, damage={}, ammo={}\naabb={s} firing={s}", .{
                self.id, self.team,
                self.level, self.damage, self.ammo,
                try self.aabb.string(), try self.firingRangeAABB.string()
        });
    }
};


