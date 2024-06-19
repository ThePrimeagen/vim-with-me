const std = @import("std");
const types = @import("types.zig");

const Position = types.Position;
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
    playing: bool,
    round: i32,
    one: i32,
    two: i32,

    time: i64,
    loopDelta: i64,
    updates: i64,

    towers: TowerList,
    creeps: CreepList,
    projectile: ProjectileList,
    allocator: Allocator,

    pub fn init(alloc: Allocator) GameState {
        return .{
            .playing = true,
            .round = 1,
            .one = 0,
            .two = 0,
            .time = 0,
            .loopDelta = 0,

            .towers = TowerList.init(alloc),
            .creeps = CreepList.init(alloc),
            .projectile = ProjectileList.init(alloc),
            .allocator = alloc,
        };
    }

    pub fn updateTime(state: *GameState, delta: u64) void {
        state.updates += 1;

        const diff: isize = @intCast(state.one - state.two);
        assert(diff >= -1 and diff <= 1, "some how we have multiple updates to one side but not the other");

        state.loopDelta = delta;
        state.time += delta;
    }

    pub fn message(state: *GameState, msg: Message) void {
        switch (msg) {
            .coord => |c| {
                if (c.team == '1') {
                    state.one += 1;
                } else {
                    state.two += 1;
                }
            },
            .round => |_| {
                state.nextRound();
            },
        }
    }

    pub fn nextRound(state: *GameState) void {
        const diff: isize = @intCast(state.one - state.two);
        assert(diff == 0, "next round can only be called once both players have played their turns.");

        state.playRound();
        state.round += 1;
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

    fn playRound(self: *GameState) void {
        assert(self.one == self.round and self.two == self.round, "one and two should be on the same round as round property");
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
};
