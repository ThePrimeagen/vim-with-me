const Output = @import("stdout_output.zig");
const Framer = @import("framer.zig");

pub const Stdout = Output;
pub const Cell = Framer.Cell;
pub const Color = Framer.Color;
pub const AnsiFramer = Framer.AnsiFramer;

test { _ = Framer; }
