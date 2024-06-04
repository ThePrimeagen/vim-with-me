const std = @import("std");
const print = std.debug.print;

const render = @import("engine/render.zig");
const engine = @import("engine/engine.zig");
const encoding = @import("encoding/encoding.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    const e = engine.Engine.init(allocator);
    // input
    // game
    // output
    //

    //while (!game.isDone()) {
    while (true) {
        // read input
        // play input into game
        // update game
        e.gameLoop();
        // output
    }
}

test { _ = encoding; }
test { _ = render; }
test { _ = engine; }

