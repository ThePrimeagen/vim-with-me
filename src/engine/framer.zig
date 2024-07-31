const objects = @import("../objects/objects.zig");
const a = @import("../assert/assert.zig");
const std = @import("std");

const assert = a.assert;
const never = a.never;
const io = std.io;
const colors = objects.colors;
const Color = colors.Color;
const Cell = colors.Cell;
const Values = objects.Values;

const initialChar: u8 = '';
const topLeftFull: [6]u8 = .{'', '[', '1', ';', '1', 'H'};
const clear: [4]u8 = .{'', '[', '2', 'J'};
const newFrame: [4]u8 = .{ '', '[', ';', 'H', };
const initialClear: [14]u8 = topLeftFull ++ clear ++ newFrame;
const foregroundColor: [7]u8 = .{ '', '[', '3', '8', ';', '2', ';', };
const backgroundColor: [7]u8 = .{ '', '[', '4', '8', ';', '2', ';', };
const backgroundClear: [5]u8 = .{ '', '[', '4', '9', 'm' };
const newline: [2]u8 = .{ '\r', '\n', };
const clearStyle = .{'', '[', '0', 'm'};

const DEFAULT_CELL: Cell = .{.text=' ', .color = undefined, .background = null};

pub fn writeAnsiColor(color: Color, ansiColor: []const u8, out: []u8, o: usize) !usize {
    var offset = o;

    const maxOffset = offset + 12 + ansiColor.len;
    assert(maxOffset < out.len, "potential buffer overflow");

    offset = write(out, offset, ansiColor);
    const slice = try std.fmt.bufPrint(out[offset..], "{};{};{}m", .{color.r, color.g, color.b});

    offset += slice.len;

    return offset;
}

pub fn write(out: []u8, o: usize, bytes: []const u8) usize {
    assert(out.len > o + bytes.len, "buffer overflowed");

    var offset = o;
    for (bytes) |b| {
        out[offset] = b;
        offset += 1;
    }

    return offset;
}

pub fn writeByte(out: []u8, offset: usize, byte: u8) usize {
    assert(out.len > offset, "buffer overflowed");
    out[offset] = byte;
    return offset + 1;
}

pub fn ansiCodeLength(buffer: []const u8) usize {
    const items: []const []const u8 = &.{
        &topLeftFull,
        &clear,
        &newFrame,
        &clearStyle,
    };

    for (items) |item| {
        if (std.mem.indexOf(u8, buffer, item)) |idx| {
            if (idx == 0) {
                return item.len;
            }
        }
    }

    if (std.mem.indexOf(u8, buffer, &foregroundColor)) |idx| {
        if (idx == 0) {
            if (std.mem.indexOf(u8, buffer, &.{'m'})) |end| {
                return end + 1;
            }
        }
    }

    never("couldn't find code");
    return 0;
}

pub fn ansiBackgroundClear(buffer: []u8, offset: usize) !usize {
    return write(buffer, offset, &backgroundClear);
}

pub const AnsiFramer = struct {
    firstPrint: bool,
    previous: Cell = DEFAULT_CELL,
    values: *const Values,

    pub fn init(values: *const Values) AnsiFramer {
        return .{
            .firstPrint = true,
            .values = values,
        };
    }

    pub fn frame(self: *AnsiFramer, f: []const Cell, out: []u8) !usize {
        assert(f.len == self.values.size, "you must hand in a frame that matches rows and cols");
        var hasBackground = self.previous.background != null;

        var offset: usize = 0;
        if (self.firstPrint) {
            offset = write(out, offset, &initialClear);
            self.firstPrint = false;
        } else {
            offset = write(out, offset, &newFrame);
        }

        var newLineCount: usize = 0;
        for (f, 1..) |*c, idx| {
            const text = c.text;

            if (!self.previous.sameColors(c.*)) {
                self.previous = c.*;
                offset = try writeAnsiColor(c.color, &foregroundColor, out, offset);
                if (c.background) |b| {
                    offset = try writeAnsiColor(b, &backgroundColor, out, offset);
                    hasBackground = true;
                } else if (hasBackground) {
                    offset = try ansiBackgroundClear(out, offset);
                    hasBackground = false;
                }
            }

            offset = writeByte(out, offset, text);

            if (idx % self.values.cols == 0) {
                offset = write(out, offset, &newline);
                newLineCount += 1;
            }
        }

        // clear everything at the end
        offset = write(out, offset, "\x1b[0m");
        assert(newLineCount == self.values.rows, "should have produced self.rows amount of rows, did not");
        self.previous = DEFAULT_CELL;
        return offset;
    }

    pub fn parseText(buf: []const u8, out: []u8) usize {
        var idx: usize = 0;
        var i: usize = 0;

        while (i < buf.len) {
            if (buf[i] == initialChar) {
                i += ansiCodeLength(buf[i..]);
                continue;
            }

            out[idx] = buf[i];
            idx += 1;
            i += 1;
        }

        return idx;
    }
};

const testing = std.testing;
test "AnsiFramer should color and newline a 3x3" {
    var values = Values{.rows = 3, .cols = 3};
    values.init();

    var frame = AnsiFramer.init(&values);
    var out = [1]u8{0} ** 100;
    var colors1 = [9]Cell{
        .{.text = 'a', .color = .{.r = 69, .g = 42, .b = 0}},
        .{.text = 'b', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'c', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'd', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'e', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'f', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'g', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'h', .color = .{.r = 70, .g = 43, .b = 1}},
        .{.text = 'i', .color = .{.r = 71, .g = 44, .b = 2}},
    };

    const len1 = try frame.frame(&colors1, &out);

    const expected =
        initialClear ++
        foregroundColor ++
        "69;42;0ma".* ++
        foregroundColor ++
        "70;43;1mbc\r\ndef\r\ngh".* ++
        foregroundColor ++
        "71;44;2mi\r\n".* ++
        clearStyle;

    try testing.expectEqualSlices(u8, &expected, out[0..len1]);

    var text: [9 + 6]u8 = .{' '} ** 15;

    _ = AnsiFramer.parseText(out[0..len1], &text);
    try testing.expectEqualStrings("abc\r\ndef\r\nghi\r\n", &text);

    var colors2 = [_]Cell{
        .{.text = 'i', .color = .{.r = 71, .g = 44, .b = 2}}
    } ** 9;

    const expected2 =
        newFrame ++
        foregroundColor ++
        "71;44;2miii\r\niii\r\niii\r\n".* ++
        clearStyle;

    const len2 = try frame.frame(&colors2, &out);
    try testing.expectEqualSlices(u8, &expected2, out[0..len2]);

    _ = AnsiFramer.parseText(out[0..len2], &text);
    try testing.expectEqualStrings("iii\r\niii\r\niii\r\n", &text);
}
