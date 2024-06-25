const assert = @import("assert").assert;
const std = @import("std");
const print = std.debug.print;

const objects = @import("objects");
const math = @import("math");

const renderer = @import("engine/renderer.zig");
const time = @import("engine/time.zig");
const canvas = @import("engine/canvas.zig");
const input = @import("engine/input/input.zig");
const framer = @import("engine/framer.zig");
const stdout = @import("engine/stdout_output.zig");
const encoding = @import("encoding/encoding.zig");

const GameState = objects.gamestate.GameState;
const Message = objects.message.Message;

const Coord = input.Coord;
const NextRound = input.NextRound;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var gs = GameState.init(allocator);
    var stdin = input.StdinInputter.init();
    var stdinInputter = stdin.inputter();

    const inputter = try input.createInputRunner(allocator, &stdinInputter);
    defer inputter.deinit();

    var render = try renderer.Renderer.init(30, 30, allocator);
    defer render.deinit();

    const out = stdout.output;

    var fps = time.FPS.init(166_666);
    _ = fps.delta();

    while (true) {
        fps.sleep();
        const delta = fps.delta();
        gs.updateTime(delta);

        const msgInput = inputter.pop();
        if (msgInput) |msg| {
            if (Message.init(msg.input[0..msg.length])) |p| {
                try gs.message(p);
            }
        }

        try render.render(&gs);
        try out(render.output);
    }

}

test { _ = encoding; }
test { _ = time; }
test { _ = framer; }
test { _ = canvas; }
test { _ = objects; }
test { _ = math; }

