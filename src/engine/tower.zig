const math = @import("math");
const objects = @import("objects");
const std = @import("std");
const unwrap = @import("assert").unwrap;

const Tower = objects.tower.Tower;
const colors = objects.colors;
const Color = colors.Color;
const Cell = colors.Cell;
const Black = colors.Black;
const GS = objects.gamestate.GameState;
const Target = objects.gamestate.Target;

const TowerSize = 3;
const TowerCell: [TowerSize]Cell = .{
    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
};

pub fn contains(self: *Tower, pos: math.Vec2) bool {
    if (!self.alive) {
        return false;
    }

    if (pos.row() != self.pos.row()) {
        return false;
    }

    return pos.sub(self.pos.add(.{.x = 1, .y = 0})).lenSq() <= 1;

}

pub fn find(self: *Tower, gs: *GS) void {
    _ = self;
    _ = gs;
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

    pub fn tower(t: TowerBuilder) Tower {
        return .{
            .id = t._id,
            .pos = t._pos,
            .team = t._team,
            .rSized = .{
                .pos = t._pos.position(),
                .cols = TowerSize,
            },
            .rCells = TowerCell,
        };
    }
};

pub fn upgrade(self: *Tower) void {
    if (self.level < 9) {
        self.level += 1;
    }
}

pub fn projectile(self: *Tower, gs: *GS, target: Target) void {
    _ = self;
    _ = gs;
    _ = target;
}

pub fn update(self: *Tower, gs: *GS) void {
    if (!self.alive) {
        return;
    }

}

pub fn render(self: *Tower, gs: *GS) void {

    const life = getLifePercent(self);
    const sqLife = life * life;

    self.rCells[1].text = '0' + self.level;
    color(self, .{
        .r = @intFromFloat(255.0 * life),
        .b = @intFromFloat(255.0 * sqLife),
        .g = @intFromFloat(255.0 * sqLife),
    });

    _ = gs;
}

fn getLifePercent(self: *Tower) f64 {
    const max: f64 = @floatFromInt(self.maxAmmo);
    const ammo: f64 = @floatFromInt(self.ammo);
    return ammo / max;
}

var testId: usize = 0;
fn getTestId() usize {
    const out = testId;
    out += 1;
    return out;
}

fn createTestTower() Tower {
    return .{
        .id = getTestId(),
        .pos = math.ZERO_VEC2,
        .team = 0,
        .rSized = .{
            .pos = math.ZERO_POS,
            .cols = TowerSize,
        },
        .rCells = TowerCell,
    };
}

const testing = std.testing;
test "tower contains" {
    var t = createTestTower();
    try testing.expect(contains(&t, .{.x = -0.9999, .y = 0}));
    try testing.expect(contains(&t, .{.x = 0.9999, .y = 0}));
    try testing.expect(contains(&t, .{.x = 0, .y = 0}));

    // col
    try testing.expect(!contains(&t, .{.x = 1.1, .y = 0}));
    try testing.expect(!contains(&t, .{.x = -1.1, .y = 0}));

    // row
    try testing.expect(!contains(&t, .{.x = 0, .y = 1}));
    try testing.expect(!contains(&t, .{.x = 0, .y = -1}));
}

