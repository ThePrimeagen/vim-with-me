package main

import (
	"log"
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func main() {
    server, err := testies.CreateServerFromArgs()
    if err != nil {
        log.Fatalf("Error creating server: %s", err)
    }

    win := window.NewWindow(80, 24)
    server.ToSockets.Welcome(window.OpenCommand(win))

    count := 0
    for {
        count++
        <-server.FromSockets

        number := count % 10
        row := count / win.Cols
        col := count % win.Cols

        // number to rune
        numStr := strconv.Itoa(number)
        numRune := []rune(numStr)[0]
        err := win.Set(row, col, numRune)
        if err != nil {
            log.Fatalf("Error setting rune: %s", err)
        }

        renders := win.Flush()
        for _, render := range renders {
            server.ToSockets.Spread(render.Command())
        }

    }
}

