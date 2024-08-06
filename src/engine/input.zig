const scratchBuf = @import("../scratch/scratch.zig").scratchBuf;
const a = @import("../assert/assert.zig");
const assert = a.assert;
const never = a.never;
const std = @import("std");
const Allocator = std.mem.Allocator;

// ASSUMPTION: no input should be bigger than 256 characters
pub const InputSize = 4096;

pub const Input = struct {
    input: [InputSize]u8,
    length: usize,

    pub fn slice(self: *const Input) []const u8 {
        return self.input[0..self.length];
    }

    pub fn string(self: *const Input) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(10 + self.length), "{}: \"{s}\"", .{self.length, self.slice()});
    }
};

pub const InputList = std.ArrayList(Input);

pub const BufferedInputter = struct {
    lines: []const []const u8 = undefined,
    idx: usize = 0,

    pub fn inputter(self: *BufferedInputter) Inputter {
        return Inputter{.buffered = self};
    }

    pub fn data(self: *BufferedInputter, lines: []const []const u8) void {
        self.lines = lines;
        self.idx = 0;
    }

    pub fn next(self: *BufferedInputter, buf: []u8) !?usize {
        if (self.idx >= self.lines.len) {
            return null;
        }

        const n = self.lines[self.idx].len;
        @memcpy(buf[0..n], self.lines[self.idx]);
        self.idx += 1;

        return n;
    }
};

pub const StdinInputter = struct {
    reader: std.io.BufferedReader(4096, std.fs.File.Reader),

    pub fn init() StdinInputter {
        const in = std.io.getStdIn();
        const buf = std.io.bufferedReader(in.reader());

        return .{
            .reader = buf
        };
    }

    pub fn inputter(self: *StdinInputter) Inputter {
        return Inputter{.stdin = self};
    }

    pub fn next(self: *StdinInputter, buf: []u8) !?usize {
        var fsb = std.io.fixedBufferStream(buf);
        try self.reader.reader().streamUntilDelimiter(fsb.writer(), '\n', buf.len);
        const len = fsb.getWritten().len;
        if (len == 0) {
            return null;
        }
        return len;
    }

};

pub const Inputter = union(enum) {
    stdin: *StdinInputter,
    buffered: *BufferedInputter,

    pub fn next(self: Inputter, buf: []u8) !?usize {
        return switch (self) {
            inline else => |s| try s.next(buf),
        };
    }
};

pub const InputRunner = struct {
    inputs: InputList,
    alloc: Allocator,
    mutex: std.Thread.Mutex,
    thread: std.Thread,
    alive: bool = true,
    running: usize = 0,

    pub fn deinit(self: *InputRunner) void {
        {
            self.mutex.lock();
            defer self.mutex.unlock();

            self.alive = false;
        }

        while (self.running > 0) { }

        self.inputs.deinit();
        self.alloc.destroy(self);
    }

    pub fn pop(self: *InputRunner) ?Input {
        if (self.inputs.items.len == 0) {
            return null;
        }

        // probably better ds here.. but whatever
        self.mutex.lock();
        defer self.mutex.unlock();
        const out = self.inputs.orderedRemove(0);

        return out;
    }

    pub fn runReadLoop(self: *InputRunner, inputter: *Inputter) !void {
        self.running += 1;
        while (true) {
            {
                self.mutex.lock();
                defer self.mutex.unlock();
                if (!self.alive) {
                    break;
                }
            }

            var buf: [InputSize]u8 = undefined;
            const n = inputter.next(&buf) catch |e| blk: {
                std.debug.print("read error: {any}\n", .{e});
                std.time.sleep(1_000_000_000);
                break :blk null;
            };

            if (n == null) {
                continue;
            }

            {
                self.mutex.lock();
                defer self.mutex.unlock();
                try self.inputs.append(Input{
                    .input = buf,
                    .length = n.?,
                });
            }
        }
        self.running -= 1;
    }

};

fn read(runner: *InputRunner, inputter: *Inputter) !void {
    try runner.runReadLoop(inputter);
}

pub fn createInputRunner(alloc: Allocator, inputter: *Inputter) !*InputRunner {
    const input = try alloc.create(InputRunner);
    input.* = .{
        .alloc = alloc,
        .inputs = InputList.init(alloc),
        .mutex = std.Thread.Mutex{},
        .thread = undefined,
    };

    input.thread = try std.Thread.spawn(.{}, read, .{input, inputter});

    return input;
}

const t = std.testing;
test "buffered input" {
    var buff = BufferedInputter{};
    var inputter = buff.inputter();
    var input = try createInputRunner(t.allocator, &inputter);
    defer input.deinit();

    var out = input.pop();
    try t.expect(out == null);

    const data: [2][]const u8 = .{
        "hello",
        "goodbye",
    };

    buff.data(data[0..2]);
    std.time.sleep(10_000_000);

    out = input.pop();
    try t.expectEqualStrings("hello", out.?.slice());

    out = input.pop();
    try t.expectEqualStrings("goodbye", out.?.slice());

    out = input.pop();
    try t.expect(out == null);
}


