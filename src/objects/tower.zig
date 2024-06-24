const math = @import("math");

const colors = @import("colors.zig");

const Cell = colors.Cell;
const Black = colors.Black;

const INITIAL_AMMO = 50;
const TowerSize = 3;
const TowerCell: [3]Cell = .{
    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
};

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
    rCells: [3]Cell = TowerCell,

    pub fn contains(self: *Tower, pos: math.Position) bool {
        if (self.alive) {
            return false;
        }

        const c = math.absUsize(self.pos.col, pos.col);
        return self.pos.row == pos.row and c <= 1;
    }

    pub fn color(self: *Tower, c: colors.Color) void {
        for (0..self.rCells.len) |idx| {
            self.rCells[idx].color = c;
        }
    }

    pub fn create(id: usize, team: u8, pos: math.Position) Tower {
        var p = pos;
        if (p.col == 0) {
            p.col = 1;
        }

        return .{
            .id = id,
            .team = team,
            .pos = p,
            .rSized = .{
                .cols = TowerSize,
                .pos = p
            },
        };
    }

    pub fn update(self: *Tower) void {
        if (!self.alive) {
            return;
        }
    }

    pub fn render(self: *Tower) void {

        const life = self.getLifePercent();
        const sqLife = life * life;

        self.rCells[1].text = '0' + self.level;
        self.color(.{
            .r = @intFromFloat(255.0 * life),
            .b = @intFromFloat(255.0 * sqLife),
            .g = @intFromFloat(255.0 * sqLife),
        });
    }

    fn getLifePercent(self: *Tower) f64 {
        const max: f64 = @floatFromInt(self.maxAmmo);
        const ammo: f64 = @floatFromInt(self.ammo);
        return ammo / max;
    }
};


