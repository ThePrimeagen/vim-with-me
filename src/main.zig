const std = @import("std");
const print = std.debug.print;

const engine = @import("engine/engine.zig");
const input = @import("engine/input/input.zig");
const output = @import("engine/output/output.zig");
const encoding = @import("encoding/encoding.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    //const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    //var e = engine.Engine.init(allocator);

    //var stdin = input.StdinInputter.init();
    //var stdinInputter = stdin.inputter();
    //const inputter = try input.createInputRunner(allocator, &stdinInputter);
    var ansi = output.AnsiFramer.init(3, 3);
    const out = output.Stdout.output;
    const colors: [3]output.Color = .{
        .{ .r = 255, .g = 0, .b = 0 },
        .{ .g = 255, .r = 0, .b = 0 },
        .{ .b = 255, .r = 0, .g = 0 },
    };

    var fps = engine.FPS.init(166_666);

    //while (!game.isDone()) {
    var count: usize = 0;
    while (true) : (count += 1) {
        const delta = fps.delta();
        _ = delta;

        //const msgInput = inputter.pop();
        //if (msgInput == null) {
        //    continue;
        //}

        var buffer = [_]u8{0} ** 100;
        var cells = [_]output.Cell{
            .{.text = 'a', .color = undefined},
        } ** 9;
        for (0..9) |i| {
            cells[i].color = colors[count % 3];
        }

        const len = try ansi.frame(&cells, &buffer);
        try out(buffer[0..len]);

        fps.sleep();
    }

}

test { _ = encoding; }
test { _ = engine; }
test { _ = output; }

