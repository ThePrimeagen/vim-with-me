const math = @import("../math/math.zig");
const assert = @import("../assert/assert.zig").assert;
const colors = @import("colors.zig");

const Range = math.Range;

pub const TEXT = 'o';
pub const COLOR = colors.Red;

// I LOVE DOING THIS
var cells: [8192]colors.Cell = undefined;

pub fn createCells(range: Range, cols: usize) []colors.Cell {
    const count = (range.endRow - range.startRow) * cols;
    assert(cells.len > count, "your assumptions have been broken, no build zone is much larger");

    for (0..count) |idx| {
        cells[idx] = .{
            .text = TEXT,
            .color = COLOR,
        };
    }

    return cells[0..count];
}
