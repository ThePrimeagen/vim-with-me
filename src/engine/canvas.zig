const a = @import("../assert/assert.zig");
const std = @import("std");

const assert = a.assert;
const framer = @import("framer.zig");
const math = @import("../math/math.zig");
const objects = @import("../objects/objects.zig");
const colors = objects.colors;
const Values = objects.Values;

const AnsiFramer = framer.AnsiFramer;
const Allocator = std.mem.Allocator;
const Vec2 = math.Vec2;
const Position = math.Position;
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
    alloc: Allocator,
    values: *const Values,

    pub fn init(alloc: Allocator, values: *const Values) !Canvas {
        var canvas: Canvas = .{
            .buffer = try alloc.alloc(u8, values.size * 32), // no idea how big ansi encoding is
            .renderBuffer = undefined,
            .bufferLen = 0,
            .cells = try alloc.alloc(Cell, values.size),
            .framer = AnsiFramer.init(values),
            .alloc = alloc,
            .values = values,
        };

        canvas.reset();
        return canvas;
    }

    pub fn deinit(self: *Canvas) void {
        self.alloc.free(self.buffer);
        self.alloc.free(self.cells);
    }

    pub fn writeText(self: *Canvas, pos: math.Position, text: []const u8, color: colors.Color) void {
        assert(pos.row < self.values.rows, "cannot write text off the screen rows");
        assert(pos.col + text.len < self.values.cols, "cannot write text off screen cols");

        const offset = pos.row * self.values.cols + pos.col;
        for (text, offset..) |txt, idx| {
            self.cells[idx].text = txt;
            self.cells[idx].color = color;
        }
    }

    pub fn background(self: *Canvas, pos: Position, col: Color) void {
        assert(pos.row < self.values.rows, "background color set off the grid");
        assert(pos.col < self.values.cols, "background color set off the grid");

        const offset = pos.row * self.values.cols + pos.col;
        self.cells[offset].background = col;
    }

    // TODO: I think that position could be swapped here
    pub fn place(self: *Canvas, sized: Sized, cells: []const Cell) void {
        assert(sized.cols != 0, "cannot render a 0 sized object");
        assert(cells.len > 0, "writing an empty object");
        assert(cells.len % sized.cols == 0, "must provide a square");

        // TODO: rethink these?  Just have the canvas draw what it can?
        assert(sized.pos.row + (cells.len / sized.cols - 1) < self.values.rows, "cannot write text off the screen rows");
        assert(sized.pos.col + sized.cols < self.values.cols, "cannot paint object off screen cols");
        for (cells, 0..) |cell, idx| {
            const col = idx % sized.cols;
            const row = idx / sized.cols;
            const offset = (sized.pos.row + row) * self.values.cols + sized.pos.col + col;
            self.cells[offset].text = cell.text;
            self.cells[offset].color = cell.color;
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
    var values = Values{.rows = 10, .cols = 10};
    values.init();

    var canvas = try Canvas.init(t.allocator, &values);
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
