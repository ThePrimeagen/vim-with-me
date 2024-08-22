const Values = @import("objects/objects.zig").Values;
const scratchBuf = @import("scratch/scratch.zig").scratchBuf;
const std = @import("std");
const a = @import("assert/assert.zig");
const assert = a.assert;
const never = a.never;

const Allocator = std.mem.Allocator;

rows: usize,
cols: usize,
fps: isize,
runCount: usize = 10000,
simCount: usize = 100,

seed: ?usize = 0,
viz: ?bool = true,
realtime: ?bool = false,
realtimeMultiplier: f64 = 1,
simulationType: ?[]u8 = null,

const Self = @This();
pub fn readFromArgs(alloc: Allocator) !Self {
    var args = try std.process.argsWithAllocator(alloc);
    _ = args.next();
    const pathMaybe = args.next();
    assert(pathMaybe != null, "there must be arguments");

    const path = pathMaybe.?;
    var self = try readConfig(alloc, path);
    defer self.deinit();

    // TODO: This could be a lot smarter....
    while (args.next()) |k| {
        if (std.mem.eql(u8, "--seed", k)) {
            const v = try std.fmt.parseInt(usize, args.next().?, 10);
            self.value.seed = v;
        }
    }

    if (self.value.simulationType) |sim| {
        const str = try alloc.alloc(u8, sim.len);
        @memcpy(str, sim);
        self.value.simulationType = str;
    }

    return self.value;
}

pub fn deinit(self: *Self, alloc: Allocator) void {
    if (self.simulationType) |s| {
        alloc.free(s);
    }
}

fn readConfig(allocator: Allocator, path: []const u8) !std.json.Parsed(Self) {
    const data = try std.fs.cwd().readFileAlloc(allocator, path, 1024);
    defer allocator.free(data);
    return std.json.parseFromSlice(Self, allocator, data, .{.allocate = .alloc_always});
}

pub fn string(self: *const Self) ![]u8 {
    return std.fmt.bufPrint(scratchBuf(150), "rows = {}, cols = {}, towerCount = {}, viz = {}, realtime = {}", .{
        self.rows, self.cols, self.towerCount,
        self.viz.?,
        self.realtime.?,
    });
}

pub fn values(self: *const Self) Values {
    var v = Values{
        .rows = self.rows,
        .cols = self.cols,
        .seed = self.seed orelse 42069,
        .realtimeMultiplier = self.realtimeMultiplier,
    };

    Values.init(&v);

    return v;
}
