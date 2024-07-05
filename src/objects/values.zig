const assert = @import("assert").assert;

rows: usize = 0,
cols: usize = 0,
size: usize = 0,

const Self = @This();

pub fn init(v: *Self) void {
    assert(v.rows > 0, "must set rows");
    assert(v.cols > 0, "must set cols");

    v.size = v.rows * v.cols;
}
