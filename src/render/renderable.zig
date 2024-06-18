const assert = @import("assert").assert;

pub const Location = struct {
    row: usize,
    col: usize,
};

pub const Rendered = struct {
    loc: Location,
    data: []const u8,
    cols: usize,
};

pub const Renderable = struct {
    ptr: *anyopaque,
    vtab: *const VTab,

    const VTab = struct {
        render: *const fn(ptr: *anyopaque) Rendered,
        id: *const fn(ptr: *anyopaque) usize,
        z: *const fn(ptr: *anyopaque) usize,
    };

    pub fn render(self: Renderable) Rendered {
        return self.vtab.render(self.ptr);
    }

    pub fn id(self: Renderable) usize {
        return self.vtab.id(self.ptr);
    }

    pub fn z(self: Renderable) usize {
        return self.vtab.z(self.ptr);
    }

    pub fn init(t: anytype) Renderable {
        const Ptr = @TypeOf(t);
        const PtrInfo = @typeInfo(Ptr);

        comptime assert(Ptr == .Pointer, "you must provide a pointer");
        comptime assert(PtrInfo.Pointer.size == .One, "it must be a pointer to one element");
        comptime assert(@typeInfo(PtrInfo.Pointer.child) == .Struct, "the pointer is pointing to a struct");

        const impl = struct {
            fn render(ptr: *anyopaque) Rendered {
                const self: Ptr = @ptrCast(@alignCast(ptr));
                return self.render();
            }
            fn id(ptr: *anyopaque) usize {
                const self: Ptr = @ptrCast(@alignCast(ptr));
                return self.id();
            }
            fn z(ptr: *anyopaque) usize {
                const self: Ptr = @ptrCast(@alignCast(ptr));
                return self.z();
            }
        };

        return .{
            .ptr = t,
            .vtab = &.{
                .render = impl.render,
                .id = impl.id,
                .z = impl.z,
            },
        };
    }
};

