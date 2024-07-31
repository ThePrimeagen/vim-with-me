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
pub const Slate: Color = .{.r = 50, .g = 50, .b = 50 };
pub const DarkSlate: Color = .{.r = 10, .g = 10, .b = 10 };
pub const White: Color = .{.r = 255, .g = 255, .b = 255 };
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


