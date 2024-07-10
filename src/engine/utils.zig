const std = @import("std");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;

const SECOND: isize = 1000 * 1000;
const MINUTE: isize = 60 * SECOND;

pub fn humanTime(timeUS: isize) ![]u8 {
    const minutes = @divFloor(timeUS, MINUTE);
    const seconds = @divFloor(@mod(timeUS, MINUTE), SECOND);

    return std.fmt.bufPrint(scratchBuf(50), "{}m {}s", .{minutes, seconds});
}

pub fn compare(expected: []const u8, value: []const u8) bool {
    if (expected.len != value.len) {
        return false;
    }

    for (expected, 0..) |c, i| {
        if (c != value[i]) {
            return false;
        }
    }
    return true;
}

test "displayTime" {
    const out = try humanTime(69 * 1000 * 1000);
    try std.testing.expectEqualStrings("1m 9s", out);
}
