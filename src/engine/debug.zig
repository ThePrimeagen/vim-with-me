const assert = @import("../assert/assert.zig").assert;
const colors = @import("../objects/objects.zig").colors;

const Cell = colors.Cell;
var cells: [1024]Cell = undefined;
var cellsIdx: usize = 0;

pub fn pathAStarToCells(cols: usize, current: usize, board: []const bool, parents: []const isize, seen: []const bool, pathMaybe: ?[]usize) []Cell {
    assert(board.len % cols == 0, "board is not square");
    assert(board.len == parents.len, "parents and board have to be the same length");
    assert(board.len == seen.len, "seen and board have to be the same length");

    cellsIdx = 0;

    for (0..board.len) |i| {
        const blocked = board[i];
        const parent = parents[i];
        const s = seen[i];

        cells[i] = .{
            .text = ' ',
            .color = colors.White,
        };

        if (i == current and pathMaybe == null) {
            cells[i].background = colors.Red;
        } else if (i == current and pathMaybe != null) {
            cells[i].background = colors.Green;
        } else if (blocked == false) {
            cells[i].background = colors.DarkGrey;
        } else if (s) {
            cells[i].background = colors.Grey;
        } else if (parent >= 0) {
            cells[i].background =  colors.Blue;
        } else {
            cells[i] = .{
                .text = ' ',
                .color = colors.White,
                .background =  colors.LightGrey,
            };
        }

        cellsIdx += 1;
    }

    if (pathMaybe) |path| for (path) |p| {
        cells[p].background = colors.Green;
    };

    return cells[0..cellsIdx];
}

