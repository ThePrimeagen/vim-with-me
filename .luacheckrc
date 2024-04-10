std = luajit
cache = true
codes = true
ignore = {
    "111", -- Setting an undefined global variable. (for ok, _ = pcall...)
    "211", -- Unused local variable.
    "411", -- Redefining a local variable.
}
read_globals = { "vim", "describe", "it", "bit", "assert", "before_each", "after_each" }


