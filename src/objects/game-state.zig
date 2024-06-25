const assert = @import("assert").assert;
const std = @import("std");

const math = @import("math");
const projectile = @import("projectile.zig");
const tower = @import("tower.zig");
const creep = @import("creep.zig");
const colors = @import("colors.zig");
const messages = @import("messages.zig");
const Values = @import("values.zig");

// TODO: Make this adjustable
const Position = math.Vec2;
const Coord = math.Coord;

const Message = messages.Message;

const ArrayList = std.ArrayList;
const Allocator = std.mem.Allocator;

const Tower = tower.Tower;
const Creep = creep.Creep;
const Projectile = projectile.Projectile;

const TowerList = ArrayList(Tower);
const CreepList = ArrayList(Creep);
const ProjectileList = ArrayList(Projectile);

pub const GameState = struct {
    playing: bool = true,
    round: i32 = 1,
    one: i32 = 0,
    two: i32 = 0,
    rows: usize = 0,
    cols: usize = 0,

    time: i64 = 0,
    loopDeltaUS: i64 = 0,
    updates: i64 = 0,

    towers: TowerList,
    creeps: CreepList,
    projectile: ProjectileList,
    board: []bool,
    alloc: Allocator,

    values: *const Values,

    pub fn string(gs: *GameState) ![]u8 {
        _ = gs;
        unreachable;
    }

    pub fn init(alloc: Allocator, values: *const Values) !GameState {
        return .{
            .towers = TowerList.init(alloc),
            .creeps = CreepList.init(alloc),
            .projectile = ProjectileList.init(alloc),
            .board = try alloc.alloc(bool, values.rows * values.cols),
            .alloc = alloc,
            .values = values,
        };
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
};

pub const Target = union(enum) {
    creep: usize,
    tower: usize,
};
