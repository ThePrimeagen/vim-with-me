const assert = @import("assert").assert;
const std = @import("std");
const io = std.io;

pub const Color = struct {
    r: u8,
    g: u8,
    b: u8,
};

pub const Cell = struct {
    text: u8,
    color: Color,
};

const initialClear: [14]u8 = .{
    '', '[', '1', ';', '1', 'H', '', '[', '2', 'J', '', '[', ';', 'H'
};

const newFrame: [4]u8 = .{
    '', '[', ';', 'H',
};

const foregroundColor: [4]u8 = .{
    '', '[', '3', '8', ';', '2', ';',
};


const newline: [2]u8 = .{
    '\n', '\r',
};

fn writeAnsiColor(color: Color, writer: *io.BufferedWriter(4096, []u8).Writer) !void {
    try writer.write(foregroundColor);
    try writer.print("{};{};{}m", .{color.r, color.g, color.b});
}

pub const AnsiFramer = struct {
    firstPrint: bool,
    rows: usize,
    cols: usize,

    pub fn frame(self: *AnsiFramer, f: []Cell, out: []u8) !void {
        assert(f.len == self.rows * self.cols, "you must hand in a frame that matches rows and cols");

        const buffered = io.bufferedWriter(out);
        var writer = buffered.writer();

        if (self.firstPrint) {
            try writer.write(initialClear);
            self.firstPrint = false;
        }

        var charCount = 0;
        var previous: ?*Cell = null;

        for (f) |*c| {
            const text = c.text;

            if (previous == null) {
                previous = &c.color;
                try writeAnsiColor(c.color, &writer);
            }

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
