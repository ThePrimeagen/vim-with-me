const std = @import("std");

const assert = @import("../assert/assert.zig").assert;
const Values = @import("values.zig");
const math = @import("../math/math.zig");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;

const colors = @import("colors.zig");

const Color = colors.Color;
const Cell = colors.Cell;
const Red = colors.Red;
const Allocator = std.mem.Allocator;

const INITIAL_CREEP_COLOR: Color = .{ .r = 0, .g = 0, .b = 0 };

pub const CreepSize = 1;
pub const CreepCell: [1]Cell = .{
    .{ .text = '*', .color = Red },
};

pub const Creep = struct {
    id: usize,
    team: u8,
    values: *const Values,

    pos: math.Vec2 = math.ZERO_VEC2,
    life: usize = 0,
    speed: f64 = 0,
    alive: bool = true,

    // rendered
    rLife: usize = 0,
    rColor: Color = INITIAL_CREEP_COLOR,
    rCells: [1]Cell = CreepCell,
    rSized: math.Sized = math.ZERO_SIZED,

    scratch: []isize,
    path: []usize,
    pathIdx: usize = 0,
    pathLen: usize = 0,
    alloc: Allocator,

    pub fn format(
        self: Creep,
        comptime fmt: []const u8,
        options: std.fmt.FormatOptions,
        writer: anytype,
    ) !void {
        _ = fmt;
        _ = options;

        try writer.print("creep({}, {}, {})\r\n  pos = {}\r\n  path = {}/{}, life = {}, speed = {d}\r\n", .{
            self.alive, self.id,      self.team,
            self.pos,   self.pathIdx, self.pathLen,
            self.life,  self.speed,
        });
    }

    pub fn init(alloc: Allocator, values: *const Values) !Creep {
        return .{
            .values = values,
            .path = try alloc.alloc(usize, values.size),
            .scratch = try alloc.alloc(isize, values.size),
            .alloc = alloc,
            .id = 0,
            .team = 0,

            .life = values.creep.life,
            .rLife = values.creep.life,
            .speed = values.creep.speed,
        };
    }

    pub fn deinit(self: *Creep) void {
        self.alloc.free(self.path);
        self.alloc.free(self.scratch);
    }
};
