const assert = @import("../assert/assert.zig").assert;
const std = @import("std");
const objects = @import("../objects/objects.zig");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;

const gamestate = objects.gamestate;
const Values = objects.Values;
const canvas = @import("canvas.zig");

const towers = @import("tower.zig");
const creeps = @import("creep.zig");
const projectiles = @import("projectile.zig");

const Allocator = std.mem.Allocator;
const GameState = gamestate.GameState;
const TEXT_AREA_COLS = 30;

pub const Renderer = struct {
    canvas: canvas.Canvas,
    output: []u8,
    count: u32,
    values: *const objects.Values,
    textOffset: usize,
    rendererValues: *objects.Values,
    alloc: Allocator,

    pub fn init(alloc: Allocator, values: *const Values) !Renderer {
        assert(values.cols > 30, "not enough columns for game information.  requires 30 cols for text");

        var rendererValues: *Values = try alloc.create(Values);
        Values.copyInto(values, rendererValues);
        rendererValues.cols = values.cols + TEXT_AREA_COLS;
        Values.init(rendererValues);

        return .{
            .textOffset = values.cols,
            .values = values,
            .canvas = try canvas.Canvas.init(alloc, rendererValues),
            .output = undefined,
            .count = 0,
            .rendererValues = rendererValues,
            .alloc = alloc,
        };
    }

    pub fn deinit(self: *Renderer) void {
        self.canvas.deinit();
        self.alloc.destroy(self.rendererValues);
    }

    pub fn render(self: *Renderer, gs: *GameState) !void {
        for (gs.towers.items) |*t| {
            towers.render(t, gs);
            self.canvas.place(t.rSized, &t.rCells);
        }

        for (gs.creeps.items) |*c| {
            if (creeps.completed(c) or !c.alive) {
                continue;
            }

            creeps.render(c, gs);
            self.canvas.place(c.rSized, &c.rCells);
        }

        for (gs.projectile.items) |*p| {
            if (!p.alive) {
                continue;
            }

            projectiles.render(p, gs);
            self.canvas.place(p.rSized, &p.rCells);
        }

        try self.text(gs);
        try self.canvas.render();
        self.output = self.canvas.renderBuffer;
        self.count += 1;
    }

    pub fn text(self: *Renderer, gs: *GameState) !void {
        const buf = scratchBuf(50);

        // round
        const roundBuf = try std.fmt.bufPrint(buf, "round: {}", .{gs.round});

        var row: usize = 0;

        self.canvas.writeText(.{
            .row = row,
            .col = self.textOffset,
        }, roundBuf, .{.r = 255, .g = 255, .b = 255});

        row += 1;
    }
};
