const assert = @import("../assert/assert.zig").assert;
const std = @import("std");
const print = std.debug.print;

const RLEError = error{
    RLETooBig,
};

pub const RLE = struct {
    out: ?[]u8,
    idx: usize,
    curr: u8,
    count: u8,

    const Self = @This();

    fn init() Self {
        return .{
            .out = null,
            .idx = 0,
            .curr = 0,
            .count = 0,
        };
    }

    pub fn reset(self: *Self, out: []u8) void {
        self.out = out;
        self.idx = 0;
        self.curr = 0;
        self.count = 0;
    }

    pub fn write(self: *Self, in: []const u8) !void {
        var start: usize = 0;
        if (self.count == 0) {
            self.curr = in[0];
            self.count = 1;
            start = 1;
        }

        for (in[start..]) |byte| {
            if (self.curr == byte and self.count < 255) {
                self.count += 1;
                continue;
            }
            try self.place();

            self.count = 1;
            self.curr = byte;
        }
    }

    pub fn finish(self: *Self) !void {
        assert(self.count != 0, "you have called finish without every writing data");
        try self.place();
    }

    fn place(self: *Self) !void {
        assert(self.count != 0, "you cannot place a 0 occurrence character");
        assert(self.out != null, "you cannot RLE without calling reset");

        var out = self.out.?;

        if (out.len <= self.idx + 1) {
            return RLEError.RLETooBig;
        }

        out[self.idx] = self.count;
        out[self.idx + 1] = self.curr;

        self.idx += 2;
    }
};

test "rle encodes wwwwjd" {
    var rle = RLE.init();
    var buffer: [6]u8 = [_]u8{0} ** 6;

    rle.reset(&buffer);
    try rle.write("wwwwjd");
    try rle.finish();

    const out: []const u8 = rle.out.?;

    try std.testing.expectEqualSlices(u8, out, &.{ 4, 'w', 1, 'j', 1, 'd' });
}

test "rle encodes should fail on encoding something too large" {
    var rle = RLE.init();
    var buffer: [6]u8 = [_]u8{0} ** 6;

    rle.reset(&buffer);
    try rle.write("wwwwjdx");
    const err = rle.finish();

    try std.testing.expectError(RLEError.RLETooBig, err);
}

test "rle encodes past the byte overflow amount" {
    var rle = RLE.init();
    var buffer: [6]u8 = [_]u8{0} ** 6;

    rle.reset(&buffer);
    try rle.write("w" ** 256);
    try rle.place();

    const out: []const u8 = rle.out.?;

    try std.testing.expectEqualSlices(u8, out, &.{255, 'w', 1, 'w', 0, 0});
}
