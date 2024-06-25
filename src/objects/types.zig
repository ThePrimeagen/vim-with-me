const assert = @import("assert").assert;
const std = @import("std");

const Allocator = std.mem.Allocator;
const math = std.math;
var scratchBuf: [8196]u8 = undefined;

const INITIAL_AMMO = 50;
const INITIAL_CREEP_LIFE = 10;
const INITIAL_CREEP_SPEED = 1;
const INITIAL_CREEP_COLOR: Color = .{.r = 0, .g = 0, .b = 0};

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

pub const Tower = struct {
    id: usize,
    team: u8,

    // position
    pos: Position = ZERO_POS,
    maxAmmo: u16 = INITIAL_AMMO,
    ammo: u16 = INITIAL_AMMO,
    alive: bool = true,
    level: u8 = 1,
    radius: u8 = 1,
    damage: u8 = 1,

    // rendered
    rSized: Sized = ZERO_SIZED,
    rAmmo: u16 = INITIAL_AMMO,
    rCells: [3]Cell = TowerCell,

    pub fn contains(self: *Tower, pos: Position) bool {
        if (self.alive) {
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

pub const Projectile = struct {
    id: usize,
    createdBy: usize,
    target: usize,
    targetType: u8, // i have to look up creep or tower
    team: u8,

    pos: Position,
    life: u16,
    speed: f32,
    alive: bool = true,

    // rendered
    rPos: Position,
    rLife: u16,
    rColor: Color,
    rText: u8,
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

    scratch: []isize,
    path: []usize,
    pathIdx: usize = 0,
    pathLen: usize = 0,
    alloc: Allocator,

    pub fn string(self: *Creep, buf: []u8) !usize {
        var out = try std.fmt.bufPrint(buf, "creep({}, {})\r\n", .{self.id, self.team});
        var len = out.len;

        out = try std.fmt.bufPrint(buf[len..], "  pos = ", .{});
        len += out.len;
        len += try self.pos.string(buf[len..]);
        out = try std.fmt.bufPrint(buf[len..], "  pathIdx = {}\r\nlife = {}\r\n  speed = {}\r\n  alive = {}\n\n)", .{self.pathIdx, self.life, self.speed, self.alive});
        return len + out.len;
    }

    pub fn init(alloc: Allocator, rows: usize, cols: usize) !Creep {
        return .{
            .path = try alloc.alloc(usize, rows * cols),
            .scratch = try alloc.alloc(isize, rows * cols),
            .cols = cols,
            .alloc = alloc,
            .id = 0,
            .team = 0,
        };
    }

    pub fn deinit(self: *Creep) void {
        self.alloc.free(self.path);
        self.alloc.free(self.scratch);
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

    pub fn contains(self: *Creep, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        const myPosition = self.pos.position();
        return myPosition.row == pos.row and myPosition.col == pos.col;
    }

    pub fn calculatePath(self: *Creep, board: []const bool, cols: usize) void {
        assert(board.len % cols == 0, "the length is not a rectangle");
        assert(board.len == self.path.len, "the length of the board is different than what the creep was initialized with.");

        for (0..self.scratch.len) |idx| {
            self.scratch[idx] = -1;
        }

        const pos = self.pos.position().toIdx(self.cols);

        const last = walk(pos, pos, cols, board, self.scratch);
        if (last == 0) {
            // TODO: i need to just go straight forward if i can...
        } else {
            self.pathLen = path(self.scratch, last, self.path);
            self.pathIdx = 0;
        }
    }

    pub fn completed(self: *Creep) bool {
        return self.pos.position().col == self.cols - 1;
    }

    pub fn dead(self: *Creep) bool {
        return !self.alive;
    }

    // TODO: I suck at game programming... how bad is this...?
    pub fn update(self: *Creep, delta: u64) void {
        if (self.completed()) {
            return;
        }

        if (self.pos.position().col == self.cols - 1) {
            return;
        }

        const to = Position.fromIdx(self.path[self.pathIdx], self.cols);
        const dist = self.pos.subP(to);
        const normDist = dist.norm();
        const len = dist.len();

        const deltaF: f64 = @floatFromInt(delta);
        const deltaP: f64 = min(deltaF / 1000.0 * self.speed, len);
        const change = normDist.scale(-deltaP);
        self.pos = self.pos.add(change);

        if (self.pos.eql(to.vec2())) {
            self.pathIdx += 1;
        }
    }

    fn render(self: *Creep) void {
        self.rPos = self.pos.position();
    }
};

const t = std.testing;
const testBoaordSize = 3;
const testBoard = [9]bool{
    true, false, true,
    true, false, true,
    true, true, true,
};
test "bfs" {
    var seen: [9]isize = .{-1} ** 9;

    const out = walk(0, 0, 3, &testBoard, &seen);
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

fn runUntil(creep: *Creep, to: usize, maxRun: usize) void {
    var elapsed: usize = 0;
    while (elapsed < maxRun and creep.pathIdx != to) : (elapsed += 16) {
        creep.update(16);
    }
}

test "creep movement" {
    var creep = try Creep.init(t.allocator, 3, 3);
    defer creep.deinit();
    _ = creep.
        initialPosition(.{.row = 0, .col = 0}).
        calculatePath(&testBoard, testBoaordSize);

    try t.expect(creep.pathIdx == 0);
    runUntil(&creep, 1, 1500);
    try t.expect(creep.pos.x == 0.0);
    try t.expect(creep.pos.y == 1.0);
    try t.expect(creep.pathIdx == 1);
    runUntil(&creep, 2, 1500);
    try t.expect(creep.pathIdx == 2);
    try t.expect(creep.pos.x == 0.0);
    try t.expect(creep.pos.y == 2.0);
    runUntil(&creep, 3, 1500);
    try t.expect(creep.pathIdx == 3);
    try t.expect(creep.pos.x == 1.0);
    try t.expect(creep.pos.y == 2.0);
    runUntil(&creep, 4, 1500);
    try t.expect(creep.pathIdx == 4);
    try t.expect(creep.pos.x == 2.0);
    try t.expect(creep.pos.y == 2.0);
}
