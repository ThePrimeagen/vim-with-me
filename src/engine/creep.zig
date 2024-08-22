const std = @import("std");
const math = @import("../math/math.zig");
const objects = @import("../objects/objects.zig");
const framer = @import("framer.zig");
const a = @import("../assert/assert.zig");
const assert = a.assert;
const u = a.u;
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const debug = @import("debug.zig");

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
    creep.aabb = .{.min = pos, .max = pos.add(.{.x = 1, .y = 1})};
    creep.rSized = .{
        .cols = CreepSize,
        .pos = creep.pos.position(),
    };

    return creep;
}

fn path(parents: []const isize, pos: usize, out: []usize) usize {
    var idx: usize = 0;

    var p = pos;
    while (parents[p] != p) {
        out[idx] = p;
        p = @intCast(parents[p]);

        idx += 1;
    }

    std.mem.reverse(usize, out[0..idx]);
    return idx;
}

fn walk(from: usize, pos: usize, values: *const Values, board: []const bool, parents: []isize, seen: []bool) !usize {
    const cols = values.cols;
    assert(board.len % cols == 0, "board is not a rectangle");
    assert(board.len == parents.len, "board and parents should have the same size");
    assert(board.len == seen.len, "board and seen should have the same size");

    if (pos >= parents.len or parents[pos] != -1 or board[pos] == false) {
        return 0;
    }

    parents[pos] = @intCast(from);
    seen[pos] = true;

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

        const out = try walk(pos, @intCast(iPos + d), values, board, parents, seen);
        if (out != 0) {
            return out;
        }
    }

    return 0;
}

fn printWalk(at: usize, dists: []isize, cols: usize) void {
    for (0..dists.len) |idx| {
        if (idx % cols == 0 and idx > 0) {
            std.debug.print("\n", .{});
        }
        if (at == idx) {
            std.debug.print("x", .{});
        } else {
            std.debug.print("{} ", .{dists[idx]});
        }
    }
    std.debug.print("\n", .{});
    std.debug.print("\n", .{});
    std.debug.print("\n", .{});
}

pub fn walkAStar(start: usize, p: usize, v: *const Values, board: []const bool, parents: []isize, alloc: Allocator) !usize {
    var pos = p;
    const cols = v.cols;

    assert(board.len % cols == 0, "board is not a rectangle");
    assert(board.len == parents.len, "board and seen should have the same size");

    const dists: []isize = try alloc.alloc(isize, board.len);
    const seen: []bool = try alloc.alloc(bool, board.len);
    defer alloc.free(dists);
    defer alloc.free(seen);

    for (0..dists.len) |idx| {
        dists[idx] = -1;
        seen[idx] = false;
    }

    const iCols: isize = @intCast(cols);
    dists[pos] = 0;
    seen[pos] = true;
    parents[pos] = @intCast(start);

    while (true) {
        const directions: [4]isize = .{1, -iCols, iCols, -1};
        const iPos: isize = @intCast(pos);

        // activate new position
        for (directions) |d| {
            const nextPos: isize = iPos + d;

            if (nextPos < 0 or
                nextPos >= board.len) {
                continue;
            }

            const nextPosU: usize = @intCast(nextPos);
            if (
                board[@intCast(nextPos)] == false or
                parents[@intCast(nextPos)] != -1 or
                seen[@intCast(nextPos)]) {
                continue;
            }

            const isizeCols: isize = @intCast(cols);
            const posDist = @abs(@divFloor(nextPos, isizeCols) - @divFloor(iPos, isizeCols)) +
                @abs(@mod(nextPos, isizeCols) - @mod(iPos, isizeCols));

            if (posDist > 1) {
                continue;
            }

            //std.debug.print("parents[{}] = {}\n", .{nextPosU, iPos});
            const dist: isize = @intCast(cols - nextPosU % cols);
            const city: isize = @intCast(math.getCityDistanceFromIdx(start, nextPosU, cols));
            dists[nextPosU] = dist * 2 + city;
            parents[nextPosU] = iPos;
        }

        var lowest: isize = 255;
        var lowestIdx: isize = 0;
        for (0..dists.len) |idx| {
            if (dists[idx] < lowest and dists[idx] != -1 and parents[idx] != -1 and seen[idx] == false) {
                lowest = dists[idx];
                lowestIdx = @intCast(idx);
            }
        }

        if (lowest == 255) {
            std.debug.print("breaking because we couldn't find anything\n", .{});
            break;
        }

        pos = @intCast(lowestIdx);
        seen[pos] = true;

        // I am at the end
        if (pos % cols == cols - 1) {
            return pos;
        }

    }

    return 0;
}

pub fn walkBFS(start: usize, p: usize, v: *const Values, board: []const bool, parents: []isize, alloc: Allocator) !usize {
    var pos = p;
    const cols = v.cols;

    assert(board.len % cols == 0, "board is not a rectangle");
    assert(board.len == parents.len, "board and seen should have the same size");

    const seen: []bool = try alloc.alloc(bool, board.len);
    defer alloc.free(seen);

    const iCols: isize = @intCast(cols);
    seen[pos] = true;
    parents[pos] = @intCast(start);

    const PosList = std.ArrayList(usize);
    var queue = PosList.init(alloc);
    defer queue.deinit();
    try queue.append(pos);

    while (queue.items.len > 0) {

        pos = queue.orderedRemove(0);
        seen[pos] = true;

        const directions: [4]isize = .{1, -iCols, iCols, -1};
        const iPos: isize = @intCast(pos);

        // activate new position
        for (directions) |d| {
            const nextPos: isize = iPos + d;

            if (nextPos < 0 or
                nextPos >= board.len) {
                continue;
            }

            const nextPosU: usize = @intCast(nextPos);
            if (
                board[@intCast(nextPos)] == false or
                parents[@intCast(nextPos)] != -1 or
                seen[@intCast(nextPos)]) {
                continue;
            }

            const isizeCols: isize = @intCast(cols);
            const posDist = @abs(@divFloor(nextPos, isizeCols) - @divFloor(iPos, isizeCols)) +
                @abs(@mod(nextPos, isizeCols) - @mod(iPos, isizeCols));

            if (posDist > 1) {
                continue;
            }

            parents[nextPosU] = iPos;
            try queue.append(nextPosU);
        }

        // I am at the end
        if (pos % cols == cols - 1) {
            return pos;
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

pub fn calculatePath(self: *Creep, board: []const bool) !void {
    if (!self.alive) {
        return;
    }

    const cols = self.values.cols;
    assert(board.len % cols == 0, "the length is not a rectangle");
    assert(board.len == self.path.len, "the length of the board is different than what the creep was initialized with.");
    assert(self.alive, "cannot calculate a path for a dead creep");

    for (0..self.scratch.len) |idx| {
        self.scratch[idx] = -1;
    }

    const pos = self.aabb.min.position().toIdx(cols);
    var seen: [10000]bool = undefined;
    for (0..board.len) |idx| {
        seen[idx] = false;
    }

    //const last = a.unwrap(usize, walk(pos, pos, self.values, board, self.scratch[0..board.len], seen[0..board.len]));
    //const last = a.unwrap(usize, walkBFS(pos, pos, self.values, board, self.scratch, self.alloc));
    const last = a.unwrap(usize, walkAStar(pos, pos, self.values, board, self.scratch, self.alloc));

    if (last == 0) {
        a.never("unable to move creep forward");
    } else {
        self.pathLen = path(self.scratch, last, self.path);
        self.pathIdx = 0;
    }
}

pub fn kill(self: *Creep, gs: *GS) void {
    assert(gs.fns != null, "must have functions defined");

    self.alive = false;
    gs.fns.?.creepKilled(gs, self);
}

pub fn completed(self: *Creep) bool {
    return self.pos.position().col == self.values.cols - 1;
}

pub fn dead(self: *Creep) bool {
    return !self.alive;
}

pub fn scale(self: *Creep, round: usize) void {
    const speedAddition: f64 = @as(f64, @floatFromInt(@divFloor(round, self.values.creep.scaleSpeedRounds))) * self.values.creep.scaleSpeed;
    const lifeAddition = @divFloor(round, self.values.creep.scaleLifeRounds) * self.values.creep.scaleLife;

    self.life += lifeAddition;
    self.speed += speedAddition;
}

// TODO: I suck at game programming... how bad is this...?
pub fn update(self: *Creep, gs: *GS) void {
    if (completed(self) or !self.alive) {
        return;
    }

    if (self.life == 0) {
        kill(self, gs);
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
            self.aabb = self.aabb.move(self.pos);
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
const testBoard = [15]bool{
    true, false, true,
    true, false, true,
    true, false, true,
    true, false, true,
    true, true, true,
};
var testValues = blk: {
    var values = objects.Values{
        .rows = 5,
        .cols = 3,
    };
    values.init();
    break :blk values;
};

test "bfs" {
    var parents: [15]isize = .{-1} ** 15;
    var seen: [15]bool = .{false} ** 15;

    const out = try walk(0, 0, &testValues, &testBoard, &parents, &seen);
    try t.expect(out == 14);

    var expected: [15]isize = .{
        0, -1, -1,
        0, -1, -1,
        3, -1, -1,
        6, -1, -1,
        9, 12, 13,
    };

    try t.expectEqualSlices(isize, &expected, &parents);

    var p: [9]usize = .{0} ** 9;
    const len = path(&parents, out, &p);
    var pExpect = [_]usize{ 3, 6, 9, 12, 13, 14 };

    try t.expect(len == 6);
    try t.expectEqualSlices(usize, &pExpect, p[0..len]);
}

fn runUntil(creep: *Creep, gs: *GS, to: usize, maxRun: usize) void {
    var elapsed: usize = 0;
    while (elapsed < maxRun and creep.pathIdx != to) : (elapsed += 16) {
        update(creep, gs);
    }
}

test "creep contains" {
    var creep = try create(t.allocator, 0, Values.TEAM_ONE, &testValues, .{.y = 1, .x = 1});
    defer creep.deinit();

    try t.expect(!contains(&creep, .{.x = 0.9999, .y = 0}));
    try t.expect(!contains(&creep, .{.x = 1.5, .y = 0}));
    try t.expect(!contains(&creep, .{.x = 0, .y = 1.5}));
    try t.expect(contains(&creep, .{.x = 1, .y = 1.5}));

}

test "creep walk astar" {
    var gs = try GS.init(t.allocator, &testValues);
    defer gs.deinit();

    var creep = try create(t.allocator, 0, Values.TEAM_ONE, &testValues, .{.x = 0, .y = 0});
    defer creep.deinit();
    _ = try calculatePath(&creep, &testBoard);
}
