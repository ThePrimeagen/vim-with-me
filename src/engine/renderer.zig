const std = @import("std");
const objects = @import("objects");
const gamestate = objects.gamestate;
const Values = objects.Values;
const canvas = @import("canvas.zig");

const towers = @import("tower.zig");
const creeps = @import("creep.zig");

const Allocator = std.mem.Allocator;
const GameState = gamestate.GameState;

pub const Renderer = struct {
    canvas: canvas.Canvas,
    output: []u8,
    count: u32,
    values: *const objects.Values,

    pub fn init(alloc: Allocator, values: *const Values) !Renderer {
        return .{
            .values = values,
            .canvas = try canvas.Canvas.init(alloc, values),
            .output = undefined,
            .count = 0,
        };
    }

    pub fn deinit(self: *Renderer) void {
        self.canvas.deinit();
    }

    pub fn render(self: *Renderer, gs: *GameState) !void {
        for (gs.towers.items) |*t| {
            towers.render(t, gs);
            self.canvas.place(t.rSized, &t.rCells);
        }

        for (gs.creeps.items) |*c| {
            creeps.render(c, gs);
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
