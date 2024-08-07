const a = @import("../assert/assert.zig");
const std = @import("std");
const Params = @import("../params.zig");
const objects = @import("../objects/objects.zig");
const engine = @import("../engine/engine.zig");
const math = @import("../math/math.zig");

const assert = a.assert;
const never = a.never;
const Self = @This();
const GS = objects.gamestate.GameState;
const Message = objects.message.Message;

fh: std.fs.File,

pub fn next(self: *Self, gs: *GS) !void {
    var in_stream = self.fh.reader();

    var buf: [1024]u8 = undefined;
    var count: usize = 0;
    const total: usize = @intCast(engine.gamestate.getTotalTowerPlacement(gs));

    for (0..total) |_| {
        const lineMaybe = try in_stream.readUntilDelimiterOrEof(&buf, '\n');
        if (lineMaybe == null) {
            break;
        }
        const line = lineMaybe.?;
        count += 1;
        const msg = try Message.init(line);
        if (msg) |m| {
            try engine.gamestate.message(gs, m);
        } else {
            std.debug.print("{s}\n", .{line});
            never("coord from file was null");
        }
    }

    assert(count == total, "did not execute the expected amount of tower creations");
}

pub fn deinit(self: *Self) void {
    self.fh.close();
}

pub fn fromParams(params: *const Params) !Self {
    assert(params.simulationType != null, "simulationType cannot be null");

    const st = params.simulationType.?;
    if (std.mem.eql(u8, st, "stdin")) {
        return .{
            .fh = std.io.getStdIn(),
        };
    }
    assert(std.mem.startsWith(u8, st, "file:"), "expected simulationType to have \"file:\" beginning");

    const path = st[5..];
    return .{
        .fh = try std.fs.cwd().openFile(path, .{}),
    };
}
