const assert = @import("assert");
const std = @import("std");

const Updateable = struct {
    ptr: *anyopaque,
    vtab: *const VTab,

    const VTab = struct {
        update: *const fn(ptr: *anyopaque, delta: u64) void,
        id: *const fn(ptr: *anyopaque) usize,
    };

    pub fn update(self: Updateable) void {
        self.vtab.update();
    }

    pub fn id(self: Updateable) usize {
        return self.vtab.id(self.ptr);
    }

    pub fn init(item: anytype) Updateable {

        const Ptr = @TypeOf(item);
        const PtrInfo = @typeInfo(Ptr);

        comptime assert(PtrInfo == .Pointer, "expected a pointer to be passed in");
        comptime assert(PtrInfo.Pointer.size == .One, "expected a pointer of size one");
        comptime assert(PtrInfo.Pointer.child == .Struct, "expected the pointer to point to a struct");

        const impl = struct {
            pub fn update(ptr: *anyopaque) void {
                const self: Ptr = @ptrCast(@alignCast(ptr));
                self.update();
            }

            fn id(ptr: *anyopaque) usize {
                const self: Ptr = @ptrCast(@alignCast(ptr));
                return self.id();
            }
        };

        return .{
            .ptr = item,
            .vtab = &.{
                .update = impl.update,
                .id = impl.id,
            },
        };
    }
};

const UpdateList = std.ArrayList(*Updateable);
pub const Updater = struct {
    items: UpdateList,

    pub fn init(alloc: std.mem.Allocator) Updater {
        return .{
            .items = UpdateList.init(alloc),
        };
    }

    pub fn deinit(self: *Updater) void {
        self.items.deinit();
    }

    pub fn update(self: *Updater, delta: u64) void {
        for (self.items) |*updateable| {
            updateable.update(delta);
        }
    }
};

