const assert = @import("assert").assert;
const std = @import("std");
const io = std.io;

const initialClear: [14]u8 = .{
    '', '[', '1', ';', '1', 'H', '', '[', '2', 'J', '', '[', ';', 'H'
};

const clear: [4]u8 = .{
    '', '[', ';', 'H',
};

pub const Output = struct {
    firstPrint: bool,
    rows: usize,
    cols: usize,
    fn init(rows: usize, cols: usize) Output {
        return .{
            .firstPrint = true,
            .rows = rows,
            .cols = cols,
        };
    }

    fn frame(self: *Output, f: []u8) !void {
        assert(f.len == self.rows * self.cols, "you must hand in a frame that matches rows and cols");
        if (self.firstPrint) {
        }
    }
};

pub fn initScreen() !void {
    io.getStdOut().writer().print("{s}", initialClear);
}
