const std = @import("std");
const testing = @import("testing");
const objects = @import("objects");
const engine = @import("vengine");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const alloc = gpa.allocator();
    var args = try testing.params.readFromArgs(alloc);
    const values = args.values();

    var gs = try objects.gamestate.GameState.init(alloc, &values);
    defer gs.deinit();
    const out = engine.stdout.output;

    var fps: ?engine.time.FPS = null;
    if (args.realtime.?) {
        fps = engine.time.FPS.init(args.fps);
        _ = fps.?.delta();
    }

    var render = try engine.renderer.Renderer.init(alloc, &values);
    defer render.deinit();
    var creeper = testing.gamestate.Spawner.init(&args, &gs);

    for (0..args.towerCount) |_| {

        while (true) {
            const row = args.rand(usize) % args.rows;
            const col = args.rand(usize) % args.cols;
            const pos = .{.row = row, .col = col};
            if (engine.gamestate.canPlaceTower(&gs, pos)) {
                try engine.gamestate.placeTower(&gs, pos, 0);
                break;
            }
        }

    }

    while (true) {
        var delta = args.fps;
        if (fps) |*f| {
            f.sleep();
            delta = f.delta();
        }

        try creeper.tick(delta);
        engine.gamestate.update(&gs, delta);

        if (args.viz.?) {
            try render.render(&gs);
            try out(render.output);
        }
    }
}
