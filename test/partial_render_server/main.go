package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"chat.theprimeagen.com/pkg/testies"
	"chat.theprimeagen.com/pkg/window"
)

func main() {
    server, win, err := testies.CreateServerFromArgs()
    if err != nil {
        log.Fatalf("Error creating server: %s", err)
    }

    ticker := time.NewTicker(500 * time.Millisecond)
    server.ToSockets.Welcome(window.OpenCommand(win))

    count := 0
    for {
        count++
        fmt.Printf("Count: %d -- waiting for ticker\n", count)
        Outer:
        for {
            select {
            case <-ticker.C:
                break Outer
            case command := <-server.FromSockets:
                fmt.Printf("Got command from socket: %+v\n", command)
            }
        }

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
            fmt.Printf("Sending: %+v\n", render.Command())
            server.ToSockets.Spread(render.Command())
            fmt.Printf("   sent\n")
        }

    }
}

