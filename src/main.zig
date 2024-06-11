const assert = @import("assert").assert;
const std = @import("std");
const print = std.debug.print;

const engine = @import("engine/engine.zig");
const input = @import("engine/input/input.zig");
const output = @import("engine/output/output.zig");
const encoding = @import("encoding/encoding.zig");

const needle: [1]u8 = .{','};
const Coord = struct {
    row: usize,
    col: usize,

    pub fn init(str: []const u8) !Coord {
        const idxMaybe = std.mem.indexOf(u8, str, ",");
        assert(idxMaybe != null, "must find , for input");

        const idx = idxMaybe.?;
        const row = try std.fmt.parseInt(usize, str[0..idx], 10);
        const col = try std.fmt.parseInt(usize, str[idx + 1..], 10);

        return .{
            .row = row,
            .col = col,
        };
    }
};

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    //var e = engine.Engine.init(allocator);

    var stdin = input.StdinInputter.init();
    var stdinInputter = stdin.inputter();
    const inputter = try input.createInputRunner(allocator, &stdinInputter);
    var ansi = output.AnsiFramer.init(10, 10);
    const out = output.Stdout.output;
    const colors: [3]output.Color = .{
        .{ .r = 255, .g = 0, .b = 0 },
        .{ .g = 255, .r = 0, .b = 0 },
        .{ .b = 255, .r = 0, .g = 0 },
    };

    var fps = engine.FPS.init(166_666);
    _ = fps.delta();

    //while (!game.isDone()) {
    var count: usize = 0;
    while (true) : (count += 1) {
        fps.sleep();
        const delta = fps.delta();
        _ = delta;

        const msgInput = inputter.pop();
        if (msgInput == null) {
            continue;
        }

        const msg = msgInput.?;
        const coord = try Coord.init(msg.input[0..msg.length]);

        var buffer = [_]u8{0} ** 4096;
        var cells = [_]output.Cell{
            .{.text = ' ', .color = .{.r = 0, .g = 0, .b = 0}},
        } ** 100;

        cells[coord.row * 10 + coord.col].color = colors[count % 3];
        cells[coord.row * 10 + coord.col].text = 'X';

        const len = try ansi.frame(&cells, &buffer);
        try out(buffer[0..len]);

    }

}

test { _ = encoding; }
test { _ = engine; }
test { _ = output; }

