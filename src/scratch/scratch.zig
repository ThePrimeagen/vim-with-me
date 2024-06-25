const assert = @import("assert").assert;

var scratch: [8192]u8 = [_]u8{0} ** 8192;
var idx: usize = 0;

pub fn scratchBuf(size: usize) []u8 {
    assert(size < scratch.len, "you cannot require more than the size of the scratch buffer");
    if (idx + size > scratch.len) {
        idx = 0;
    }

    const out = scratch[idx..idx + size];
    idx += out.len;
    return out;
}

