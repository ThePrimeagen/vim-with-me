const assert = @import("assert").assert;
const std = @import("std");

const Allocator = std.mem.Allocator;

const INITIAL_AMMO = 50;
const INITIAL_CREEP_LIFE = 10;
const INITIAL_CREEP_SPEED = 1;
const INITIAL_CREEP_COLOR: Color = .{.r = 0, .g = 0, .b = 0};

fn usizeToIsizeSub(a: usize, b: usize) isize {
    const ai: isize = @intCast(a);
    const bi: isize = @intCast(b);
    return a - b;
}

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

    pub fn toIdx(self: *Position, cols: usize) usize {
        return self.row * cols + self.col;
    }

    pub fn fromIdx(idx: usize, cols: usize) Position {
        return .{
            .row = idx / cols,
            .col = idx % cols,
        };
    }

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
pub const Red: Color = .{.r = 255, .g = 0, .b = 0 };
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
const ZERO_POS_F: Vec2 = .{.row = 0.0, .col = 0.0};
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

pub const Vec2 = struct {
    x: f64,
    y: f64,

    pub fn subP(self: Vec2, b: Position) Vec2 {
        const rf = @floatFromInt(b.row);
        const cf = @floatFromInt(b.col);

        return .{
            .x = self.x - cf,
            .y = self.y - rf,
        };
    }

    pub fn add(self: Vec2, b: Vec2) Vec2 {
        return .{
            .x = self.x + b.x,
            .y = self.y + b.y,
        };
    }

    pub fn scale(self: Vec2, s: f64) Vec2 {
        return .{
            .x = self.x * s,
            .y = self.y * s,
        };
    }

    pub fn position(self: *Vec2) Position {
        assert(self.x >= 0, "x cannot be negative");
        assert(self.y >= 0, "y cannot be negative");

        return .{
            .row = @intFromFloat(self.y),
            .col = @intFromFloat(self.x),
        };
    }

    pub fn fromPosition(pos: Position) Vec2 {
        return .{
            .row = @floatFromInt(pos.row),
            .col = @floatFromInt(pos.col),
        };
    }
};

fn path(seen: []const isize, pos: usize, out: []usize) usize {
    var idx: usize = 0;

    var p = pos;
    while (seen[p] != p) {
        out[idx] = p;
        p = @intCast(seen[p]);

        idx += 1;
    }

    std.mem.reverse(usize, out[0..idx]);
    return idx;
}

fn walk(from: usize, pos: usize, cols: usize, board: []const bool, seen: []isize) usize {
    assert(board.len % cols == 0, "board is not a rectangle");
    assert(board.len == seen.len, "board and seen should have the same size");

    if (pos >= seen.len or seen[pos] != -1 or board[pos] == false) {
        return 0;
    }

    seen[pos] = @intCast(from);

    // I am at the end
    if (pos % cols == cols - 1) {
        return pos;
    }

    const iCols: isize = @intCast(cols);
    const directions: [4]isize = .{1, -iCols, iCols, -1};
    for (directions) |d| {
        const iPos: isize = @intCast(pos);
        if (iPos + d < 0) {
            continue;
        }

        const out = walk(pos, @intCast(iPos + d), cols, board, seen);
        if (out != 0) {
            return out;
        }
    }

    return 0;
}


const CreepSize = 1;
const CreepCell: [1]Cell = .{
    .{.text = '*', .color = Red },
};

pub const Creep = struct {
    id: usize,
    team: u8,
    cols: usize,

    pos: Vec2 = ZERO_POS_F,
    life: u16 = INITIAL_CREEP_LIFE,
    speed: f32 = INITIAL_CREEP_SPEED,
    alive: bool = true,

    // rendered
    rPos: Position = ZERO_POS,
    rLife: u16 = INITIAL_CREEP_LIFE,
    rColor: Color = INITIAL_CREEP_COLOR,
    rCells: [1]Cell = CreepCell,
    rSized: Sized = ZERO_SIZED,

    path: []usize,
    pathIdx: usize = 0,
    pathLen: usize = 0,

    pub fn init(alloc: Allocator, rows: usize, cols: usize) !Creep {
        return .{
            .path = try alloc.alloc(u8, rows * cols),
            .cols = cols,
        };
    }

    pub fn initialPosition(creep: *Creep, pos: Position) *Creep {
        creep.pos = Vec2.fromPosition(pos);
        creep.rPos = creep.pos.position();
        creep.rSized = .{
            .cols = CreepSize,
            .pos = creep.rPos,
        };

        return creep;
    }

    fn contains(self: *Creep, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        const myPosition = self.pos.position();
        return myPosition.row == pos.row and myPosition.col == pos.col;
    }

    fn calculatePath(self: *Creep, board: []bool, cols: usize) *Creep {
        assert(board.len % cols == 0, "the length is not a rectangle");

        var seen: []isize = .{-1} ** board.len;
        const pos = self.pos.position();


        const last = walk(pos, pos, cols, board, &seen);
        if (last == 0) {
            // TODO: i need to just go straight forward if i can...
        } else {
            self.pathLen = path(&seen, last, self.path);
            self.pathIdx = 0;
        }

        return self;
    }

    fn update(self: *Creep, delta: u64) void {
        const deltaF: f64 = @floatFromInt(delta);
        const deltaP: f64 = deltaF / 1000.0;

        const to = Position.fromIdx(self.path[self.pathIdx], self.cols);

        self.pos = self.pos.subP(to).scale(-deltaP).add(self.pos);
    }

    fn render(self: *Creep) void {
        self.rPos = self.pos.position();
    }
};

const t = std.testing;
test "bfs" {
    const board = [9]bool{
        true, false, true,
        true, false, true,
        true, true, true,
    };
    var seen: [9]isize = .{-1} ** 9;

    const out = walk(0, 0, 3, &board, &seen);
    try t.expect(out == 8);

    var expected: [9]isize = .{
        0, -1, -1,
        0, -1, -1,
        3, 6, 7,
    };

    try t.expectEqualSlices(isize, &expected, &seen);

    var p: [9]usize = .{0} ** 9;
    const len = path(&seen, out, &p);
    var pExpect = [_]usize{ 3, 6, 7, 8 };

    try t.expect(len == 4);
    try t.expectEqualSlices(usize, &pExpect, p[0..len]);
}
