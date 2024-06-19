const assert = @import("assert").assert;
const output = @import("output/output.zig");
const types = @import("types.zig");
const std = @import("std");

const AnsiFramer = output.AnsiFramer;
const Allocator = std.mem.Allocator;
const Position = types.Position;

const EMPTY_CELL = output.Cell{
    .text = ' ',
    .color = .{
        .r = 0,
        .g = 0,
        .b = 0,
    },
};

pub const Canvas = struct {
    framer: output.AnsiFramer,
    renderBuffer: []u8,
    buffer: []u8,
    bufferLen: usize,
    cells: []output.Cell,
    rows: usize,
    cols: usize,
    alloc: Allocator,

    pub fn init(rows: usize, cols: usize, alloc: Allocator) !Canvas {
        const size = rows * cols;
        var canvas: Canvas = .{
            .buffer = try alloc.alloc(u8, size * 16), // no idea how big ansi encoding is
            .renderBuffer = undefined,
            .bufferLen = 0,
            .cells = try alloc.alloc(output.Cell, size),
            .rows = rows,
            .cols = cols,
            .framer = output.AnsiFramer.init(rows, cols),
            .alloc = alloc,
        };

        canvas.reset();
        return canvas;
    }

    pub fn deinit(self: *Canvas) void {
        self.alloc.free(self.buffer);
        self.alloc.free(self.cells);
    }

    pub fn writeText(self: *Canvas, pos: Position, text: []const u8, color: output.Color) void {
        assert(pos.row < self.rows, "cannot write text off the screen rows");
        assert(pos.col + text.len < self.cols, "cannot write text off screen cols");

        const offset = pos.row * self.cols + pos.col;
        for (text, offset..) |txt, idx| {
            self.cells[idx] = .{
                .text = txt,
                .color = color,
            };
        }
    }

    pub fn place(self: *Canvas, pos: Position, cells: []const []const output.Cell) void {
        assert(cells.len > 0, "writing an empty object");
        assert(pos.row + cells.len < self.rows, "cannot write text off the screen rows");
        assert(pos.col + cells[0].len < self.cols, "cannot write text off screen cols");

        for (cells, 0..) |cellList, rowOffset| {
            const offset = (pos.row + rowOffset) * self.cols + pos.col;
            for (cellList, offset..) |cell, idx| {
                self.cells[idx] = cell;
            }
        }
    }

    pub fn render(state: *Canvas) !void {
        state.bufferLen = try state.framer.frame(state.cells, state.buffer);
        state.renderBuffer = state.buffer[0..state.bufferLen];
        state.reset();
    }

    fn reset(state: *Canvas) void {
        for (0..state.cells.len) |i| {
            state.cells[i] = EMPTY_CELL;
        }
    }

};


const t = std.testing;
test "i think this test is terribly written, i need help or a doctor" {
    var canvas = try Canvas.init(10, 10, t.allocator);
    defer canvas.deinit();

    const newLine = "\r\n";
    const emptyLine = " " ** 10;
    const x: []const output.Cell = &(.{
        .{.text = 'x', .color = .{.r = 69, .g = 69, .b = 69}}
    } ** 3);
    const y: []const output.Cell = &(.{
        .{.text = 'y', .color = .{.r = 69, .g = 69, .b = 69}}
    } ** 3);
    const z: []const output.Cell = &(.{
        .{.text = 'z', .color = .{.r = 69, .g = 69, .b = 69}}
    } ** 3);

    const image: []const []const output.Cell = &.{
        x,
        y,
        z,
    };

    canvas.place(.{.row = 3, .col = 3}, image);
    try canvas.render();

    var text: [200]u8 = undefined;
    {
        const len = AnsiFramer.parseText(canvas.renderBuffer, &text);
        try t.expectEqualStrings(
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            "   xxx    \r\n".* ++
            "   yyy    \r\n".* ++
            "   zzz    \r\n".* ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine,

            text[0..len]
        );
    }


    canvas.place(.{.row = 4, .col = 4}, &.{x});
    try canvas.render();
    {
        const len = AnsiFramer.parseText(canvas.renderBuffer, &text);

        try t.expectEqualStrings(
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            "    xxx   \r\n".* ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine ++
            emptyLine ++ newLine,

            text[0..len]
        );
    }

}
