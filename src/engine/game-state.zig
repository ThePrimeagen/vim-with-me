const std = @import("std");
const types = @import("types.zig");

// TODO: Make this adjustable
const Position = types.Vec2;
const Message = types.Message;
const Coord = types.Coord;
const assert = @import("assert").assert;

const ArrayList = std.ArrayList;
const Allocator = std.mem.Allocator;

const Tower = types.Tower;
const Creep = types.Creep;
const Projectile = types.Projectile;

const TowerList = ArrayList(types.Tower);
const CreepList = ArrayList(types.Creep);
const ProjectileList = ArrayList(types.Projectile);

pub const GameState = struct {
    playing: bool = true,
    round: i32 = 1,
    one: i32 = 0,
    two: i32 = 0,

    time: i64 = 0,
    loopDelta: i64 = 0,
    updates: i64 = 0,

    towers: TowerList,
    creeps: CreepList,
    projectile: ProjectileList,
    allocator: Allocator,

    pub fn init(alloc: Allocator) GameState {
        return .{
            .towers = TowerList.init(alloc),
            .creeps = CreepList.init(alloc),
            .projectile = ProjectileList.init(alloc),
            .allocator = alloc,
        };
    }

    pub fn update(state: *GameState, delta: i64) void {
        state.updates += 1;

        const diff: isize = @intCast(state.one - state.two);
        assert(diff >= -1 and diff <= 1, "some how we have multiple updates to one side but not the other");

        state.loopDelta = delta;
        state.time += delta;

        if (state.playing) {
            state.runUpdate();
        }
    }

    pub fn play(state: *GameState) void {
        assert(state.one == state.two, "player one and two must have same play count");
        state.playing = true;
    }

    pub fn pause(state: *GameState) void {
        assert(state.one == state.two, "player one and two must have same play count");
        state.playing = false;
    }

    pub fn message(state: *GameState, msg: Message) !void {
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

    pub fn clone(self: *GameState) GameState {
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

    fn tower(self: *GameState, pos: Position) ?usize {
        for (self.towers.items, 0..) |*t, i| {
            if (t.contains(pos)) {
                return i;
            }
        }
        return null;
    }

    fn creep(self: *GameState, pos: Position) ?usize {
        for (self.creeps.items, 0..) |*c, i| {
            if (c.contains(pos)) {
                return i;
            }
        }
        return null;
    }

    fn runUpdate(self: *GameState) void {
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
};

