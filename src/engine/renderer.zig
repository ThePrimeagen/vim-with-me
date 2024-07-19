const assert = @import("../assert/assert.zig").assert;
const std = @import("std");
const objects = @import("../objects/objects.zig");
const engine = @import("../engine/engine.zig");
const scratch = @import("../scratch/scratch.zig");
const utils = @import("utils.zig");
const math = @import("../math/math.zig");

const gamestate = objects.gamestate;
const colors = objects.colors;
const Values = objects.Values;
const canvas = @import("canvas.zig");

const towers = @import("tower.zig");
const creeps = @import("creep.zig");
const projectiles = @import("projectile.zig");

const Allocator = std.mem.Allocator;
const GameState = gamestate.GameState;
const TEXT_AREA_COLS = 30;
const GRID_AREA_COLS = 2;
const GRID_AREA_ROWS = 1;
const PLAYER_ONE = "Player One";
const PLAYER_TWO = "Player Two";
const scratchBuf = scratch.scratchBuf;
const toNumber = scratch.toNumber;

pub const Renderer = struct {
    canvas: canvas.Canvas,
    output: []u8,
    count: u32,
    values: *const objects.Values,
    gridOffsetX: usize,
    gridOffsetY: usize,
    textOffset: usize,
    rendererValues: *objects.Values,
    alloc: Allocator,

    pub fn init(alloc: Allocator, values: *const Values) !Renderer {
        assert(values.cols > 30, "not enough columns for game information.  requires 30 cols for text");

        var rendererValues: *Values = try alloc.create(Values);
        Values.copyInto(values, rendererValues);
        rendererValues.cols = values.cols + TEXT_AREA_COLS + GRID_AREA_COLS;
        rendererValues.rows = values.rows + GRID_AREA_ROWS;
        Values.init(rendererValues);

        return .{
            .textOffset = values.cols + 2,
            .values = values,
            .canvas = try canvas.Canvas.init(alloc, rendererValues),
            .output = undefined,
            .count = 0,
            .rendererValues = rendererValues,
            .alloc = alloc,
            .gridOffsetX = 2,
            .gridOffsetY = 1,
        };
    }

    pub fn deinit(self: *Renderer) void {
        self.canvas.deinit();
        self.alloc.destroy(self.rendererValues);
    }

    pub fn render(self: *Renderer, gs: *GameState) !void {
        const gridOffset = .{.row = self.gridOffsetY, .col = self.gridOffsetX};

        if (gs.noBuildZone) {
            const cells = objects.nobuild.createCells(gs.noBuildRange, gs.values.cols);
            const sized = gs.noBuildRange.sized(gs.values.cols).add(gridOffset);
            self.canvas.place(sized, cells);
        }

        try self.renderGrid(gs);

        for (gs.towers.items) |*t| {
            if (!t.alive) {
                continue;
            }

            towers.render(t, gs);
            self.canvas.place(t.rSized.add(gridOffset), &t.rCells);
        }

        for (gs.creeps.items) |*c| {
            if (creeps.completed(c) or !c.alive) {
                continue;
            }

            creeps.render(c, gs);
            self.canvas.place(c.rSized.add(gridOffset), &c.rCells);
        }

        for (gs.projectile.items) |*p| {
            if (!p.alive) {
                continue;
            }

            projectiles.render(p, gs);
            self.canvas.place(p.rSized.add(gridOffset), &p.rCells);
        }

        try self.gameStateText(gs);
        try self.canvas.render();
        self.output = self.canvas.renderBuffer;
        self.count += 1;
    }

    fn renderGrid(self: *Renderer, gs: *GameState) !void {

        const columns = @divFloor(gs.values.cols, 5);
        var offset: math.Position = .{.row = 0, .col = 0};
        for (0..columns) |idx| {
            offset.col = self.gridOffsetX + idx * 5;
            self.canvas.writeText(offset, try toNumber(idx * 5), colors.White);
        }

        const rows = @divFloor(gs.values.rows, 3);
        offset.col = 0;
        for (0..rows) |idx| {
            offset.row = self.gridOffsetY + idx * 3;
            self.canvas.writeText(offset, try toNumber(idx * 3), colors.White);
        }


        for (0..self.values.size) |idx| {
        }
    }

    // WAP
    fn gameStateText(self: *Renderer, gs: *GameState) !void {
        const roundBuf = try std.fmt.bufPrint(scratchBuf(50), "round: {}", .{gs.round});
        var row: usize = 0;
        self.canvas.writeText(.{
            .row = row,
            .col = self.textOffset,
        }, roundBuf, .{.r = 255, .g = 255, .b = 255});

        row += 1;

        const elapsed = try std.fmt.bufPrint(scratchBuf(50), "time: {s}", .{try utils.humanTime(gs.time)});
        self.canvas.writeText(.{
            .row = row,
            .col = self.textOffset,
        }, elapsed, .{.r = 255, .g = 255, .b = 255});

        if (engine.gamestate.waitingForTowers(gs)) {
            const towersRemaining = try std.fmt.bufPrint(scratchBuf(50), "mode: Tower Selecting", .{});
            self.canvas.writeText(.{
                .row = row,
                .col = self.textOffset,
            }, towersRemaining, .{.r = 255, .g = 255, .b = 255});
            row += 1;

            const oneRemaining = try std.fmt.bufPrint(scratchBuf(50), "one: {}", .{gs.oneAvailableTower});
            self.canvas.writeText(.{
                .row = row,
                .col = self.textOffset,
            }, oneRemaining, .{.r = 255, .g = 255, .b = 255});
            row += 1;

            const twoRemaining = try std.fmt.bufPrint(scratchBuf(50), "two: {}", .{gs.twoAvailableTower});
            self.canvas.writeText(.{
                .row = row,
                .col = self.textOffset,
            }, twoRemaining, .{.r = 255, .g = 255, .b = 255});
            row += 1;

        } else {
            const playing = try std.fmt.bufPrint(scratchBuf(50), "mode: Playing", .{});
            self.canvas.writeText(.{
                .row = row,
                .col = self.textOffset,
            }, playing, .{.r = 255, .g = 255, .b = 255});
            row += 1;
            self.canvas.writeText(.{
                .row = row,
                .col = self.textOffset,
            }, "", .{.r = 255, .g = 255, .b = 255});
            row += 1;
            self.canvas.writeText(.{
                .row = row,
                .col = self.textOffset,
            }, "", .{.r = 255, .g = 255, .b = 255});
            row += 1;
        }
    }

    pub fn completed(self: *Renderer, gs: *GameState) !void {
        const roundBuf = try std.fmt.bufPrint(scratchBuf(50), "round: {}", .{gs.round});
        var row: usize = 0;
        self.canvas.writeText(.{
            .row = row,
            .col = 0,
        }, roundBuf, .{.r = 255, .g = 255, .b = 255});

        row += 1;

        const elapsed = try std.fmt.bufPrint(scratchBuf(50), "time: {s}", .{try utils.humanTime(gs.time)});
        self.canvas.writeText(.{
            .row = row,
            .col = 0,
        }, elapsed, .{.r = 255, .g = 255, .b = 255});

        var winner = PLAYER_ONE;
        if (gs.oneTowerCount == 0) {
            winner = PLAYER_TWO;
        }

        row += 1;

        const winnerTxt = try std.fmt.bufPrint(scratchBuf(50), "winner: {s}", .{winner});
        self.canvas.writeText(.{
            .row = row,
            .col = 0,
        }, winnerTxt, .{.r = 255, .g = 255, .b = 255});

        try self.canvas.render();
        self.output = self.canvas.renderBuffer;
        self.count += 1;
    }
};
