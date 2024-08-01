const std = @import("std");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const objects = @import("../objects/objects.zig");

const GS = objects.gamestate.GameState;

fn printTowers(gs: *GS, team: u8) !void {
    const err = std.io.getStdErr();
    var printed = false;
    for (gs.towers.items) |t| {
        if (t.team == team and t.alive) {
            if (printed) {
                try err.writeAll(", ");
            }
            printed = true;

            try err.writeAll(try std.fmt.bufPrint(scratchBuf(100), "({}, {}, {}, {})", .{
                t.pos.position().row,
                t.pos.position().col,
                t.ammo,
                t.level,
            }));
        }
    }
}

// NOTE: purely hard coded... don't think i need to change that...
// later you will hate this :)
pub fn printPrompt(gs: *GS) !void {
    const err = std.io.getStdErr();

    try err.writeAll("prompt:\n");
    try err.writeAll(try std.fmt.bufPrint(scratchBuf(50), "rows: {}\n", .{gs.values.rows}));
    try err.writeAll(try std.fmt.bufPrint(scratchBuf(50), "cols: {}\n", .{gs.values.cols}));
    try err.writeAll(try std.fmt.bufPrint(scratchBuf(50), "allowed towers: {}\n", .{gs.twoAvailableTower}));
    try err.writeAll(try std.fmt.bufPrint(scratchBuf(50), "your creep damage: {}\n", .{gs.twoCreepDamage}));
    try err.writeAll(try std.fmt.bufPrint(scratchBuf(50), "enemy creep damage: {}\n", .{gs.oneCreepDamage}));
    try err.writeAll("your towers: ");
    try printTowers(gs, objects.Values.TEAM_TWO);
    try err.writeAll("\n");
    try err.writeAll("enemy towers: ");
    try printTowers(gs, objects.Values.TEAM_ONE);
    try err.writeAll("\n");


    const trPos = gs.twoNoBuildTowerRange;
    try err.writeAll(try std.fmt.bufPrint(scratchBuf(100), "tower placement range: TL={},{} BR={},{}\n", .{
        trPos.startRow, 0,
        trPos.endRow, gs.values.cols,
    }));

    try err.writeAll(try std.fmt.bufPrint(scratchBuf(100), "creep spawn range: SR={} ER={}\n", .{
        gs.twoCreepRange.startRow,
        gs.twoCreepRange.endRow,
    }));

    try err.writeAll(try std.fmt.bufPrint(scratchBuf(50), "round: {}", .{gs.round}));
    try err.writeAll("\n");
    try err.writeAll("\n");

}

