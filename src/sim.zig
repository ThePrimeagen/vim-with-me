const std = @import("std");
const Params = @import("test/params.zig");
const utils = @import("test/utils.zig");
const objects = @import("objects/objects.zig");
const engine = @import("engine/engine.zig");
const assert = @import("assert/assert.zig");
const math = @import("math/math.zig");
const rand = @import("sim/rand.zig");
const scratchBuf = @import("scratch/scratch.zig").scratchBuf;
const simulation = @import("sim/sim.zig");

const never = assert.never;
const Values = objects.Values;
const Allocator = std.mem.Allocator;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const alloc = gpa.allocator();
    var args = try Params.readFromArgs(alloc);

    var inGameTime: i64 = 0;
    for (0..args.simCount) |_| {
        args.seed.? += 1;

        const timings = try runSimulation(alloc, &args);
        std.debug.print("{}\n", .{timings.rounds});

        inGameTime += timings.time;
    }
}

const Timings = struct {
    rounds: usize,
    time: i64,

    pub fn string(self: Timings) ![]u8 {
        return std.fmt.bufPrint(scratchBuf(150), "round = {}, = time = {s}", .{self.rounds, try engine.utils.humanTime(self.time)});
    }
};

fn runSimulation(alloc: Allocator, args: *Params) !Timings {
    var values = args.values();
    var gs = try objects.gamestate.GameState.init(alloc, &values);
    var sim = try simulation.fromParams(args);

    defer gs.deinit();
    defer args.deinit(alloc);

    engine.gamestate.init(&gs);

    const gsDump = gs.dumper();
    assert.addDump(&gsDump);
    defer assert.removeDump(&gsDump);

    const out = engine.stdout.output;

    var fps: ?engine.time.FPS = null;
    if (args.realtime.?) {
        fps = engine.time.FPS.init(args.fps);
        _ = fps.?.delta();
    }

    var render = try engine.renderer.Renderer.init(alloc, &values);
    defer render.deinit();

    var spawner = engine.rounds.CreepSpawner.init(&gs);

    var count: usize = 0;
    while (args.runCount > count) : (count += 1) {
        var delta = args.fps;
        var multipliedDelta = delta;
        if (fps) |*f| {
            f.sleep();
            delta = f.delta();
            multipliedDelta = @intFromFloat(@as(f64, @floatFromInt(delta)) * args.realtimeMultiplier);
        }

        if (args.viz.?) {
            engine.stdout.resetColor();
        }

        if (engine.gamestate.hasActiveCreeps(&gs)) {
            while (multipliedDelta > 0) {
                const innerDelta = @min(multipliedDelta, delta);
                try engine.gamestate.update(&gs, innerDelta);
                try spawner.tick();

                multipliedDelta -= innerDelta;
            }
        } else {

            // TODO: Move tower count into end round and creeper spawn into
            // start round
            engine.gamestate.endRound(&gs);

            const cnt = engine.rounds.towerCount(&gs);
            engine.gamestate.setTowerPlacementCount(&gs, cnt);

            try simulation.simulate(&sim, &gs);
            engine.gamestate.startRound(&gs);

            spawner.startRound();

            // Note: Future me... remember my spawner spawns 2 creeps PER spawnCount
            const creepCount: isize = @intCast(spawner.spawnCount);
            engine.gamestate.setActiveCreeps(&gs, creepCount * 2);
        }

        if (args.viz.?) {
            try render.render(&gs);
            try out(render.output);
        }

        engine.gamestate.validateState(&gs);

        if (engine.gamestate.completed(&gs)) {
            if (args.viz.?) {
                try render.completed(&gs);
                try out(render.output);
            }
            break;
        }
    }

    return .{
        .time = gs.time,
        .rounds = gs.round,
    };
}
