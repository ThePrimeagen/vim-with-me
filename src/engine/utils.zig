const std = @import("std");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;

pub const SECOND: isize = 1000 * 1000;
pub const MINUTE: isize = 60 * SECOND;

pub fn humanTime(timeUS: isize) ![]u8 {
    const minutes = @divFloor(timeUS, MINUTE);
    const seconds = @divFloor(@mod(timeUS, MINUTE), SECOND);

    return std.fmt.bufPrint(scratchBuf(50), "{}m {}s", .{minutes, seconds});
}

pub fn normalize(comptime t: type, current: t, total: t) f64 {
    const cF: f64 = @floatFromInt(current);
    const tF: f64 = @floatFromInt(total);

    return cF / tF;
}

test "displayTime" {
    const out = try humanTime(69 * 1000 * 1000);
    try std.testing.expectEqualStrings("1m 9s", out);
}
