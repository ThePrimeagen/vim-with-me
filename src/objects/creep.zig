const std = @import("std");

const assert = @import("assert").assert;
const math = @import("math");

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
    rPos: math.Position = math.ZERO_POS,
    rLife: u16 = INITIAL_CREEP_LIFE,
    rColor: Color = INITIAL_CREEP_COLOR,
    rCells: [1]Cell = CreepCell,
    rSized: math.Sized = math.ZERO_SIZED,

    scratch: []isize,
    path: []usize,
    pathIdx: usize = 0,
    pathLen: usize = 0,
    alloc: Allocator,

    pub fn string(self: *Creep, buf: []u8) !usize {
        var out = try std.fmt.bufPrint(buf, "creep({}, {})\r\n", .{self.id, self.team});
        var len = out.len;

        out = try std.fmt.bufPrint(buf[len..], "  pos = ", .{});
        len += out.len;
        len += try self.pos.string(buf[len..]);
        out = try std.fmt.bufPrint(buf[len..], "  pathIdx = {}\r\nlife = {}\r\n  speed = {}\r\n  alive = {}\n\n)", .{self.pathIdx, self.life, self.speed, self.alive});
        return len + out.len;
    }

    pub fn init(alloc: Allocator, rows: usize, cols: usize) !Creep {
        return .{
            .path = try alloc.alloc(usize, rows * cols),
            .scratch = try alloc.alloc(isize, rows * cols),
            .cols = cols,
            .alloc = alloc,
            .id = 0,
            .team = 0,
        };
    }

    pub fn deinit(self: *Creep) void {
        self.alloc.free(self.path);
        self.alloc.free(self.scratch);
    }
};

