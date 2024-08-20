const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const utils = @import("../test/utils.zig");
const math = @import("../math/math.zig");
const objects = @import("../objects/objects.zig");
const std = @import("std");
const a = @import("../assert/assert.zig");
const creeps = @import("creep.zig");

const never = a.never;
const unwrap = a.unwrap;
const assert = a.assert;
const Values = objects.Values;
const Tower = objects.tower.Tower;
const colors = objects.colors;
const Color = colors.Color;
const Cell = colors.Cell;
const Black = colors.Black;
const GS = objects.gamestate.GameState;
const Target = objects.gamestate.Target;
const Creep = objects.creep.Creep;

pub const OnePlacementTowerCell: [objects.tower.TOWER_CELL_COUNT]Cell = .{
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},

    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},

    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamOneTowerColor},
};

pub const TwoPlacementTowerCell: [objects.tower.TOWER_CELL_COUNT]Cell = .{
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},

    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},

    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
    .{.text = ' ', .color = Black, .background = colors.TeamTwoTowerColor},
};

const TowerCell: [objects.tower.TOWER_CELL_COUNT]Cell = .{
    .{.text = ' ', .color = Black },
    .{.text = ' ', .color = Black },
    .{.text = '^', .color = Black },
    .{.text = ' ', .color = Black },
    .{.text = ' ', .color = Black },

    .{.text = ' ', .color = Black },
    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
    .{.text = ' ', .color = Black },

    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
};

pub fn placementAABB(pos: math.Vec2) math.AABB {
    return objects.tower.TOWER_AABB.move(pos);
}

pub fn contains(self: *Tower, pos: math.Vec2) bool {
    if (!self.alive) {
        return false;
    }

    return self.aabb.contains(pos);
}

pub fn withinRange(self: *Tower, pos: math.Vec2) bool {
    if (!self.alive) {
        return false;
    }

    return self.firingRangeAABB.contains(pos);
}

pub fn color(self: *Tower, c: Color) void {
    for (0..self.rCells.len) |idx| {
        self.rCells[idx].color = c;
    }
}

pub const TowerBuilder = struct {
    _id: usize = 0,
    _team: u8 = 0,
    _pos: math.Vec2 = math.ZERO_VEC2,
    _rSized: math.Sized = math.ZERO_SIZED,

    pub fn start() TowerBuilder {
        return .{ };
    }

    pub fn id(t: TowerBuilder, _id: usize) TowerBuilder {
        var tow = t;
        tow._id = _id;
        return tow;
    }

    pub fn team(t: TowerBuilder, myTeam: u8) TowerBuilder {
        var tow = t;
        tow._team = myTeam;
        return tow;
    }

    pub fn pos(t: TowerBuilder, p: math.Position) TowerBuilder {
        var tow = t;
        tow._pos = p.vec2();
        return tow;
    }

    pub fn tower(t: TowerBuilder, values: *const objects.Values) Tower {
        return .{
            .id = t._id,
            .pos = t._pos,
            .team = t._team,
            .rSized = .{
                .pos = t._pos.position(),
                .cols = objects.tower.TowerSize,
            },
            .aabb = objects.tower.TOWER_AABB.move(t._pos),

            // TODO(render): when i make the tower pretty... this will change
            .firingRangeAABB = t._pos.
                sub(.{.x = 1, .y = 1}).
                aabb(.{.x = objects.tower.TOWER_COL_COUNT + 2, .y = objects.tower.TOWER_ROW_COUNT + 2}),

            .rCells = TowerCell,

            .firingDurationUS = values.tower.firingDurationUS,
            .fireRateUS = values.tower.fireRateUS,
            .ammo = values.tower.ammo,
            .maxAmmo = values.tower.ammo,

            .values = values,
        };
    }
};

const ONE_VEC: math.Vec2 = .{.x = 1, .y = 1};
const TWO_VEC: math.Vec2 = .{.x = 2, .y = 2};
pub fn upgrade(self: *Tower) void {
    if (self.level < 9) {
        self.level += 1;
        if (self.level == 3) {
            self.firingRangeAABB.min = self.firingRangeAABB.min.sub(ONE_VEC);
            self.firingRangeAABB.max = self.firingRangeAABB.max.add(ONE_VEC);
        }
        if (self.level == 5) {
            self.firingRangeAABB.min = self.firingRangeAABB.min.sub(ONE_VEC);
            self.firingRangeAABB.max = self.firingRangeAABB.max.add(ONE_VEC);
        }
        if (self.level == 7) {
            self.firingRangeAABB.min = self.firingRangeAABB.min.sub(ONE_VEC);
            self.firingRangeAABB.max = self.firingRangeAABB.max.add(ONE_VEC);
        }
        if (self.level == 9) {
            self.firingRangeAABB.min = self.firingRangeAABB.min.sub(TWO_VEC);
            self.firingRangeAABB.max = self.firingRangeAABB.max.add(TWO_VEC);
        }
    }

    const ammo = self.values.tower.ammo + (self.level - 1) * self.values.tower.ammoPerLevel;
    self.ammo = ammo;
    self.maxAmmo = ammo;

    const level: i64 = @intCast(self.level);
    self.fireRateUS = self.values.tower.fireRateUS -
        self.values.tower.scaleFireRateUS * level;

    self.damage = self.values.tower.damage;
    if (self.level >= 5) {
        self.damage *= 2;
    }
    if (self.level == 9) {
        self.damage *= 3;
    }

    assert(self.fireRateUS > self.firingDurationUS, "cannot shoot quicker than animation");
}

pub fn update(self: *Tower, gs: *GS) !void {

    assert(gs.fns != null, "fns must be defined to use this function");

    if (!self.alive) {
        return;
    }

    if (self.fired) {
        self.fired = false;
        self.firing = false;
    }

    if (self.ammo == 0) {
        kill(self, gs);
        return;
    }

    const creepMaybe = creepWithinRange(self, gs);
    const towerMaybe = towerWithinRange(self, gs);
    if (creepMaybe == null and towerMaybe == null) {
        if (self.firing) {
            self.firing = false;
        }
        return;
    }
    const target: Target = if (creepMaybe) |c| .{.creep = c.id} else .{.tower = towerMaybe.?.id};

    if (self.firing and self.lastFiringUS + self.firingDurationUS < gs.time) {
        _ = try gs.fns.?.placeProjectile(gs, self, target);

        self.ammo -= 1;
        self.fired = true;
    }

    if (self.lastFiringUS + self.fireRateUS > gs.time) {
        return;
    }

    self.lastFiringUS += self.fireRateUS;
    self.firing = true;
}

pub fn render(self: *Tower, gs: *GS) !void {
    _ = gs;

    self.rCells[7].text = '0' + self.level;
    const buf = try std.fmt.bufPrint(scratchBuf(3), "{:03}", .{@min(999, self.ammo)});
    const offset = 11;
    for (offset..offset + 3) |idx| {
        self.rCells[idx].text = buf[idx - offset];
    }

    const c = switch (self.team) {
        '1' => colors.Blue,
        '2' => colors.Orange,
        else => {
            never("unknown tower team");
            unreachable;
        }
    };

    color(self, c);
}

pub fn hurt(self: *Tower, damage: usize) usize {
    const prevAmmo = self.ammo;
    self.ammo -|= damage;

    return @min(prevAmmo, damage);
}

fn getLifePercent(self: *Tower) f64 {
    const max: f64 = @floatFromInt(self.maxAmmo);
    const ammo: f64 = @floatFromInt(self.ammo);
    return ammo / max;
}

var testId: usize = 0;
fn getTestId() usize {
    const out = testId;
    testId += 1;
    return out;
}

pub fn towerWithinRange(self: *Tower, gs: *GS) ?*Tower {
    var out: ?*Tower = null;

    var minAmmo: usize = 4000;
    for (gs.towers.items) |*t| {
        if (!t.alive or t.team == self.team) {
            continue;
        }
        if (t.ammo < minAmmo and self.firingRangeAABB.overlaps(t.aabb)) {
            minAmmo = t.ammo;
            out = t;
        }
    }

    return out;
}

pub fn creepWithinRange(self: *Tower, gs: *GS) ?*Creep {
    var maxDist: f64 = std.math.floatMax(f64);
    var out: ?*Creep = null;

    for (gs.creeps.items) |*c| {
        if (!c.alive or c.team != self.team) {
            continue;
        }

        if (!creeps.completed(c) and c.alive and withinRange(self, c.pos)) {
            const dist = creeps.distanceToExit(c, gs);
            if (dist < maxDist) {
                maxDist = dist;
                out = c;
            }
        }
    }

    return out;
}

pub fn killById(tid: usize, gs: *GS) void {
    kill(&gs.towers.items[tid], gs);
}

pub fn kill(t: *Tower, gs: *GS) void {
    assert(gs.fns != null, "gs must have functions on it");

    t.alive = false;
    t.deadTimeUS = gs.time;

    gs.fns.?.towerDied(gs, t);
}

fn createTestTower() Tower {
    return TowerBuilder.
        start().
        pos(math.ZERO_POS).
        id(getTestId()).
        team(objects.Values.TEAM_ONE).
        tower(utils.values());
}

const testing = std.testing;
test "tower contains" {
    var t = createTestTower();
    const X = 0;
    const Y = 0;
    const ROW = objects.tower.TOWER_ROW_COUNT;
    const COL = objects.tower.TOWER_COL_COUNT;
    const SMALL = 0.0001;

    try testing.expect(!contains(&t, .{.x = -SMALL, .y = Y}));
    try testing.expect(contains(&t, .{.x = SMALL, .y = Y}));
    try testing.expect(contains(&t, .{.x = X, .y = Y}));

    // col
    try testing.expect(contains(&t, .{.x = COL - SMALL, .y = Y}));
    try testing.expect(!contains(&t, .{.x = COL + SMALL, .y = Y}));

    // row
    try testing.expect(contains(&t, .{.x = X, .y = ROW - SMALL}));
    try testing.expect(!contains(&t, .{.x = X, .y = ROW + SMALL}));
}

