const std = @import("std");

const print = std.debug.print;

fn clear() void {
    std.debug.print("\x1b[?25h", .{});
    std.debug.print("\x1b[0m", .{});
    _ = std.io.getStdOut().write("\x1b[?25h") catch |e| {
        std.debug.print("error while clearing stdout: {}\n", .{e});
    };
    _ = std.io.getStdOut().write("\x1b[0m") catch |e| {
        std.debug.print("error while clearing stdout: {}\n", .{e});
    };
}

pub const Dump = struct {
    ptr: *anyopaque,
    vtab: *const VTab,
    const VTab = struct {
        dump: *const fn (ptr: *anyopaque) void,
    };

    pub fn dump(self: Dump) void {
        self.vtab.dump(self.ptr);
    }

    // cast concrete implementation types/objs to interface
    pub fn init(obj: anytype) Dump {
        const Ptr = @TypeOf(obj);
        const PtrInfo = @typeInfo(Ptr);

        assert(PtrInfo == .Pointer, "must pass in pointer"); // Must be a pointer
        assert(PtrInfo.Pointer.size == .One, "Pointer to one object"); // Must be a single-item pointer
        assert(@typeInfo(PtrInfo.Pointer.child) == .Struct, "must be pointing to a struct"); // Must point to a struct

        const impl = struct {
            fn dump(ptr: *anyopaque) void {
                const self: Ptr = @ptrCast(@alignCast(ptr));
                self.dump();
            }
        };

        return .{
            .ptr = obj,
            .vtab = &.{
                .dump = impl.dump,
            },
        };
    }
};

const dumpLen = 10;
var dumps: [dumpLen]?*const Dump = .{null} ** dumpLen;

pub fn addDump(dump: *const Dump) void {
    for (0..dumpLen) |idx| {
        if (dumps[idx] == null) {
            dumps[idx] = dump;
            return;
        }
    }
    assert(false, "could not add dumps, you have exceeded the hardcoded size");
}

pub fn removeDump(dump: *const Dump) void {
    for (0..dumpLen) |idx| {
        if (dumps[idx] == dump) {
            dumps[idx] = null;
            return;
        }
    }
    assert(false, "could not find dump");
}

pub fn unwrap(comptime T: type, val: anyerror!T) T {
    if (val) |v| {
        return v;
    } else |err| {
        std.debug.panic("unwrap error: {any}", .{err});
    }
}

pub fn never(msg: []const u8) void {
    clear();

    for (0..dumpLen) |dump| {
        if (dumps[dump]) |d| {
            d.dump();
        }
    }

    @panic(msg);
}

pub fn option(comptime T: type, val: ?T) T {
    if (val) |v| {
        return v;
    }
    never("option is null");
    unreachable;
}

pub fn u(v: anyerror![]u8) []u8 {
    return unwrap([]u8, v);
}

pub fn assert(truthy: bool, msg: []const u8) void {
    if (truthy) {
        return;
    }

    never(msg);
}

// TODO: DO SOMETHING WITH THIS...
pub fn printZZZ(toPrint: anytype) ![]u8 {
    const MyType = @TypeOf(toPrint);
    const hasStr = @hasDecl(MyType, "string");

    if (hasStr) {
        return toPrint.string();
    }
    return .{};
}

