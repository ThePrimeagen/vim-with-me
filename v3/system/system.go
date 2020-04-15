package system

import (
    "vim-with-me/scheduled"
)

type SystemCommand struct {
    Program string
    OnArgs []string
    OffArgs []string
    Scheduled scheduled.Scheduled
}

func createCommand(prog string, on []string, off []string, seconds int64) SystemCommand {
    c := SystemCommand{
        prog,
        on,
        off,
        scheduled.NewScheduled(seconds),
    }

    return c
}

func CreateASDF() SystemCommand {
    return createCommand(
        "setxkbmap",
        []string{"us"},
        []string{"us", "real-prog-dvorak"},
        0,
    )
}

func CreateXrandr() SystemCommand {
    return createCommand(
        "xrandr",
        []string{"--output", "HDMI-0", "--brightness", "0.05"},
        []string{"--output", "HDMI-0", "--brightness", "1"},
        0,
    )
}

