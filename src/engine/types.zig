const std = @import("std");

const needle: [1]u8 = .{','};
pub const Position = struct {
    row: usize,
    col: usize,

    pub fn init(str: []const u8) ?Position {
        const idxMaybe = std.mem.indexOf(u8, str, ",");
        if (idxMaybe == null) {
            return null;
        }

        const idx = idxMaybe.?;
        const row = std.fmt.parseInt(usize, str[0..idx], 10) catch {
            return null;
        };
        const col = std.fmt.parseInt(usize, str[idx + 1..], 10) catch {
            return null;
        };

        return .{
            .row = row,
            .col = col,
        };
    }
};

pub const Coord = struct {
    pos: Position,
    team: u8,

    pub fn init(msg: []const u8) ?Coord {
        const teamNumber = msg[0];
        if (teamNumber != '1' and teamNumber != '2') {
            return null;
        }
        const pos = Position.init(msg[1..]);
        if (pos == null) {
            return null;
        }

        return .{
            .pos = pos.?,
            .team = teamNumber,
        };
    }
};

pub const NextRound = struct {
    pub fn init(msg: []const u8) ?NextRound {
        if (msg.len == 1 and msg.ptr[0] == 'n') {
            return .{ };
        }
        return null;
    }
};

pub const Message = union(enum) {
    round: NextRound,
    coord: Coord,

    pub fn init(msg: []const u8) ?Message {
        const coord = Coord.init(msg);
        if (coord) |c| {
            return .{.coord = c};
        }

        const next = NextRound.init(msg);
        if (next) |n| {
            return .{.round = n};
        }

        return null;
    }
};

pub const Color = struct {
    r: u8,
    g: u8,
    b: u8,
};

pub const Tower = struct {
    id: usize,
    team: u8,

    // position
    pos: Position,
    ammo: u16,
    dead: bool,
    level: u8,
    radius: u8,
    damage: u8,

    // rendered
    rPos: Position,
    rAmmo: u16,
    rColor: [9]Color,
    rText: [9]u8,

    fn contains(self: *Tower, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        const r = @abs(self.row - pos.row);
        const c = @abs(self.col - pos.col);
        return r <= 1 and c <= 1;
    }

};

pub const Projectile = struct {
    id: usize,
    createdBy: usize,
    team: u8,

    pos: Position,
    life: u16,
    speed: f32,
    dead: bool,

    // rendered
    rPos: Position,
    rLife: u16,
    rColor: Color,
    rText: u8,
};

pub const Creep = struct {
    id: usize,
    team: u8,

    pos: Position,
    life: u16,
    speed: f32,
    dead: bool,

    // rendered
    rPos: Position,
    rLife: u16,
    rColor: Color,
    rText: u8,

    fn contains(self: *Creep, pos: Position) bool {
        if (self.dead) {
            return false;
        }

        const r = @abs(self.row - pos.row);
        const c = @abs(self.col - pos.col);
        return r == 0 and c == 0;
    }
};
