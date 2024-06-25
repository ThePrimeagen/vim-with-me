const a = @import("assert");
const assert = a.assert;
const framer = @import("framer.zig");
const math = @import("math");
const colors = @import("objects").colors;
const std = @import("std");

const AnsiFramer = framer.AnsiFramer;
const Allocator = std.mem.Allocator;
const Position = math.Vec2;
const Sized = math.Sized;
const Color = colors.Color;
const Cell = colors.Cell;

const EMPTY_CELL = Cell{
    .text = ' ',
    .color = .{
        .r = 0,
        .g = 0,
        .b = 0,
    },
};

pub const Canvas = struct {
    framer: AnsiFramer,
    renderBuffer: []u8,
    buffer: []u8,
    bufferLen: usize,
    cells: []Cell,
    rows: usize,
    cols: usize,
    alloc: Allocator,

    pub fn init(rows: usize, cols: usize, alloc: Allocator) !Canvas {
        const size = rows * cols;
        var canvas: Canvas = .{
            .buffer = try alloc.alloc(u8, size * 16), // no idea how big ansi encoding is
            .renderBuffer = undefined,
            .bufferLen = 0,
            .cells = try alloc.alloc(Cell, size),
            .rows = rows,
            .cols = cols,
            .framer = AnsiFramer.init(rows, cols),
            .alloc = alloc,
        };

        canvas.reset();
        return canvas;
    }

    pub fn deinit(self: *Canvas) void {
        self.alloc.free(self.buffer);
        self.alloc.free(self.cells);
    }

    pub fn writeText(self: *Canvas, pos: math.Position, text: []const u8, color: colors.Color) void {
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

    // TODO: I think that position could be swapped here
    pub fn place(self: *Canvas, sized: Sized, cells: []const Cell) void {

        std.debug.print("{s}\n", .{a.u(sized.string())});

        assert(sized.cols != 0, "cannot render a 0 sized object");
        assert(cells.len > 0, "writing an empty object");
        assert(cells.len % sized.cols == 0, "must provide a square");

        // TODO: rethink these?  Just have the canvas draw what it can?
        std.debug.print("here we are: {} + {} / {} < {}\n", .{sized.pos.row, cells.len, sized.cols, self.rows});
        assert(sized.pos.row + cells.len / sized.cols < self.rows, "cannot write text off the screen rows");
        assert(sized.pos.col + sized.cols < self.cols, "cannot paint object off screen cols");

        for (cells, 0..) |cell, idx| {
            const col = idx % sized.cols;
            const row = idx / sized.cols;
            const offset = (sized.pos.row + row) * self.cols + sized.pos.col + col;
            self.cells[offset] = cell;
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
    const x: []const Cell = &(.{
        .{.text = 'x', .color = .{.r = 69, .g = 69, .b = 69}}
    } ** 3);
    const y: []const Cell = &(.{
        .{.text = 'y', .color = .{.r = 69, .g = 69, .b = 69}}
    } ** 3);
    const z: []const Cell = &(.{
        .{.text = 'z', .color = .{.r = 69, .g = 69, .b = 69}}
    } ** 3);

    const image: []const Cell = x ++ y ++ z;

    canvas.place(.{
        .cols = 3,
        .pos = .{.row = 3, .col = 3}
    }, image);

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


    canvas.place(.{.cols = 3, .pos = .{.row = 4, .col = 4}}, x);
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
