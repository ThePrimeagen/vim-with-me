const math = @import("math");
const objects = @import("objects");

const Tower = objects.tower.Tower;
const colors = objects.colors;
const Color = colors.Color;
const Cell = colors.Cell;
const Black = colors.Black;
const GS = objects.gamestate.GameState;

const TowerSize = 3;
const TowerCell: [TowerSize]Cell = .{
    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
};

var id: usize = 0;
fn nextTowerId() usize {
    const out = id;
    id += 1;
    return out;
}

pub fn contains(self: *Tower, pos: math.Position) bool {
    if (self.alive) {
        return false;
    }

    const c = math.absUsize(self.pos.col, pos.col);
    return self.pos.row == pos.row and c <= 1;
}

pub fn color(self: *Tower, c: Color) void {
    for (0..self.rCells.len) |idx| {
        self.rCells[idx].color = c;
    }
}

pub const TowerBuilder = struct {
    _id: usize = 0,
    _team: u8 = 0,
    _pos: math.Vec2 = math.ZERO_POS_F,
    _rSized: math.Sized = math.ZERO_SIZED,

    pub fn start() TowerBuilder {
        return .{
            .id = nextTowerId(),
        };
    }

    pub fn team(t: TowerBuilder, te: u8) TowerBuilder {
        t.team = te;
        return t;
    }

    pub fn pos(t: TowerBuilder, p: math.Position) TowerBuilder {
        t.pos = p;
        return t;
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

pub fn update(self: *Tower, gs: *GS) void {
    if (!self.alive) {
        return;
    }
    _ = gs;
}

pub fn render(self: *Tower, gs: *GS) void {

    const life = self.getLifePercent();
    const sqLife = life * life;

    self.rCells[1].text = '0' + self.level;
    self.color(.{
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
