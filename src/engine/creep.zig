const std = @import("std");
const math = @import("../math/math.zig");
const objects = @import("../objects/objects.zig");
const a = @import("../assert/assert.zig");
const assert = a.assert;
const u = a.u;
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;

const Values = objects.Values;
const Allocator = std.mem.Allocator;
const Creep = objects.creep.Creep;
const CreepSize = objects.creep.CreepSize;
const colors = objects.colors;
const GS = objects.gamestate.GameState;
const Position = math.Position;
const MICROSECOND = 1_000_000.0;

pub fn distanceToExit(creep: *Creep, gs: *GS) f64 {
    assert(creep.alive, "you cannot call distance to exit if the creep is dead");
    assert(!completed(creep), "expected the creep to be still within the maze");
    assert(creep.pathLen > creep.pathIdx, "pathLen should always be larger than pathIdx");

    // 1 or larger
    const diff: f64 = @floatFromInt(creep.pathLen - creep.pathIdx);

    // distance from where we are to the _NEXT_ path and add that
    const dist = Position.fromIdx(creep.path[creep.pathIdx], gs.values.cols).
        vec2().sub(creep.pos).lenSq();

    return diff + dist;
}

// TODO: Params object STAT (just not now)
pub fn create(alloc: Allocator, id: usize, team: u8, values: *const Values, pos: math.Vec2) !Creep {
    assert(team == Values.TEAM_ONE or team == Values.TEAM_TWO, "invalid team");

    var creep: Creep = try Creep.init(alloc, values);
    creep.id = id;
    creep.team = team;

    creep.pos = pos;
    creep.rSized = .{
        .cols = CreepSize,
        .pos = creep.pos.position(),
    };

    return creep;
}

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


pub fn contains(self: *Creep, pos: math.Vec2) bool {
    if (!self.alive) {
        return false;
    }

    // TODO: i should probably make this work with negative numbers
    const myPosition = self.pos.position();
    const inPos = pos.position();

    return myPosition.row == inPos.row and myPosition.col == inPos.col;
}

pub fn calculatePath(self: *Creep, board: []const bool) void {
    const cols = self.values.cols;
    assert(board.len % cols == 0, "the length is not a rectangle");
    assert(board.len == self.path.len, "the length of the board is different than what the creep was initialized with.");

    for (0..self.scratch.len) |idx| {
        self.scratch[idx] = -1;
    }

    const pos = self.pos.position().toIdx(cols);
    const last = walk(pos, pos, cols, board, self.scratch);

    if (last == 0) {
        unreachable;
    } else {
        self.pathLen = path(self.scratch, last, self.path);
        self.pathIdx = 0;
    }
}

pub fn completed(self: *Creep) bool {
    return self.pos.position().col == self.values.cols - 1;
}

pub fn dead(self: *Creep) bool {
    return !self.alive;
}

// TODO: I suck at game programming... how bad is this...?
pub fn update(self: *Creep, gs: *GS) void {
    if (completed(self) or !self.alive) {
        return;
    }

    if (self.life == 0) {
        self.alive = false;
        return;
    }

    var consumedUS: i64 = 0;
    var count: usize = 0;
    while (consumedUS < gs.loopDeltaUS and !completed(self)) {
        const delta = gs.loopDeltaUS - consumedUS;

        const to = math.Position.fromIdx(self.path[self.pathIdx], self.values.cols);
        const dist = self.pos.subP(to);
        const normDist = dist.norm();
        const len = dist.len();
        const maxUS: i64 = @intFromFloat(@ceil(len / self.speed * MICROSECOND));
        const usConsumed = @min(maxUS, delta);
        assert(usConsumed > 0, "we must consume some amount of the remaining microsecnds");

        const deltaF: f64 = @floatFromInt(usConsumed);
        const deltaP: f64 = deltaF / MICROSECOND * self.speed;
        const change = normDist.scale(-deltaP);
        self.pos = self.pos.add(change);

        if (self.pos.closeEnough(to.vec2(), 0.001)) {
            self.pos = Position.fromIdx(self.path[self.pathIdx], self.values.cols).vec2();
            self.pathIdx += 1;
        }

        consumedUS += usConsumed;
        count += 1;
    }
}

pub fn render(self: *Creep, gs: *GS) void {
    self.rSized.pos = self.pos.position();
    self.rCells[0].text = '0' + @as(u8, @intCast(self.life));
    _ = gs;
}

const t = std.testing;
const testBoaordSize = 3;
const testBoard = [9]bool{
    true, false, true,
    true, false, true,
    true, true, true,
};
const testEmptyBoard = [9]bool{
    true, true, true,
    true, true, true,
    true, true, true,
};
var testValues = blk: {
    var values = objects.Values{
        .rows = 3,
        .cols = 3,
    };
    values.init();
    break :blk values;
};

test "bfs" {
    testValues.init();
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

fn runUntil(creep: *Creep, gs: *GS, to: usize, maxRun: usize) void {
    var elapsed: usize = 0;
    while (elapsed < maxRun and creep.pathIdx != to) : (elapsed += 16) {
        update(creep, gs);
    }
}

test "creep movement" {
    var gs = try GS.init(t.allocator, &testValues);
    defer gs.deinit();

    var creep = try create(t.allocator, 0, Values.TEAM_ONE, &testValues, .{.x = 0, .y = 0});
    defer creep.deinit();
    _ = calculatePath(&creep, &testBoard);

    gs.loopDeltaUS = 16_000;

    // Note: with how update works we can "bleed" a bit over
    // I find that this amount captures any oopsy amount of movement.
    const closeEnough = 0.05;
    try t.expect(creep.pathIdx == 0);

    runUntil(&creep, &gs, 1, 1500);
    try t.expect(creep.pos.closeEnough(.{.x = 0, .y = 1}, closeEnough));
    try t.expect(creep.pathIdx == 1);

    runUntil(&creep, &gs, 2, 1500);
    try t.expect(creep.pathIdx == 2);
    try t.expect(creep.pos.closeEnough(.{.x = 0, .y = 2}, closeEnough));

    runUntil(&creep, &gs, 3, 1500);
    try t.expect(creep.pathIdx == 3);
    try t.expect(creep.pos.closeEnough(.{.x = 1, .y = 2}, closeEnough));

    runUntil(&creep, &gs, 4, 1500);
    try t.expect(creep.pathIdx == 4);
    try t.expect(creep.pos.closeEnough(.{.x = 2, .y = 2}, closeEnough));

}

test "creep contains" {
    var creep = try create(t.allocator, 0, Values.TEAM_ONE, &testValues, .{.y = 1, .x = 1});
    defer creep.deinit();

    try t.expect(!contains(&creep, .{.x = 0.9999, .y = 0}));
    try t.expect(!contains(&creep, .{.x = 1.5, .y = 0}));
    try t.expect(!contains(&creep, .{.x = 0, .y = 1.5}));
    try t.expect(contains(&creep, .{.x = 1, .y = 1.5}));

}
