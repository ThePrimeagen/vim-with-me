const assert = @import("assert").assert;
const std = @import("std");
const objects = @import("objects");

// TODO: Make this adjustable
const Position = objects.Vec2;
const Message = objects.Message;
const Coord = objects.Coord;

const ArrayList = std.ArrayList;
const Allocator = std.mem.Allocator;

const Tower = objects.Tower;
const Creep = objects.Creep;
const Projectile = objects.Projectile;

const TowerList = ArrayList(objects.Tower);
const CreepList = ArrayList(objects.Creep);
const ProjectileList = ArrayList(objects.Projectile);

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

    pub fn deinit(gs: *GameState) void {
        gs.allocator.free(gs.towers);
        gs.allocator.free(gs.creeps);
        gs.allocator.free(gs.projectile);
    }
};

