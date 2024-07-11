const std = @import("std");

const objects = @import("../objects/objects.zig");
const math = @import("../math/math.zig");
const a = @import("../assert/assert.zig");
const assert = a.assert;
const towers = @import("tower.zig");
const creeps = @import("creep.zig");
const projectiles = @import("projectile.zig");

const GS = objects.gamestate.GameState;
const Message = objects.message.Message;
const Tower = objects.tower.Tower;
const Creep = objects.creep.Creep;
const Projectile = objects.projectile.Projectile;
const Vec2 = math.Vec2;
const Allocator = std.mem.Allocator;

pub fn update(state: *GS, delta: i64) !void {
    state.updates += 1;

    const changed = state.boardChanged;

    const diff: isize = @intCast(state.one - state.two);
    assert(diff >= -1 and diff <= 1, "some how we have multiple updates to one side but not the other");

    state.loopDeltaUS = delta;

    if (!state.playing) {
        return;
    }

    state.time += delta;

    for (state.towers.items) |*t| {
        try towers.update(t, state);
    }

    for (state.creeps.items) |*c| {
        creeps.update(c, state);
    }

    for (state.projectile.items) |*p| {
        try projectiles.update(p, state);
    }

    if (state.boardChanged > changed) {
        updateBoard(state);
        for (state.creeps.items) |*c| {
            creeps.calculatePath(c, state.board);
        }
    }

    if (!state.values.debug) {
        return;
    }

    if (state.towers.items.len < 2 or state.creeps.items.len < 1) {
        return;
    }

    const one = &state.towers.items[1];
    const c = state.creeps.items[0];

    std.debug.print("within: {} -- creep: {s} tower: {s}\n", .{
        towers.withinRange(one, c.pos),
        a.u(c.pos.string()),
        a.u(one.pos.string()),
    });
}

pub fn init(self: *GS) void {
    self.fns = &.{
        .placeProjectile = placeProjectile,
        .towerDied = towerDied,
        .creepKilled = creepKilled,
        .shot = shot,
        .strike = strike,
    };

    const rows = self.values.rows;
    const teamSpace = rows / 3;

    self.oneRange.endRow = teamSpace;

    self.twoRange.startRow = rows - teamSpace;
    self.twoRange.endRow = rows;

    self.noBuildRange.startRow = self.oneRange.endRow;
    self.noBuildRange.endRow = self.twoRange.startRow;

    self.playing = false;
}

pub fn towerDied(self: *GS, t: *Tower) void {
    if (t.team == objects.Values.TEAM_ONE) {
        self.oneStats.towersLost += 1;
    } else {
        self.twoStats.towersLost += 1;
    }
    self.boardChanged += 1;
}

pub fn creepKilled(self: *GS, c: *Creep) void {
    if (c.team == objects.Values.TEAM_ONE) {
        self.oneStats.creepsKilled += 1;
    } else {
        self.twoStats.creepsKilled += 1;
    }
}

pub fn shot(self: *GS, t: *Tower) void {
    if (t.team == objects.Values.TEAM_ONE) {
        self.oneStats.shots += 1;
    } else {
        self.twoStats.shots += 1;
    }
}

pub fn strike(self: *GS, p: *Projectile) void {
    switch (p.target) {
        .creep => |c| self.creeps.items[c].life -|= p.damage,
        .tower => |t| self.towers.items[t].ammo -|= p.damage,
    }
}

pub fn play(state: *GS) void {
    assert(state.one == state.two, "player one and two must have same play count");
    state.playing = true;
}

pub fn pause(state: *GS) void {
    assert(state.one == state.two, "player one and two must have same play count");
    state.playing = false;
}

pub fn message(state: *GS, msg: Message) !void {
    switch (msg) {
        .coord => |c| {

            //if (c.team == '1') {
            //    state.one += 1;
            //} else {
            //    state.two += 1;
            //}

            state.one += 1;
            state.two += 1;

            if (tower(state, c.pos.vec2())) |idx| {
                towers.upgrade(&state.towers.items[idx]);
                return;
            }

            //if (canPlaceTower(state, c.pos)) {
            //    const tt = towers.TowerBuilder.start().
            //        team(c.team).
            //        pos(c.pos).
            //        id(state.towers.items.len).
            //        tower(state);

            //    try state.towers.append(tt);
            //} else {
                // TODO: randomly place tower
                unreachable;
            //}

        },
        .round => |_| {
            // not sure what to do here...
            // probably need to think about "playing/pausing"
            // play(state);
        },
    }
}

pub fn clone(self: *GS) !GS {
    const diff: isize = @intCast(self.one - self.two);
    assert(diff == 0, "next round can only be called once both players have played their turns.");

    var board = try self.alloc.alloc(bool, self.board.len);
    std.mem.copyForwards(bool, board[0..], self.board);

    return .{
        .round = self.round,
        .values = self.values,

        .one = self.one,
        .oneCoords = self.oneCoords,

        .two = self.two,
        .twoCoords = self.twoCoords,

        .time = self.time,
        .loopDeltaUS = self.time,

        .towers = try self.towers.clone(),
        .creeps = try self.creeps.clone(),
        .projectile = try self.projectile.clone(),
        .board = board,
        .alloc = self.alloc,
    };
}

fn tower(self: *GS, pos: Vec2) ?usize {
    for (self.towers.items, 0..) |*t, i| {
        if (towers.contains(t, pos)) {
            return i;
        }
    }
    return null;
}

fn creep(self: *GS, pos: Vec2) ?usize {
    for (self.creeps.items, 0..) |*c, i| {
        if (creeps.contains(c, pos)) {
            return i;
        }
    }
    return null;
}

pub fn calculateBoard(self: *GS) void {
    for (self.board, 0..) |_, idx| {
        self.board[idx] = true;
    }

    for (self.towers.items) |*t| {
        const cells = t.rCells;
        const sized = t.rSized;

        for (cells, 0..) |_, idx| {
            const col = idx % sized.cols;
            const row = idx / sized.cols;
            const offset = (sized.pos.row + row) * self.values.cols + sized.pos.col + col;
            self.board[offset] = true;
        }
    }
}

pub fn placeCreep(self: *GS, pos: math.Position, team: u8) !usize {
    switch (team) {
        objects.Values.TEAM_ONE => assert(self.oneRange.contains(pos), "invalid team one position"),
        objects.Values.TEAM_TWO => assert(self.twoRange.contains(pos), "invalid team one position"),
        else => a.never("invalid team"),
    }

    const id = self.creeps.items.len;
    var c = try creeps.create(
        self.alloc, id, team, self.values, pos.vec2()
    );

    errdefer c.deinit();
    try self.creeps.append(c);

    creeps.calculatePath(&self.creeps.items[id], self.board);

    return id;
}

pub fn updateBoard(self: *GS) void {
    for (0..self.board.len) |i| {
        self.board[i] = true;
    }

    for (self.towers.items) |*t| {
        if (!t.alive) {
            continue;
        }

        const start = t.rSized.pos.row * self.values.cols + t.rSized.pos.col;

        for (0..t.rows) |r| {
            const rowStart = start + r * self.values.cols;
            for (0..t.cols) |c| {
                self.board[rowStart + c] = false;
            }
        }
    }
}

fn canPlaceTower(self: *GS, pos: math.Position, team: u8) bool {

    if (self.noBuildZone) {
        const range = switch (team) {
            '1' => self.oneRange,
            '2' => self.twoRange,
            else => {
                assert(false, "inTeam is an invalid value");
                unreachable;
            }
        };

        if (!range.contains(pos)) {
            return false;
        }
    }

    if (pos.col == 0 or pos.col == self.values.cols - 1) {
        return false;
    }

    if (tower(self, pos.vec2())) |_| {
        return false;
    }

    if (creep(self, pos.vec2())) |_| {
        return false;
    }

    return true;
}

pub fn placeTower(self: *GS, pos: math.Position, team: u8) !?usize {
    if (!canPlaceTower(self, pos, team)) {
        return null;
    }

    const id = self.towers.items.len;
    const t = towers.TowerBuilder.start()
        .pos(pos)
        .team(team)
        .id(id)
        .tower(self);

    try self.towers.append(t);

    updateBoard(self);

    for (self.creeps.items) |*c| {
        creeps.calculatePath(c, self.board);
    }

    return id;
}

pub fn placeProjectile(self: *GS, t: *Tower, target: objects.Target) Allocator.Error!usize {
    const id = self.projectile.items.len;
    const projectile = objects.projectile.Projectile {
        .pos = t.pos,
        .target = target,
        .id = id,
        .damage = t.damage,
        .speed = self.values.projectile.speed,
    };

    try self.projectile.append(projectile);
    shot(self, t);

    return id;
}

pub fn towerById(self: *GS, id: usize) *Tower {
    assert(self.towers.items.len > id, "grabbing a tower outside the size of the tower list");

    const t = &self.towers.items[id];
    assert(t.alive, "cannot retrieve a dead tower");
    return t;
}


pub fn validateState(self: *GS) void {
    for (self.creeps.items) |*c| {
        if (tower(self, c.pos)) |t| {
            std.debug.print("tower: {s} collided with creep {s}\n", .{a.u(self.towers.items[t].pos.string()), a.u(c.string())});
            assert(false, "a creep is within a tower");
        }
    }
}

const testing = std.testing;
test "calculate the board" {
    var values = objects.Values{.rows = 3, .cols = 3};
    values.init();

    var gs = try GS.init(testing.allocator, &values);
    defer gs.deinit();

    calculateBoard(&gs);

    try testing.expectEqualSlices(bool, &.{
        true, true, true,
        true, true, true,
        true, true, true,
    }, gs.board);
}

test "place creep calculates positions" {
    var values = objects.Values{.rows = 3, .cols = 3};
    values.init();

    var gs = try GS.init(testing.allocator, &values);
    defer gs.deinit();
    calculateBoard(&gs);
    init(&gs);

    _ = try placeCreep(&gs, .{.row = 0, .col = 0}, objects.Values.TEAM_ONE);
    try testing.expect(gs.creeps.items[0].pathLen == 2);
}
