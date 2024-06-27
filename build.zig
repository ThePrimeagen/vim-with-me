const std = @import("std");

// Although this function looks imperative, note that its job is to
// declaratively construct a build graph that will be executed by an external
// runner.
pub fn build(b: *std.Build) void {
    // Standard target options allows the person running `zig build` to choose
    // what target to build for. Here we do not override the defaults, which
    // means any target is allowed, and the default is native. Other options
    // for restricting supported target set are available.
    const target = b.standardTargetOptions(.{});

    // Standard optimization options allow the person running `zig build` to select
    // between Debug, ReleaseSafe, ReleaseFast, and ReleaseSmall. Here we do not
    // set a preferred release mode, allowing the user to decide how to optimize.
    const optimize = b.standardOptimizeOption(.{});

    const assert = b.addModule("assert", .{
        .root_source_file = .{ .path = "src/assert/assert.zig" },
    });

    const scratch = b.addModule("scratch", .{
        .root_source_file = .{ .path = "src/scratch/scratch.zig" },
    });

    const objects = b.addModule("objects", .{
        .root_source_file = .{ .path = "src/objects/objects.zig" },
    });

    const math = b.addModule("math", .{
        .root_source_file = .{ .path = "src/math/math.zig" },
    });

    const testing = b.addModule("testing", .{
        .root_source_file = .{ .path = "src/test/test.zig" },
    });

    const vengine = b.addModule("vengine", .{
        .root_source_file = .{ .path = "src/engine/engine.zig" },
    });

    const exe = b.addExecutable(.{
        .name = "to",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    const testExe = b.addExecutable(.{
        .name = "test_to",
        .root_source_file = b.path("src/test/main.zig"),
        .target = target,
        .optimize = optimize,
    });


    exe.root_module.addImport("assert", assert);
    exe.root_module.addImport("scratch", scratch);
    scratch.addImport("assert", assert);

    exe.root_module.addImport("vengine", vengine);
    vengine.addImport("assert", assert);
    vengine.addImport("scratch", scratch);
    vengine.addImport("math", math);
    vengine.addImport("objects", objects);

    exe.root_module.addImport("testing", testing);
    testing.addImport("assert", assert);
    testing.addImport("scratch", scratch);
    testing.addImport("math", math);
    testing.addImport("objects", objects);
    testing.addImport("vengine", vengine);

    exe.root_module.addImport("math", math);
    math.addImport("assert", assert);
    math.addImport("scratch", scratch);

    exe.root_module.addImport("objects", objects);
    objects.addImport("math", math);
    objects.addImport("assert", assert);
    objects.addImport("scratch", scratch);

    testExe.root_module.addImport("testing", testing);
    testExe.root_module.addImport("assert", assert);
    testExe.root_module.addImport("scratch", scratch);
    testExe.root_module.addImport("math", math);
    testExe.root_module.addImport("objects", objects);
    testExe.root_module.addImport("vengine", vengine);

    {
        // This declares intent for the executable to be installed into the
        // standard location when the user invokes the "install" step (the default
        // step when running `zig build`).
        b.installArtifact(exe);

        // This *creates* a Run step in the build graph, to be executed when another
        // step is evaluated that depends on it. The next line below will establish
        // such a dependency.
        const run_cmd = b.addRunArtifact(exe);

        // By making the run step depend on the install step, it will be run from the
        // installation directory rather than directly from within the cache directory.
        // This is not necessary, however, if the application depends on other installed
        // files, this ensures they will be present and in the expected location.
        run_cmd.step.dependOn(b.getInstallStep());

        // This allows the user to pass arguments to the application in the build
        // command itself, like this: `zig build run -- arg1 arg2 etc`
        if (b.args) |args| {
            run_cmd.addArgs(args);
        }

        // This creates a build step. It will be visible in the `zig build --help` menu,
        // and can be selected like this: `zig build run`
        // This will evaluate the `run` step rather than the default, which is "install".
        const run_step = b.step("run", "Run the app");
        run_step.dependOn(&run_cmd.step);
    }

    {
        // This declares intent for the executable to be installed into the
        // standard location when the user invokes the "install" step (the default
        // step when running `zig build`).
        b.installArtifact(testExe);

        // This *creates* a Run step in the build graph, to be executed when another
        // step is evaluated that depends on it. The next line below will establish
        // such a dependency.
        const run_cmd = b.addRunArtifact(testExe);

        // By making the run step depend on the install step, it will be run from the
        // installation directory rather than directly from within the cache directory.
        // This is not necessary, however, if the application depends on other installed
        // files, this ensures they will be present and in the expected location.
        run_cmd.step.dependOn(b.getInstallStep());

        // This allows the user to pass arguments to the application in the build
        // command itself, like this: `zig build run -- arg1 arg2 etc`
        if (b.args) |args| {
            run_cmd.addArgs(args);
        }

        // This creates a build step. It will be visible in the `zig build --help` menu,
        // and can be selected like this: `zig build run`
        // This will evaluate the `run` step rather than the default, which is "install".
        const run_step = b.step("simulate", "runs the simulation (which we are in)");
        run_step.dependOn(&run_cmd.step);
    }

    const exe_unit_tests = b.addTest(.{
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    exe_unit_tests.root_module.addImport("assert", assert);
    exe_unit_tests.root_module.addImport("objects", objects);
    exe_unit_tests.root_module.addImport("math", math);
    exe_unit_tests.root_module.addImport("scratch", scratch);
    exe_unit_tests.root_module.addImport("testing", testing);
    exe_unit_tests.root_module.addImport("vengine", vengine);

    const test_exe_unit_tests = b.addTest(.{
        .root_source_file = b.path("src/test/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    test_exe_unit_tests.root_module.addImport("assert", assert);
    test_exe_unit_tests.root_module.addImport("objects", objects);
    test_exe_unit_tests.root_module.addImport("math", math);
    test_exe_unit_tests.root_module.addImport("scratch", scratch);
    test_exe_unit_tests.root_module.addImport("testing", testing);
    test_exe_unit_tests.root_module.addImport("vengine", vengine);

    const run_exe_unit_tests = b.addRunArtifact(exe_unit_tests);
    run_exe_unit_tests.has_side_effects = true;

    const test_run_exe_unit_tests = b.addRunArtifact(test_exe_unit_tests);
    test_run_exe_unit_tests.has_side_effects = true;

    // Similar to creating the run step earlier, this exposes a `test` step to
    // the `zig build --help` menu, providing a way for the user to request
    // running the unit tests.
    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_exe_unit_tests.step);

    const test_test_step = b.step("test-sim", "Run simulations");
    test_test_step.dependOn(&test_run_exe_unit_tests.step);
}
