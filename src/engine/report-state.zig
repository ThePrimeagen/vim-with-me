const std = @import("std");
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const objects = @import("../objects/objects.zig");

const GS = objects.gamestate.GameState;

const NONE = "none";
const PLAYING = "playing\n";
const WAITING_TO_PLAY = "waiting";

pub const ReportState = struct {
    state: []const u8 = NONE,

    pub fn playing(self: *ReportState) !void {
        if (std.mem.eql(u8, self.state, PLAYING)) {
            return;
        }

        self.state = PLAYING;
        _ = try std.io.getStdErr().write(self.state);
    }

    pub fn waiting(self: *ReportState, gs: *GS) !void {
        if (std.mem.eql(u8, self.state, WAITING_TO_PLAY)) {
            return;
        }

        self.state = WAITING_TO_PLAY;
        const contents = try std.fmt.bufPrint(scratchBuf(50), "{s}-{}\n", .{WAITING_TO_PLAY, gs.oneAvailableTower});
        _ = try std.io.getStdErr().write(contents);
    }

    pub fn final(self: *ReportState, gs: *GS) !void {
        _ = self;
        const err = std.io.getStdErr();

        try err.writeAll("final\n");
        if (gs.oneTowerCount == 0) {
            try err.writeAll("1: lost\n");
            try err.writeAll("2: won\n");
        } else {
            try err.writeAll("1: won\n");
            try err.writeAll("2: lost\n");
        }
    }
};


