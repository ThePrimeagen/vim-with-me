const std = @import("std");
const gamestate = @import("game-state.zig");
const canvas = @import("canvas.zig");

const Allocator = std.mem.Allocator;
const GameState = gamestate.GameState;

pub const Renderer = struct {
    canvas: canvas.Canvas,
    output: []u8,
    count: u32,

    pub fn init(rows: usize, cols: usize, alloc: Allocator) !Renderer {
        return .{
            .canvas = try canvas.Canvas.init(rows, cols, alloc),
            .output = undefined,
            .count = 0,
        };
    }

    pub fn deinit(self: *Renderer) void {
        self.canvas.deinit();
    }

    pub fn render(self: *Renderer, gs: *GameState) void {
        for (gs.towers.items) |*t| {
            self.canvas.place(t.rSized, &t.rCells);
        }

        var buff: [10]u8 = undefined;
        try std.fmt.format("renders: {}", &buff, self.count);

        self.canvas.writeText(.{
            .rows = 0,
            .cols = 0,
        }, &buff);

        self.canvas.render();
        self.output = self.canvas.renderBuffer;
    }
};

