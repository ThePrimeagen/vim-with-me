const std = @import("std");

const print = std.debug.print;

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

var dumps: [10]?*const Dump = .{null} ** 10;
var dumpsIdx: usize = 0;

pub fn addDump(dump: *const Dump) void {
    assert(dumpsIdx < dumps.len, "you have added too many dumps, please readjust the hard coded size");
    dumps[dumpsIdx] = dump;
    dumpsIdx += 1;
}

pub fn unwrap(comptime T: type, val: anyerror!T) T {
    if (val) |v| {
        return v;
    } else |err| {
        std.debug.panic("unwrap error: {any}", .{err});
    }
}

pub fn u(v: anyerror![]u8) []u8 {
    return unwrap([]u8, v);
}

pub fn assert(truthy: bool, msg: []const u8) void {
    if (!truthy) {
        std.debug.print("\x1b[0m", .{});
        _ = std.io.getStdOut().write("\x1b[0m") catch |e| {
            std.debug.print("error while clearing stdout: {}\n", .{e});
        };

        for (0..dumpsIdx) |d| {
            dumps[d].?.dump();
        }

        @panic(msg);
    }
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

