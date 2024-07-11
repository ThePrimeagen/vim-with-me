const utils = @import("../test/utils.zig");
const math = @import("../math/math.zig");
const objects = @import("../objects/objects.zig");
const std = @import("std");
const a = @import("../assert/assert.zig");
const unwrap = a.unwrap;
const assert = a.assert;
const creeps = @import("creep.zig");

const Values = objects.tower.Tower;
const Tower = objects.tower.Tower;
const colors = objects.colors;
const Color = colors.Color;
const Cell = colors.Cell;
const Black = colors.Black;
const GS = objects.gamestate.GameState;
const Target = objects.gamestate.Target;
const Creep = objects.creep.Creep;

const TowerCell: [objects.tower.TowerSize]Cell = .{
    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
};

pub fn placementAABB(pos: math.Vec2) math.AABB {
    return objects.tower.TOWER_AABB.move(pos).add(.{
        .x = -1,
        .y = 0,
    });
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
                aabb(.{.x = objects.tower.TowerSize + 2, .y = 3}),

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
pub fn upgrade(self: *Tower) void {
    if (true) {
    if (self.level < 9) {
        self.level += 1;
        if (self.level == 3 or self.level == 9) {
            self.firingRangeAABB.min = self.firingRangeAABB.min.sub(ONE_VEC);
            self.firingRangeAABB.max = self.firingRangeAABB.max.add(ONE_VEC);
        }
    }


    const ammo = self.values.tower.ammo * self.level;
    self.ammo = ammo;
    self.maxAmmo = ammo;

    const level: i64 = @intCast(self.level);
    self.fireRateUS = self.values.tower.fireRateUS -
        self.values.tower.scaleFireRateUS * level;

    assert(self.fireRateUS > self.firingDurationUS, "cannot shoot quicker than animation");
    }
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
    if (creepMaybe == null) {
        if (self.firing) {
            self.firing = false;
        }
        return;
    }
    const creep = creepMaybe.?;

    if (self.firing and self.lastFiringUS + self.firingDurationUS < gs.time) {
        _ = try gs.fns.?.placeProjectile(gs, self, .{.creep = creep.id});

        self.ammo -= 1;
        self.fired = true;
    }

    if (self.lastFiringUS + self.fireRateUS > gs.time) {
        return;
    }

    self.lastFiringUS += self.fireRateUS;
    self.firing = true;
}

pub fn render(self: *Tower, gs: *GS) void {

    if (self.firing) {
        const delta: f64 = @floatFromInt(gs.time - self.lastFiringUS);
        const fDuration: f64 = @floatFromInt(self.firingDurationUS);
        const percent = @min(1, delta / fDuration);
        const sqrt = @sqrt(percent);
        const sq = percent * percent;

        color(self, .{
            .r = @intFromFloat(255.0 * sqrt),
            .g = @intFromFloat(255.0 * sqrt),
            .b = @intFromFloat(255.0 * sq),
        });

        return;
    }

    const life = getLifePercent(self);
    const sqLife = life * life;

    self.rCells[1].text = '0' + self.level;
    color(self, .{
        .r = @intFromFloat(255.0 * life),
        .g = @intFromFloat(255.0 * sqLife),
        .b = @intFromFloat(255.0 * sqLife),
    });
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

    try testing.expect(!contains(&t, .{.x = -0.9999, .y = 0}));

    try testing.expect(contains(&t, .{.x = 0.9999, .y = 0}));
    try testing.expect(contains(&t, .{.x = 0, .y = 0}));

    // col
    try testing.expect(contains(&t, .{.x = 2.99, .y = 0}));
    try testing.expect(!contains(&t, .{.x = 3.1, .y = 0}));

    // row
    try testing.expect(!contains(&t, .{.x = 0, .y = 1}));
    try testing.expect(!contains(&t, .{.x = 0, .y = -1}));
}

