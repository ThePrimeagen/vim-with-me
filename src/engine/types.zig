const std = @import("std");

const INITIAL_AMMO = 50;

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
const TowerCell: [9]Cell = .{
    .{.text = ' ', .color = Black },
    .{.text = '^', .color = Black },
    .{.text = ' ', .color = Black },

    .{.text = '/', .color = Black },
    .{.text = '*', .color = Black },
    .{.text = '\\', .color = Black },

    .{.text = '<', .color = Black },
    .{.text = '_', .color = Black },
    .{.text = '>', .color = Black },
};
const ZERO_POS: Position = .{.row = 0, .col = 0};
const ZERO_SIZED: Sized = .{.cols = 3, .pos = ZERO_POS};

pub const Tower = struct {
    id: usize,
    team: u8,

    // position
    pos: Position = ZERO_POS,
    ammo: u16 = INITIAL_AMMO,
    dead: bool = false,
    level: u8 = 1,
    radius: u8 = 1,
    damage: u8 = 1,

    // rendered
    rSized: Sized = ZERO_SIZED,
    rAmmo: u16 = INITIAL_AMMO,
    rCells: [9]Cell = TowerCell,

    pub fn contains(self: *Tower, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        const r = @abs(self.row - pos.row);
        const c = @abs(self.col - pos.col);
        return r <= 1 and c <= 1;
    }

    pub fn color(self: *Tower, c: Color) void {
        for (self.rCells) |*cell| {
            cell.color = c;
        }
    }

    pub fn create(id: usize, team: u8, pos: Position) Tower {
        return .{
            .id = id,
            .team = team,
            .pos = pos,
            .rSized = .{
                .cols = TowerSize,
                .pos = pos
            },
        };
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

        const r = @abs(self.row - pos.row);
        const c = @abs(self.col - pos.col);
        return r == 0 and c == 0;
    }
};
