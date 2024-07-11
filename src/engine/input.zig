const assert = @import("../assert/assert.zig").assert;
const std = @import("std");
const Allocator = std.mem.Allocator;

// ASSUMPTION: no input should be bigger than 256 characters
pub const InputSize = 4096;
pub const Input = struct {
    input: [InputSize]u8,
    length: usize,
};
pub const InputList = std.ArrayList(Input);

pub const BufferedInputter = struct {
    lines: [][]const u8,
    idx: usize,

    pub fn init() BufferedInputter {
        return .{
            .lines = .{},
            .idx = 0,
        };
    }

    pub fn data(self: *BufferedInputter, lines: [][]const u8) void {
        self.lines = lines;
        self.idx = 0;
    }

    pub fn next(self: *BufferedInputter, buf: []u8) !?usize {
        if (self.idx >= self.lines.len) {
            return 0;
        }

        const n = self.lines[self.idx].len;
        @memcpy(buf, self.lines[self.idx]);
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

    pub fn deinit(self: *InputRunner) void {
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
        while (true) {
            var buf: [InputSize]u8 = undefined;
            const n = try inputter.next(&buf) orelse break;

            {
                self.mutex.lock();
                defer self.mutex.unlock();
                try self.inputs.append(Input{
                    .input = buf,
                    .length = n,
                });
            }
        }
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

test "buffered input" {
}
