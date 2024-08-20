const std = @import("std");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const objects = @import("../objects/objects.zig");
const math = @import("../math/math.zig");
const gamestate = @import("game-state.zig");

const GS = objects.gamestate.GameState;
const Allocator = std.mem.Allocator;
const Values = objects.Values;

const TowerValue = struct {
    row: usize,
    col: usize,
    level: usize,
    ammo: usize,
};

const TowerList = std.ArrayList(TowerValue);
const GameStateJSON = struct {
    rows: usize,
    cols: usize,
    allowedTowers: isize,

    oneCreepDamage: usize,
    oneTowers: []TowerValue,
    oneTowerPlacementRange: math.Range,
    oneCreepSpawnRange: math.Range,
    oneTotalTowersBuild: usize = 0,
    oneTotalProjectiles: usize = 0,
    oneTotalTowerUpgrades: usize = 0,
    oneTotalCreepDamage: usize = 0,
    oneTotalTowerDamage: usize = 0,
    oneTotalDamageFromCreeps: usize = 0,

    twoCreepDamage: usize,
    twoTowers: []TowerValue,
    twoTowerPlacementRange: math.Range,
    twoCreepSpawnRange: math.Range,
    twoTotalTowersBuild: usize = 0,
    twoTotalProjectiles: usize = 0,
    twoTotalTowerUpgrades: usize = 0,
    twoTotalCreepDamage: usize = 0,
    twoTotalTowerDamage: usize = 0,
    twoTotalDamageFromCreeps: usize = 0,

    round: usize,
    playing: bool,
    finished: bool,
    winner: u8,
};

fn fromGameState(gs: *GS, oneTowers: *TowerList, twoTowers: *TowerList) !GameStateJSON {
    for (gs.towers.items) |t| {
        if (!t.alive) {
            continue;
        }

        const tower = .{
            .row = t.pos.position().row,
            .col = t.pos.position().col,
            .ammo = t.ammo,
            .level = t.level,
        };
        switch (t.team) {
            Values.TEAM_ONE => try oneTowers.append(tower),
            Values.TEAM_TWO => try twoTowers.append(tower),
            else => unreachable,
        }
    }

    return .{
        .rows = gs.values.rows,
        .cols = gs.values.cols,
        .allowedTowers = gs.oneAvailableTower,

        .oneCreepDamage = gs.oneCreepDamage,
        .oneTowers = oneTowers.items,
        .oneTowerPlacementRange = gs.oneNoBuildTowerRange,
        .oneCreepSpawnRange = gs.oneCreepRange,
        .oneTotalTowersBuild = gs.oneTotalTowersBuild,
        .oneTotalProjectiles = gs.oneTotalProjectiles,
        .oneTotalTowerUpgrades = gs.oneTotalTowerUpgrades,
        .oneTotalCreepDamage = gs.oneTotalCreepDamage,
        .oneTotalTowerDamage = gs.oneTotalTowerDamage,
        .oneTotalDamageFromCreeps = gs.oneTotalDamageFromCreeps,

        .twoTowerPlacementRange = gs.twoNoBuildTowerRange,
        .twoCreepSpawnRange = gs.twoCreepRange,
        .twoTowers = twoTowers.items,
        .twoCreepDamage = gs.twoCreepDamage,
        .twoTotalTowersBuild = gs.twoTotalTowersBuild,
        .twoTotalProjectiles = gs.twoTotalProjectiles,
        .twoTotalTowerUpgrades = gs.twoTotalTowerUpgrades,
        .twoTotalCreepDamage = gs.twoTotalCreepDamage,
        .twoTotalTowerDamage = gs.twoTotalTowerDamage,
        .twoTotalDamageFromCreeps = gs.twoTotalDamageFromCreeps,

        .round = gs.round,
        .playing = gs.playing,
        .finished = gamestate.completed(gs),
        .winner = gamestate.getWinner(gs),
    };
}

// NOTE: purely hard coded... don't think i need to change that...
// later you will hate this :)
pub fn writeState(alloc: Allocator, gs: *GS, out: std.fs.File) !void {
    var one = TowerList.init(alloc);
    var two = TowerList.init(alloc);
    defer one.deinit();
    defer two.deinit();

    var v = try fromGameState(gs, &one, &two);
    var buf: [8196]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&buf);
    var string = std.ArrayList(u8).init(fba.allocator());
    try std.json.stringify(&v, .{}, string.writer());

    try out.writeAll(string.items);
    try out.writeAll("\n");
}

pub fn writeValues(gs: *GS) !void {
    var buf: [8196]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&buf);
    var string = std.ArrayList(u8).init(fba.allocator());
    try std.json.stringify(gs.values, .{}, string.writer());
}
