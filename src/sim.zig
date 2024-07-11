const std = @import("std");
const Params = @import("test/params.zig");
const utils = @import("test/utils.zig");
const objects = @import("objects/objects.zig");
const engine = @import("engine/engine.zig");
const assert = @import("assert/assert.zig");
const math = @import("math/math.zig");
const Values = objects.Values;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const alloc = gpa.allocator();
    var args = try Params.readFromArgs(alloc);
    var values = args.values();

    var gs = try objects.gamestate.GameState.init(alloc, &values);
    defer gs.deinit();

    engine.gamestate.init(&gs);
    engine.gamestate.endRound(&gs);

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

    var count: usize = 0;
    while (args.runCount > count) : (count += 1) {
        var delta = args.fps;
        if (fps) |*f| {
            f.sleep();
            delta = f.delta();
        }

        engine.stdout.resetColor();
        if (engine.gamestate.hasActiveCreeps(&gs)) {
            //try creeper.tick();
            try engine.gamestate.update(&gs, delta);
        } else {
            engine.gamestate.endRound(&gs);

            const cnt = engine.rounds.towerCount(&gs);
            engine.gamestate.setTowerPlacementCount(&gs, cnt);

            for (0..cnt) |_| {
                const one = utils.positionInRange(&gs, Values.TEAM_ONE);
                const two = utils.positionInRange(&gs, Values.TEAM_TWO);
                try engine.gamestate.message(&gs, .{ .coord = .{ .pos = one, .team = Values.TEAM_ONE, }, });
                try engine.gamestate.message(&gs, .{ .coord = .{ .pos = two, .team = Values.TEAM_TWO, }, });
            }

            engine.gamestate.startRound(&gs);
        }

        if (args.viz.?) {
            try render.render(&gs);
            try out(render.output);
        }

        engine.gamestate.validateState(&gs);

        if (engine.gamestate.completed(&gs)) {
            try render.completed(&gs);
            try out(render.output);
            break;
        }
    }

}
