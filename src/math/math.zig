const std = @import("std");
const assert = @import("../assert/assert.zig");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;

pub const ZERO_POS: Position = .{.row = 0, .col = 0};
pub const ZERO_VEC2: Vec2 = .{.x = 0.0, .y = 0.0};
pub const ZERO_AABB: AABB = .{.min = ZERO_VEC2, .max = ZERO_VEC2};
pub const ZERO_SIZED: Sized = .{.cols = 3, .pos = ZERO_POS};

pub fn floor(a: f64, precision: usize) f64 {
    const p: f64 = @floatFromInt(precision);
    return std.math.floor(a * p) / p;
}

pub fn min(a: f64, b: f64) f64 {
    return if (a > b) b else a;
}

pub fn usizeToIsizeSub(a: usize, b: usize) isize {
    const ai: isize = @intCast(a);
    const bi: isize = @intCast(b);
    return ai - bi;
}

pub fn getCityDistanceFromIdx(a: usize, b: usize, cols: usize) usize {
    const ai: isize = @intCast(a);
    const bi: isize = @intCast(b);
    const colsi: isize = @intCast(cols);

    const rowDiff = @abs(@divFloor(ai, colsi) - @divFloor(bi, colsi));
    const colDiff = @abs(@mod(ai, colsi) - @mod(bi, colsi));

    return  rowDiff + colDiff;
}

pub fn absUsize(a: usize, b: usize) usize {
    if (a > b) {
        return a - b;
    }
    return b - a;
}

const needle: [1]u8 = .{','};

pub const Position = struct {
    row: usize,
    col: usize,

    pub fn vec2(self: Position) Vec2 {
        return .{
            .x = @floatFromInt(self.col),
            .y = @floatFromInt(self.row),
        };
    }

    pub fn toIdx(self: Position, cols: usize) usize {
        return self.row * cols + self.col;
    }

    pub fn add(self: Position, other: Position) Position {
        return .{
            .row = self.row + other.row,
            .col = self.col + other.col,
        };
    }

    pub fn fromIdx(idx: usize, cols: usize) Position {
        return .{
            .row = idx / cols,
            .col = idx % cols,
        };
    }

    pub fn init(str: []const u8, pos: *?Position) usize {
        const idx = std.mem.indexOf(u8, str, ",") orelse return 0;
        if (str.len <= idx + 1) {
            return 0;
        }

        const endIdx = idx + 1 + (std.mem.indexOf(u8, str[idx + 1..], ",") orelse str[idx + 1..].len);

        assert.unwrap(void, std.io.getStdErr().writeAll(assert.u(std.fmt.bufPrint(scratchBuf(150), "str = {s} idx = {}, endIdx = {}", .{str, idx, endIdx}))));
        const row = std.fmt.parseInt(usize, str[0..idx], 10) catch {
            return 0;
        };
        const col = std.fmt.parseInt(usize, str[idx + 1..endIdx], 10) catch {
            return 0;
        };

        pos.* = .{
            .row = row,
            .col = col,
        };

        return endIdx;
    }

    pub fn string(self: Position) ![]u8 {
        const tmp = scratchBuf(50);
        return try std.fmt.bufPrint(tmp, "vec(r = {}, c = {})", .{self.row, self.col});
    }
};

pub const Range = struct {
    startRow: usize = 0,
    endRow: usize = 0,

    pub fn string(self: Range) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(100), "Range{{s={}, e={}}}", .{self.startRow, self.endRow});
    }

    pub fn len(self: Range) usize {
        if (self.endRow < self.startRow) {
            return 0;
        }
        return self.endRow - self.startRow;
    }

    pub fn position(self: Range) Position {
        return .{
            .row = self.startRow,
            .col = 0,
        };
    }

    pub fn sized(self: Range, cols: usize) Sized {
        return .{
            .pos = self.position(),
            .cols = cols,
        };
    }

    pub fn invalid(self: Range) bool {
        return self.startRow >= self.endRow;
    }

    pub fn contains(self: Range, pos: Position) bool {
        return self.startRow <= pos.row and self.endRow > pos.row;
    }

    pub fn containsAABB(self: Range, aabb: AABB) bool {
        const mY: usize = @intFromFloat(aabb.min.y);
        const xY: usize = @intFromFloat(aabb.max.y);
        return self.startRow <= mY and self.endRow >= xY;
    }
};

pub const Sized = struct {
    cols: usize,
    pos: Position,

    pub fn string(self: Sized) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(70), "Sized(cols={}, pos={s})", .{
            self.cols,
            try self.pos.string(),
        });
    }

    pub fn add(self: Sized, pos: Position) Sized {
        var s = self;
        s.pos = s.pos.add(pos);
        return s;
    }
};

pub const PossiblePositions = struct {
    positions: [50]?Position,
    len: usize,
    team: u8,

    pub fn init(msg: []const u8) ?PossiblePositions {
        if (msg.len == 0 or msg[0] != '1' and msg[0] != '2') {
            return null;
        }

        const team = msg[0];
        var posSet: PossiblePositions = .{
            .team = team,
            .len = 0,
            .positions = undefined,
        };

        var curr = msg[1..];
        while (true) {
            var pos: ?Position = null;
            const nextOffset = Position.init(curr, &pos);

            if (nextOffset == 0) {
                break;
            }

            posSet.positions[posSet.len] = pos;
            posSet.len += 1;

            if (curr.len <= nextOffset + 1) {
                break;
            }

            curr = curr[nextOffset + 1..];
        }

        if (posSet.len == 0) {
            return null;
        }

        return posSet;
    }

    pub fn empty() PossiblePositions {
        return .{
            .positions = undefined,
            .len = 0,
            .team = 0,
        };
    }
};


pub const AABB = struct {
    min: Vec2,
    max: Vec2,

    pub fn aabb(m: Vec2, max: Vec2) AABB {
        return .{
            .min = m,
            .max = max,
        };
    }

    pub fn contains(self: AABB, pos: Vec2) bool {
        return pos.x >= self.min.x and pos.x < self.max.x and
            pos.y >= self.min.y and pos.y < self.max.y;
    }

    pub fn containsAABB(self: AABB, other: AABB) bool {
        return self.contains(other.min) and self.contains(other.max);
    }

    pub fn overlaps(self: AABB, other: AABB) bool {
        return !(
            self.min.x >= other.max.x or
            self.min.y >= other.max.y or
            other.min.x >= self.max.x or
            other.min.y >= self.max.y
        );
    }

    pub fn add(self: AABB, to: Vec2) AABB {
        var out = self;
        out.min = out.min.add(to);
        out.max = out.max.add(to);
        return out;
    }

    pub fn move(self: AABB, to: Vec2) AABB {
        const xDist = self.max.x - self.min.x;
        const yDist = self.max.y - self.min.y;
        return .{
            .min = to,
            .max = to.add(.{.x = xDist, .y = yDist }),
        };
    }

    pub fn string(self: AABB) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(100), "AABB({s}, {s})", .{try self.min.string(), try self.max.string()});
    }

};

pub const Vec2 = struct {
    x: f64,
    y: f64,

    pub fn eql(self: Vec2, b: Vec2) bool {
        return self.x == b.x and self.y == b.y;
    }

    pub fn closeEnough(self: Vec2, b: Vec2, enough: f64) bool {
        return @abs(self.x - b.x) < enough and
            @abs(self.y - b.y) < enough;
    }

    pub fn aabb(self: Vec2, max: Vec2) AABB {
        return AABB.aabb(self, .{
            .x = self.x + max.x,
            .y = self.y + max.y,
        });
    }

    pub fn norm(self: Vec2) Vec2 {
        const l = self.len();
        return .{
            .x = self.x / l,
            .y = self.y / l,
        };
    }

    pub fn subP(self: Vec2, b: Position) Vec2 {
        const rf: f64 = @floatFromInt(b.row);
        const cf: f64 = @floatFromInt(b.col);

        return .{
            .x = self.x - cf,
            .y = self.y - rf,
        };
    }

    pub fn row(self: Vec2) isize {
        return @intFromFloat(self.y);
    }

    pub fn sub(self: Vec2, b: Vec2) Vec2 {
        return .{
            .x = self.x - b.x,
            .y = self.y - b.y,
        };
    }

    pub fn len(self: Vec2) f64 {
        return std.math.sqrt(self.x * self.x + self.y * self.y);
    }

    pub fn lenSq(self: Vec2) f64 {
        return self.x * self.x + self.y * self.y;
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

    pub fn position(self: *const Vec2) Position {
        assert.assert(self.x >= 0, "x cannot be negative");
        assert.assert(self.y >= 0, "y cannot be negative");

        return .{
            .row = @intFromFloat(self.y),
            .col = @intFromFloat(self.x),
        };
    }

    pub fn fromPosition(pos: Position) Vec2 {
        return .{
            .y = @floatFromInt(pos.row),
            .x = @floatFromInt(pos.col),
        };
    }

    pub fn string(self: Vec2) ![]u8 {
        return try std.fmt.bufPrint(scratchBuf(50), "x = {}, y = {}", .{floor(self.x, 1000), floor(self.y, 1000)});
    }

};

const t = std.testing;
test "vec2 add" {
    const a = Vec2{.x = 1, .y = 2};
    const b = Vec2{.x = 68, .y = 418};
    try t.expect(a.add(b).eql(Vec2{.x = 69, .y = 420}));
}

test "aabb contains" {
    const a = Vec2{.x = 1, .y = 1};
    const box = a.aabb(.{.x = 3, .y = 3});

    try t.expect(!box.contains(.{.x = 1, .y = 0.9999}));
    try t.expect(!box.contains(.{.x = 0.9999, .y = 1}));
    try t.expect(box.contains(.{.x = 1, .y = 1}));
    try t.expect(box.contains(.{.x = 3.9999, .y = 1}));
    try t.expect(box.contains(.{.x = 1, .y = 3.9999}));

    try t.expect(!box.contains(.{.x = 3.9999, .y = 4}));
    try t.expect(!box.contains(.{.x = 4, .y = 3.9999}));

}

test "aabb overlap" {
    const a = AABB{.min = .{.x = 1, .y = 1}, .max = .{.x = 3, .y = 3}};
    const b_false_x = AABB{.min = .{.x = 3, .y = 1}, .max = .{.x = 4, .y = 4}};
    const b_false_y = AABB{.min = .{.x = 1, .y = 3}, .max = .{.x = 4, .y = 4}};
    const b_true = AABB{.min = .{.x = 2, .y = 1}, .max = .{.x = 4, .y = 4}};

    try t.expect(!a.overlaps(b_false_x));
    try t.expect(!a.overlaps(b_false_y));
    try t.expect(a.overlaps(b_true));

    try t.expect(!b_false_x.overlaps(a));
    try t.expect(!b_false_y.overlaps(a));
    try t.expect(b_true.overlaps(a));
}

test "aabb collision failure creep tower issue" {
    const tower = AABB{.min = .{.x = 2.8e1, .y = 0e0}, .max = .{.x = 3.1e1, .y = 1e0}};
    const creep = AABB{.min = .{.x = 3e1, .y = 1e0}, .max = .{.x = 3.1e1, .y = 2e0}};

    try t.expect(!tower.overlaps(creep));
    try t.expect(!creep.overlaps(tower));
}

