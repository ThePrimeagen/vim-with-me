const objects = @import("objects");
const math = @import("math");
const assert = @import("assert").assert;
const towers = @import("tower.zig");

const GS = objects.gamestate.GameState;
const Message = objects.message.Message;
const Tower = objects.tower.Tower;
const Vec2 = math.Vec2;

pub fn update(state: *GS, delta: i64) void {
    state.updates += 1;

    const diff: isize = @intCast(state.one - state.two);
    assert(diff >= -1 and diff <= 1, "some how we have multiple updates to one side but not the other");

    state.loopDelta = delta;
    state.time += delta;

    if (state.playing) {
        state.runUpdate();
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

            if (state.tower(c.pos)) |idx| {
                state.towers.items[idx].level += 1;
                return;
            }

            // a tower may not be able to fit between two towers...
            // i may need to "fit" them in
            const id = state.towers.items.len;
            try state.towers.append(Tower.create(id, c.team, c.pos));

        },
        .round => |_| {
            state.nextRound();
        },
    }
}

pub fn clone(self: *GS) GS {
    const diff: isize = @intCast(self.one - self.two);
    assert(diff == 0, "next round can only be called once both players have played their turns.");

    return .{
        .round = self.round,
        .one = self.one,
        .two = self.two,
        .time = self.time,
        .loopDelta = self.time,

        .towers = self.towers.clone(),
        .creeps = self.creeps.clone(),
        .projectile = self.projectile.clone(),
        .allocator = self.allocator,
    };
}

fn tower(self: *GS, pos: Vec2) ?usize {
    for (self.towers.items, 0..) |*t, i| {
        if (towers.contains(t, ) {
            return i;
        }
    }
    return null;
}

fn creep(self: *GS, pos: Vec2) ?usize {
    for (self.creeps.items, 0..) |*c, i| {
        if (c.contains(pos.)) {
            return i;
        }
    }
    return null;
}

fn runUpdate(self: *GS) void {
    for (self.towers.items) |*t| {
        t.update();
    }
    for (self.creeps.items) |*c| {
        c.update();
    }
    for (self.projectile.items) |*p| {
        p.update();
    }
}

