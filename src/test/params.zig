const Values = @import("../objects/objects.zig").Values;
const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const std = @import("std");
const assert = @import("../assert/assert.zig").assert;

const Allocator = std.mem.Allocator;
const RndGen = std.rand.DefaultPrng;

rows: usize,
cols: usize,
creepRate: usize,
fps: isize,
runCount: usize = 10000,

seed: ?usize = 0,
viz: ?bool = true,
realtime: ?bool = false,

_rand: ?RndGen = null,

const Self = @This();
pub fn readFromArgs(alloc: Allocator) !Self {
    var args = try std.process.argsWithAllocator(alloc);
    _ = args.next();
    const pathMaybe = args.next();
    assert(pathMaybe != null, "there must be arguments");

    const path = pathMaybe.?;
    const self = try readConfig(alloc, path);
    defer self.deinit();

    var out = self.value;
    out._rand = RndGen.init(@intCast(out.seed.?));

    return out;
}

fn readConfig(allocator: Allocator, path: []const u8) !std.json.Parsed(Self) {
    const data = try std.fs.cwd().readFileAlloc(allocator, path, 1024);
    defer allocator.free(data);
    return std.json.parseFromSlice(Self, allocator, data, .{.allocate = .alloc_always});
}

pub fn string(self: *const Self) ![]u8 {
    return std.fmt.bufPrint(scratchBuf(150), "rows = {}, cols = {}, creepRate = {}, towerCount = {}, viz = {}, realtime = {}", .{
        self.rows, self.cols, self.creepRate, self.towerCount,
        self.viz.?,
        self.realtime.?,
    });
}

pub fn rand(self: *Self, comptime T: type) T {
    return self._rand.?.random().int(T);
}

pub fn randRange(self: *Self, comptime T: type, start: T, end: T) T {
    assert(start < end, "end must be greater than start");
    return start + (end - start) % self._rand.?.random().int(T);
}

pub fn values(self: *const Self) Values {
    var v = Values{
        .rows = self.rows,
        .cols = self.cols,
    };

    Values.init(&v);

    return v;
}
