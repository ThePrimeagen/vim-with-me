const std = @import("std");

const assert = @import("assert").assert;
const math = @import("math");
const scratchBuf = @import("scratch").scratchBuf;

const colors = @import("colors.zig");

const Color = colors.Color;
const Cell = colors.Cell;
const Red = colors.Red;
const Allocator = std.mem.Allocator;

const INITIAL_CREEP_LIFE = 10;
const INITIAL_CREEP_SPEED = 1;
const INITIAL_CREEP_COLOR: Color = .{.r = 0, .g = 0, .b = 0};

pub const CreepSize = 1;
pub const CreepCell: [1]Cell = .{
    .{.text = '*', .color = Red },
};

pub const Creep = struct {
    id: usize,
    team: u8,
    cols: usize,

    pos: math.Vec2 = math.ZERO_VEC2,
    life: u16 = INITIAL_CREEP_LIFE,
    speed: f32 = INITIAL_CREEP_SPEED,
    alive: bool = true,

    // rendered
    rLife: u16 = INITIAL_CREEP_LIFE,
    rColor: Color = INITIAL_CREEP_COLOR,
    rCells: [1]Cell = CreepCell,
    rSized: math.Sized = math.ZERO_SIZED,

    path: []usize,
    pathIdx: usize = 0,
    pathLen: usize = 0,
    alloc: Allocator,

    pub fn string(self: *Creep) ![]u8 {
        const buf = scratchBuf(150);
        return std.fmt.bufPrint(buf, "creep({}, {}, {})\r\n  pos = {s}\r\n  pathIdx = {}, life = {}, speed = {}\r\n", .{
            self.alive, self.id, self.team,
            try self.pos.string(),
            self.pathIdx, self.life, self.speed,
        });
    }

    pub fn init(alloc: Allocator, rows: usize, cols: usize) !Creep {
        return .{
            .path = try alloc.alloc(usize, rows * cols),
            .cols = cols,
            .alloc = alloc,
            .id = 0,
            .team = 0,
        };
    }

    pub fn deinit(self: *Creep) void {
        self.alloc.free(self.path);
    }
};

