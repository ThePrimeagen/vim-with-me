const assert = @import("assert/assert.zig").assert;
const std = @import("std");
const print = std.debug.print;

const objects = @import("objects/objects.zig");
const math = @import("math/math.zig");

const engine = @import("engine/engine.zig");

const GameState = objects.gamestate.GameState;
const Message = objects.message.Message;

const Coord = engine.input.Coord;
const NextRound = engine.input.NextRound;

pub fn main() !void {
    engine.stdout.hideCursor();
    try engine.stdout.showCursorOnSigInt();

    var values = objects.Values{};
    values.rows = 24;
    values.cols = 80;
    values.seed = 42069;
    values.fps = 33_333;
    values.realtimeMultiplier = 3;
    values.init();

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var gs = try GameState.init(allocator, &values);
    defer gs.deinit();
    engine.gamestate.init(&gs);

    var stdin = engine.input.StdinInputter.init();
    var stdinInputter = stdin.inputter();

    const inputter = try engine.input.createInputRunner(allocator, &stdinInputter);
    defer inputter.deinit();

    var render = try engine.renderer.Renderer.init(allocator, &values);

    defer render.deinit();

    const out = engine.stdout.output;

    var fps = engine.time.FPS.init(values.fps);
    _ = fps.delta();

    // TODO: Figure out something better with the creep spawner... this sucks
    var spawner = engine.rounds.CreepSpawner.init(&gs);
    var reportState = engine.reportState.ReportState{};

    try reportState.waiting(&gs);

    while (!engine.gamestate.completed(&gs)) {
        fps.sleep();
        const delta = fps.delta();
        var multipliedDelta: isize = @intFromFloat(@as(f64, @floatFromInt(delta)) * values.realtimeMultiplier);

        const msgInput = inputter.pop();
        if (msgInput) |msg| {
            if (Message.init(msg.input[0..msg.length])) |p| {
                try engine.gamestate.message(&gs, p);
            }
        }

        if (engine.gamestate.hasActiveCreeps(&gs)) {
            while (multipliedDelta > 0) {
                const innerDelta = @min(multipliedDelta, delta);
                try engine.gamestate.update(&gs, innerDelta);
                try spawner.tick();
                multipliedDelta -= innerDelta;
            }
        } else if (gs.playing) {
            engine.gamestate.endRound(&gs);

            // TODO: this sucks... but it rocks...
            try reportState.waiting(&gs);

        } else if (!engine.gamestate.waitingForTowers(&gs)) {
            engine.gamestate.startRound(&gs, &spawner);
            try reportState.playing();
        }

        try render.render(&gs);
        try out(render.output);

        engine.gamestate.validateState(&gs);
    }

    try render.completed(&gs);
    try out(render.output);
    engine.stdout.showCursor();
}

test { _ = objects; }
test { _ = math; }
test { _ = engine; }

