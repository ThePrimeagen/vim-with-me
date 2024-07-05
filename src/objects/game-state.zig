const a = @import("../assert/assert.zig");
const std = @import("std");

const assert = a.assert;
const Dump = a.Dump;

const math = @import("../math/math.zig");
const projectile = @import("projectile.zig");
const tower = @import("tower.zig");
const creep = @import("creep.zig");
const colors = @import("colors.zig");
const messages = @import("messages.zig");
const Values = @import("values.zig");
const Target = @import("target.zig").Target;

// TODO: Make this adjustable
const Vec2 = math.Vec2;
const Coord = math.Coord;
const Range = math.Range;

const Message = messages.Message;

const ArrayList = std.ArrayList;
const Allocator = std.mem.Allocator;

const Tower = tower.Tower;
const Creep = creep.Creep;
const Projectile = projectile.Projectile;

const TowerList = ArrayList(Tower);
const CreepList = ArrayList(Creep);
const ProjectileList = ArrayList(Projectile);

pub const Stats = struct {
    creepsKilled: usize = 0,
    creepsEscaped: usize = 0,
    towersLost: usize = 0,
    towersKilled: usize = 0,
    shots: usize = 0,
};

pub const GameState = struct {
    playing: bool = true,
    round: usize = 1,

    noBuildZone: bool = true,
    noBuildRange: Range = Range{},

    boardChanged: usize = 0,

    one: usize = 0,
    oneCoords: [3]?Coord,
    oneStats: Stats = Stats{},
    oneRange: Range = Range{},

    two: usize = 0,
    twoCoords: [3]?Coord,
    twoStats: Stats = Stats{},
    twoRange: Range = Range{},

    time: i64 = 0,
    loopDeltaUS: i64 = 0,
    updates: i64 = 0,

    towers: TowerList,
    creeps: CreepList,
    projectile: ProjectileList,
    board: []bool,
    alloc: Allocator,

    fns: ?*const GameStateFunctions = null,

    values: *const Values,

    pub fn string(gs: *GameState) ![]u8 {
        _ = gs;
        unreachable;
    }

    pub fn init(alloc: Allocator, values: *const Values) !GameState {
        var gs = .{
            .towers = TowerList.init(alloc),
            .creeps = CreepList.init(alloc),
            .projectile = ProjectileList.init(alloc),
            .board = try alloc.alloc(bool, values.size),
            .alloc = alloc,
            .values = values,
            .oneCoords = .{null, null, null},
            .twoCoords = .{null, null, null},
        };

        for (0..values.size) |i| {
            gs.board[i] = true;
        }

        return gs;
    }

    pub fn deinit(gs: *GameState) void {
        gs.towers.deinit();
        for (gs.creeps.items) |*c| {
            c.deinit();
        }

        gs.creeps.deinit();
        gs.projectile.deinit();
        gs.alloc.free(gs.board);
    }

    pub fn dumper(self: *GameState) Dump {
        return Dump.init(self);
    }

    pub fn dump(self: *GameState) void {
        std.debug.print("------ GameState ------\n", .{});
        std.debug.print("values: {s}\n", .{a.u(self.values.string())});
        std.debug.print("round = {}, one = {}, two = {}\n", .{self.round, self.one, self.two});
        std.debug.print("time = {}, loopDeltaUS = {}\n", .{self.time, self.loopDeltaUS});

        std.debug.print("one coords:\n", .{});
        for (&self.oneCoords) |c| {
            if (c) |coord| {
                std.debug.print("  {s}\n", .{a.u(coord.string())});
            }
        }
        std.debug.print("\n", .{});

        std.debug.print("two coords:\n", .{});
        for (&self.twoCoords) |c| {
            if (c) |coord| {
                std.debug.print("  {s}\n", .{a.u(coord.string())});
            }
        }

        self.debugBoard();

        std.debug.print("Towers:\n", .{});
        for (self.towers.items) |*t| {
            std.debug.print("  {s}\n", .{a.u(t.pos.position().string())});
        }
        std.debug.print("\nCreeps:\n", .{});
        for (self.creeps.items) |*c| {
            std.debug.print("  {s}\n", .{a.u(c.string())});
        }
        std.debug.print("\nProjectiles\n", .{});
        for (self.projectile.items) |*p| {
            std.debug.print("  {s}\n", .{a.u(p.pos.position().string())});
        }
        std.debug.print("\n", .{});

    }

    pub fn debugBoard(self: *GameState) void {
        std.debug.print("\nBoard:\n", .{});
        outer:
        for (self.board, 0..) |b, idx| {
            if (idx > 0 and idx % self.values.cols == 0) {
                std.debug.print("\n", .{});
            }

            for (self.creeps.items) |*c| {
                if (idx == c.pos.position().toIdx(self.values.cols)) {
                    std.debug.print("c ", .{});
                    continue :outer;
                }
            }

            const v: usize = if (b) 1 else 0;
            std.debug.print("{} ", .{v});
        }
        std.debug.print("\n", .{});
    }
};

pub const GameStateFunctions = struct {
    placeProjectile: *const fn(self: *GameState, tower: *Tower, target: Target) Allocator.Error!usize,
    towerDied: *const fn(self: *GameState, tower: *Tower) void,
    creepKilled: *const fn(self: *GameState, creep: *Creep) void,
    shot: *const fn(self: *GameState, tower: *Tower) void,
    strike: *const fn(self: *GameState, p: *Projectile) void,
};
