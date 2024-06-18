const std = @import("std");
const print = std.debug.print;

const render = @import("engine/render.zig");
const engine = @import("engine/engine.zig");
const input = @import("engine/input/input.zig");
const output = @import("engine/output/output.zig");
const encoding = @import("encoding/encoding.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var e = engine.Engine.init(allocator);

    var stdin = input.StdinInputter.init();
    var stdinInputter = stdin.inputter();
    const inputter = try input.createInputRunner(allocator, &stdinInputter);
    var outputter = output.init(3, 3);

    //while (!game.isDone()) {
    while (true) {
        while (true) {
            const msg = inputter.pop();
            if (msg) |_| {
            } else {
                break;
            }
        }

        e.gameLoop();
        // render
        outputter.frame();
    }

}

test { _ = encoding; }
test { _ = render; }
test { _ = engine; }

