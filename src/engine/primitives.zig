const assert = @import("assert");

pub const Updateable = struct {
    pub const VTable = struct {
        update: *const fn(*Updateable, delta: u64) void,
    };

    vtable: *const VTable,

    pub fn update(self: *Updateable, delta: u64) void {
        self.vtable.update(self, delta);
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
        render: *const fn(*Renderable) Rendered,
        deinit: *const fn(*Renderable) void,
    };

    vtable: *const VTable,
    id: usize,
    z: usize,

    pub fn init(vtable: *const VTable, id: usize, z: usize) Renderable {
        return .{
            .z = z,
            .id = id,
            .vtable = vtable,
        };
    }

    pub fn render(self: *Renderable) Rendered {
        return self.vtable.render(self);
    }

    pub fn deinit(self: *Renderable) void {
        self.vtable.deinit(self);
    }
};
