pub const Color = struct {
    r: u8,
    g: u8,
    b: u8,

    pub fn equal(self: Color, other: Color) bool {
        return self.r == other.r and
            self.g == other.g and
            self.b == other.b;
    }

};

pub const Black: Color = .{.r = 0, .g = 0, .b = 0 };
pub const Red: Color = .{.r = 255, .g = 0, .b = 0 };
pub const DarkDarkGrey: Color = .{.r = 25, .g = 25, .b = 25 };
pub const DarkGrey: Color = .{.r = 50, .g = 50, .b = 50 };
pub const Grey: Color = .{.r = 100, .g = 100, .b = 100 };
pub const Green: Color = .{.r = 0, .g = 255, .b = 0 };
pub const LightGrey: Color = .{.r = 175, .g = 175, .b = 175 };
pub const White: Color = .{.r = 255, .g = 255, .b = 255 };
pub const Blue: Color = .{ .r = 0x3f, .g = 0xa9, .b = 0xff, };
pub const Orange: Color = .{ .r = 245, .g = 164, .b = 66, };

pub const Cell = struct {
    text: u8,
    color: Color,
    background: ?Color = null,

    pub fn sameColors(self: Cell, other: Cell) bool {
        return self.color.equal(other.color) and (
            (self.background == null and other.background == null) or
            (self.background != null and other.background != null and
             self.background.?.equal(other.background.?))
        );
    }
};


