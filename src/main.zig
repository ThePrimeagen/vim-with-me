const std = @import("std");
const print = std.debug.print;

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
    //var outputter = output.init(3, 3);
    var time = engine.RealTime.init();
    time.reset();

    //while (!game.isDone()) {
    while (true) {
        const delta = time.tick();
        while (true) {
            const msg = inputter.pop();
            if (msg) |_| {
            } else {
                break;
            }
        }

        e.gameLoop(delta);
        // render
        // outputter.frame();
    }

}

test { _ = encoding; }
test { _ = engine; }
test { _ = output; }

