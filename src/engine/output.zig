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
    twoCreepDamage: usize,
    oneTowers: []TowerValue,
    twoTowers: []TowerValue,
    oneTowerPlacementRange: math.Range,
    twoTowerPlacementRange: math.Range,
    oneCreepSpawnRange: math.Range,
    twoCreepSpawnRange: math.Range,
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
        .twoCreepDamage = gs.twoCreepDamage,
        .oneTowers = oneTowers.items,
        .twoTowers = twoTowers.items,
        .oneTowerPlacementRange = gs.oneNoBuildTowerRange,
        .twoTowerPlacementRange = gs.twoNoBuildTowerRange,
        .oneCreepSpawnRange = gs.oneCreepRange,
        .twoCreepSpawnRange = gs.twoCreepRange,
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
