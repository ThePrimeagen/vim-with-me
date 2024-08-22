const a = @import("../assert/assert.zig");
const std = @import("std");

const assert = a.assert;
const Dump = a.Dump;

const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const math = @import("../math/math.zig");
const projectile = @import("projectile.zig");
const tower = @import("tower.zig");
const creep = @import("creep.zig");
const colors = @import("colors.zig");
const messages = @import("messages.zig");
const Values = @import("values.zig");
pub const Target = @import("target.zig").Target;

// TODO: Make this adjustable
const Vec2 = math.Vec2;
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
    playingStartUS: i64 = 0,
    round: usize = 0,
    countdown: i64 = 0,

    noBuildZone: bool = true,
    noBuildRange: Range = Range{},

    boardChanged: usize = 0,
    activeCreepCount: isize = 0,

    oneAvailableTower: isize = 0,
    oneTowerCount: usize = 0,
    onePositions: math.PossiblePositions,
    oneStats: Stats = Stats{},
    oneCreepRange: Range = Range{},
    oneNoBuildTowerRange: Range = Range{},
    oneCreepDamage: usize = 1,
    oneTotalTowersBuild: usize = 0,
    oneTotalProjectiles: usize = 0,
    oneTotalTowerUpgrades: usize = 0,
    oneTotalCreepDamage: usize = 0,
    oneTotalTowerDamage: usize = 0,
    oneTotalDamageFromCreeps: usize = 0,

    twoAvailableTower: isize = 0,
    twoTowerCount: usize = 0,
    twoPositions: math.PossiblePositions,
    twoStats: Stats = Stats{},
    twoCreepRange: Range = Range{},
    twoNoBuildTowerRange: Range = Range{},
    twoCreepDamage: usize = 1,
    twoTotalTowersBuild: usize = 0,
    twoTotalProjectiles: usize = 0,
    twoTotalTowerUpgrades: usize = 0,
    twoTotalCreepDamage: usize = 0,
    twoTotalTowerDamage: usize = 0,
    twoTotalDamageFromCreeps: usize = 0,

    time: i64 = 0,
    loopDeltaUS: i64 = 0,
    updates: i64 = 0,

    towers: TowerList,
    creeps: CreepList,
    projectile: ProjectileList,
    board: []bool,
    alloc: Allocator,

    fns: ?*const GameStateFunctions = null,

    values: *Values,

    pub fn playingString(gs: *GameState) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(150), "GameState({}): round={} one={} two={}\ntime={} playingTime={} overTime={}", .{
            gs.playing, gs.round, gs.one, gs.two,
            gs.time, gs.playingStartUS, gs.playingStartUS + gs.values.roundTimeUS,
        });
    }

    pub fn init(alloc: Allocator, values: *Values) !GameState {
        var gs = .{
            .towers = TowerList.init(alloc),
            .creeps = CreepList.init(alloc),
            .projectile = ProjectileList.init(alloc),
            .board = try alloc.alloc(bool, values.size),
            .alloc = alloc,
            .values = values,
            .onePositions = .{.team = Values.TEAM_ONE, .len = 0, .positions = undefined },
            .twoPositions = .{.team = Values.TEAM_TWO, .len = 0, .positions = undefined },
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
        std.debug.print("round = {}, oneAvailable = {}, twoAvailable = {}\n", .{self.round, self.oneAvailableTower, self.twoAvailableTower});
        std.debug.print("time = {}, loopDeltaUS = {}\n", .{self.time, self.loopDeltaUS});

        std.debug.print("one coords:\n", .{});
        for (0..self.onePositions.len) |idx| {
            const pos = self.onePositions.positions[idx];
            if (pos) |p| {
                std.debug.print("  {s}\n", .{a.u(p.string())});
            }
        }
        std.debug.print("\n", .{});

        std.debug.print("two coords:\n", .{});
        for (0..self.twoPositions.len) |idx| {
            const pos = self.twoPositions.positions[idx];
            if (pos) |p| {
                std.debug.print("  {s}\n", .{a.u(p.string())});
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

        for (0..self.values.cols) |c| {
                std.debug.print("{:0>3}", .{c});
        }
        std.debug.print("\n", .{});

        outer:
        for (self.board, 0..) |b, idx| {
            if (idx > 0 and idx % self.values.cols == 0) {
                std.debug.print("\n", .{});
            }

            for (self.creeps.items) |*c| {
                if (idx == c.pos.position().toIdx(self.values.cols)) {
                    std.debug.print("c  ", .{});
                    continue :outer;
                }
            }

            const v: usize = if (b) 1 else 0;
            std.debug.print("{}  ", .{v});
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
