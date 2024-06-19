const std = @import("std");

const INITIAL_AMMO = 50;

fn absUsize(a: usize, b: usize) usize {
    if (a > b) {
        return a - b;
    }
    return b - a;
}

const needle: [1]u8 = .{','};
pub const Position = struct {
    row: usize,
    col: usize,

    pub fn init(str: []const u8) ?Position {
        const idxMaybe = std.mem.indexOf(u8, str, ",");
        if (idxMaybe == null) {
            return null;
        }

        const idx = idxMaybe.?;
        const row = std.fmt.parseInt(usize, str[0..idx], 10) catch {
            return null;
        };
        const col = std.fmt.parseInt(usize, str[idx + 1..], 10) catch {
            return null;
        };

        return .{
            .row = row,
            .col = col,
        };
    }
};

pub const Sized = struct {
    cols: usize,
    pos: Position,
};

pub const Coord = struct {
    pos: Position,
    team: u8,

    pub fn init(msg: []const u8) ?Coord {
        const teamNumber = msg[0];
        if (teamNumber != '1' and teamNumber != '2') {
            return null;
        }
        const pos = Position.init(msg[1..]);
        if (pos == null) {
            return null;
        }

        return .{
            .pos = pos.?,
            .team = teamNumber,
        };
    }
};

pub const NextRound = struct {
    pub fn init(msg: []const u8) ?NextRound {
        if (msg.len == 1 and msg.ptr[0] == 'n') {
            return .{ };
        }
        return null;
    }
};

pub const Message = union(enum) {
    round: NextRound,
    coord: Coord,

    pub fn init(msg: []const u8) ?Message {
        const coord = Coord.init(msg);
        if (coord) |c| {
            return .{.coord = c};
        }

        const next = NextRound.init(msg);
        if (next) |n| {
            return .{.round = n};
        }

        return null;
    }
};

pub const Color = struct {
    r: u8,
    g: u8,
    b: u8,

    pub fn equal(self: Color, other: Color) bool {
        return self.r == other.r and
            self.g == other.g and
            self.b == other.b;
    }

};

pub const Black: Color = .{.r = 0, .g = 0, .b = 0 };
pub const Cell = struct {
    text: u8,
    color: Color,
};

const TowerSize = 3;
const TowerCell: [3]Cell = .{
    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },
};
const ZERO_POS: Position = .{.row = 0, .col = 0};
const ZERO_SIZED: Sized = .{.cols = 3, .pos = ZERO_POS};

pub const Tower = struct {
    id: usize,
    team: u8,

    // position
    pos: Position = ZERO_POS,
    maxAmmo: u16 = INITIAL_AMMO,
    ammo: u16 = INITIAL_AMMO,
    dead: bool = false,
    level: u8 = 1,
    radius: u8 = 1,
    damage: u8 = 1,

    // rendered
    rSized: Sized = ZERO_SIZED,
    rAmmo: u16 = INITIAL_AMMO,
    rCells: [3]Cell = TowerCell,

    pub fn contains(self: *Tower, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        const c = absUsize(self.pos.col, pos.col);
        return self.pos.row == pos.row and c <= 1;
    }

    pub fn color(self: *Tower, c: Color) void {
        for (0..self.rCells.len) |idx| {
            self.rCells[idx].color = c;
        }
    }

    pub fn create(id: usize, team: u8, pos: Position) Tower {
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

pub const Projectile = struct {
    id: usize,
    createdBy: usize,
    target: usize,
    targetType: u8, // i have to look up creep or tower
    team: u8,

    pos: Position,
    life: u16,
    speed: f32,
    dead: bool,

    // rendered
    rPos: Position,
    rLife: u16,
    rColor: Color,
    rText: u8,
};

pub const Creep = struct {
    id: usize,
    team: u8,

    pos: Position,
    life: u16,
    speed: f32,
    dead: bool,

    // rendered
    rPos: Position,
    rLife: u16,
    rColor: Color,
    rText: u8,

    fn contains(self: *Creep, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        return self.pos.row == pos.row and self.pos.col == pos.col;
    }
};
