const assert = @import("assert");

pub const Updateable = struct {
    pub const VTable = struct {
        update: *const fn(*anyopaque, delta: u64) void,
        id: *const fn(*anyopaque) usize,
    };

    ptr: *anyopaque,
    vtable: *const VTable,

    pub fn update(self: Updateable, delta: u64) void {
        self.vtable.update(self.ptr, delta);
    }

    pub fn id(self: Updateable) void {
        self.vtable.id(self.ptr);
    }

};

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
    pub const VTable = struct {
        render: *const fn(*anyopaque) Rendered,
        id: *const fn(*anyopaque) usize,
        z: *const fn(*anyopaque) usize,
    };

    ptr: *anyopaque,
    vtable: *const VTable,

    pub fn render(self: Renderable) Rendered {
        return self.vtable.render(self.ptr);
    }

    pub fn id(self: Renderable) usize {
        return self.vtable.id(self.ptr);
    }

    pub fn z(self: Renderable) usize {
        return self.vtable.z(self.ptr);
    }
};
