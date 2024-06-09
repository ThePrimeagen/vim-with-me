const assert = @import("assert").assert;
const std = @import("std");
const io = std.io;

const initialClear: [14]u8 = .{
    '', '[', '1', ';', '1', 'H', '', '[', '2', 'J', '', '[', ';', 'H'
};

const clear: [4]u8 = .{
    '', '[', ';', 'H',
};

const newline: [2]u8 = .{
    '\n', '\r',
};

const ColorOutput = struct {
    text: []u8,
    r: u8,
    g: u8,
    b: u8,
};

pub const Output = struct {
    firstPrint: bool,
    rows: usize,
    cols: usize,

    pub fn frame(self: *Output, f: []ColorOutput) !void {
        assert(f.len == self.rows * self.cols, "you must hand in a frame that matches rows and cols");
        const stdout = io.getStdOut().writer();
        var buffered = io.bufferedWriter(stdout);
        var writer = buffered.writer();

        if (self.firstPrint) {
            try writer.write(initialClear);
        }

        var charCount = 0;
        for (f) |*c| {
            const text = c.text;

            for (text) |t| {
                try writer.writeByte(t);
                charCount += 1;

                if (charCount % self.cols == 0) {
                    try writer.write(newline);
                }
            }
        }

        assert(charCount == self.rows * self.cols, "did not produce the correct amount of characters");
        try buffered.flush();
    }
};

pub fn init(rows: usize, cols: usize) Output {
    return .{
        .firstPrint = true,
        .rows = rows,
        .cols = cols,
    };
}
