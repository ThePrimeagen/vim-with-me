const std = @import("std");
const print = std.debug.print;

const render = @import("engine/render.zig");
const engine = @import("engine/engine.zig");
const input = @import("engine/input/input.zig");
const encoding = @import("encoding/encoding.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    _ = engine.Engine.init(allocator);
    // input
    // game
    // output
    //

    //while (!game.isDone()) {
    //while (true) {
        // read input
        // play input into game
        // update game
        //e.gameLoop();
        // output
    //}

    //const in = std.io.getStdIn();
    //var buf = std.io.bufferedReader(in.reader());
    //const reader = buf.reader().any();

    var stdin = input.StdinInputter.init();
    var stdinInputter = stdin.inputter();
    const inputter = try input.createInputRunner(allocator, &stdinInputter);

    while (true) {
        const msg = inputter.pop();

        if (msg) | m | {
            std.debug.print("msg: {s}\n", .{m.input[0..m.length]});
        } else {
            std.debug.print("no message :(\n", .{});
        }

        std.time.sleep(1000 * 1000 * 1000);
    }
}

test { _ = encoding; }
test { _ = render; }
test { _ = engine; }

