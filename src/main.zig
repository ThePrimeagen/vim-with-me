const assert = @import("assert").assert;
const std = @import("std");
const print = std.debug.print;

const objects = @import("objects");
const math = @import("math");

const renderer = @import("engine/renderer.zig");
const engine = @import("engine/engine.zig");

// TODO: Remove encoding
const encoding = @import("encoding/encoding.zig");

const GameState = objects.gamestate.GameState;
const Message = objects.message.Message;

const Coord = engine.input.Coord;
const NextRound = engine.input.NextRound;

const ROWS = 30;
const COLS = 30;

pub fn main() !void {
    var values = objects.Values{};
    values.rows = 30;
    values.cols = 80;
    values.init();

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var gs = try GameState.init(allocator, &values);
    defer gs.deinit();

    var stdin = engine.input.StdinInputter.init();
    var stdinInputter = stdin.inputter();

    const inputter = try engine.input.createInputRunner(allocator, &stdinInputter);
    defer inputter.deinit();

    var render = try engine.renderer.Renderer.init(allocator, &values);

    defer render.deinit();

    const out = engine.stdout.output;

    var fps = engine.time.FPS.init(166_666);
    _ = fps.delta();

    try engine.gamestate.placeCreep(&gs, .{
        .row = 0,
        .col = 0,
    });

    while (true) {
        fps.sleep();
        const delta = fps.delta();

        const msgInput = inputter.pop();
        if (msgInput) |msg| {
            if (Message.init(msg.input[0..msg.length])) |p| {
                try engine.gamestate.message(&gs, p);
            }
        }

        engine.gamestate.update(&gs, delta);

        try render.render(&gs);
        try out(render.output);
    }

}

test { _ = encoding; }
test { _ = objects; }
test { _ = math; }
test { _ = engine; }

