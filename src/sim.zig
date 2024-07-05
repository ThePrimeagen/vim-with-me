const std = @import("std");
const testing = @import("test/test.zig");
const objects = @import("objects/objects.zig");
const engine = @import("engine/engine.zig");
const assert = @import("assert/assert.zig");
const math = @import("math/math.zig");
const Values = objects.Values;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const alloc = gpa.allocator();
    var args = try testing.Params.readFromArgs(alloc);
    var values = args.values();

    var gs = try objects.gamestate.GameState.init(alloc, &values);
    defer gs.deinit();

    engine.gamestate.init(&gs);

    const gsDump = gs.dumper();
    assert.addDump(&gsDump);

    const out = engine.stdout.output;

    var fps: ?engine.time.FPS = null;
    if (args.realtime.?) {
        fps = engine.time.FPS.init(args.fps);
        _ = fps.?.delta();
    }

    var render = try engine.renderer.Renderer.init(alloc, &values);
    defer render.deinit();
    var creeper = testing.gamestate.Spawner.init(&args, &gs);

    var count: usize = 0;
    while (args.runCount > count) : (count += 1) {
        var delta = args.fps;
        if (fps) |*f| {
            f.sleep();
            delta = f.delta();
        }

        try creeper.tick(delta);

        engine.stdout.resetColor();
        try engine.gamestate.update(&gs, delta);

        if (args.viz.?) {
            try render.render(&gs);
            try out(render.output);
        }

        engine.gamestate.validateState(&gs);
    }

}
