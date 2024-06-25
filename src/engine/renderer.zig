const std = @import("std");
const gamestate = @import("objects").gamestate;
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

    pub fn render(self: *Renderer, gs: *GameState) !void {
        for (gs.towers.items) |*t| {
            t.render();
            self.canvas.place(t.rSized, &t.rCells);
        }

        for (gs.creeps.items) |*c| {
            c.render();
            self.canvas.place(c.rSized, &c.rCells);
        }

        var buff: [15]u8 = undefined;
        _ = try std.fmt.bufPrint(&buff, "renders: {}", .{self.count});

        self.canvas.writeText(.{
            .row = 0,
            .col = 14,
        }, &buff, .{.r = 255, .g = 255, .b = 255});

        try self.canvas.render();
        self.output = self.canvas.renderBuffer;
        self.count += 1;
    }
};
