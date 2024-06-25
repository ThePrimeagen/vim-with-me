const objects = @import("objects");
const math = @import("math");
const assert = @import("assert").assert;
const towers = @import("tower.zig");
const creeps = @import("creep.zig");

const GS = objects.gamestate.GameState;
const Message = objects.message.Message;
const Tower = objects.tower.Tower;
const Vec2 = math.Vec2;

pub fn update(state: *GS, delta: i64) void {
    state.updates += 1;

    const diff: isize = @intCast(state.one - state.two);
    assert(diff >= -1 and diff <= 1, "some how we have multiple updates to one side but not the other");

    state.loopDeltaUS = delta;
    state.time += delta;

    if (!state.playing) {
        return;
    }

    for (state.towers.items) |*t| {
        towers.update(t, state);
    }

    for (state.creeps.items) |*c| {
        creeps.update(c, state);
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
                state.towers.items[idx].level += 1;
                return;
            }

            // a tower may not be able to fit between two towers...
            // i may need to "fit" them in
            const tt = towers.TowerBuilder.start().team(c.team).pos(c.pos).tower();
            try state.towers.append(tt);

        },
        .round => |_| {
            // not sure what to do here...
            // probably need to think about "playing/pausing"
            // play(state);
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
        .loopDeltaUS = self.time,

        .towers = self.towers.clone(),
        .creeps = self.creeps.clone(),
        .projectile = self.projectile.clone(),
        .allocator = self.allocator,
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

// TODO: vec and position?
pub fn placeCreep(self: *GS, pos: math.Position) !void {
    const c = try creeps.create(self.alloc, self.creeps.items.len, 0, self.values, pos.vec2());
    try self.creeps.append(c);
}
