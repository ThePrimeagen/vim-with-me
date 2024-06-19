const assert = @import("assert").assert;
const std = @import("std");
const print = std.debug.print;

const time = @import("engine/time.zig");
const canvas = @import("engine/canvas.zig");
const input = @import("engine/input/input.zig");
const output = @import("engine/output/output.zig");
const gamestate = @import("engine/game-state.zig");
const encoding = @import("encoding/encoding.zig");

const Message = input.Message;
const Coord = input.Coord;
const NextRound = input.NextRound;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var gs = gamestate.GameState.init(allocator);
    var stdin = input.StdinInputter.init();
    var stdinInputter = stdin.inputter();
    const inputter = try input.createInputRunner(allocator, &stdinInputter);
    //const out = output.Stdout.output;

    var fps = time.FPS.init(166_666);
    _ = fps.delta();

    while (true) {
        fps.sleep();
        const delta = fps.delta();
        gs.updateTime(delta);

        const msgInput = inputter.pop();
        if (msgInput) |msg| {
            gs.message(msg);
        }

        //const len = try ansi.frame(&cells, &buffer);
        //try out(buffer[0..len]);

    }

}

test { _ = encoding; }
test { _ = time; }
test { _ = output; }
test { _ = canvas; }

