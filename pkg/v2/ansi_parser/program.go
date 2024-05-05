package ansiparser

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/theprimeagen/vim-with-me/pkg/chat"
	chat2 "theprimeagen.com/test/pkg/chat"
)

type DoomWriter struct{}

func (d *DoomWriter) Write(b []byte) (int, error) {
	fmt.Printf("I have read: %d\n", len(b))
	return len(b), nil
}

type Stats struct {
	msgs int

	executed map[string]int
	raw      map[string]int
}

type AnsiProgram struct {
    path string
}

func NewAnsiProgram(name string) *AnsiProgram {
    return &AnsiProgram{
        path: name,
    }
}

func main() {
	outFile, err := os.Create("/tmp/test")
	stats := Stats{
        msgs: 0,
        executed: map[string]int{},
        raw: map[string]int{},
    }

	defer outFile.Close()
	defer func() {
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		fmt.Printf("stats %+v\n", stats)
		slog.Warn("game stats", "stats", stats)
	}()

	if err != nil {
		log.Fatal("could not initialize out file", err)
	}

	logger := slog.New(slog.NewTextHandler(outFile, nil))
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cmd := exec.Command("./doom_ascii/doom_ascii", "-iwad", "./DOOM.WAD", "-scaling", "2")

	ctx := context.Background()
	ch, err := chat.NewTwitchChat(ctx)
	agg := chat2.NewChatAggregator(func(msg string) bool {
		switch msg {
		case "w", "a", "s", "d", "1", "2", "3", "use", "u", "fire", "f":
			slog.Debug("filter#valid", "msg", msg)
            stats.raw[msg]++
			return true
		}
		slog.Debug("filter#invalid", "msg", msg)
		return false
	}).WithMap(func(msg string) string {
		slog.Debug("map", "msg", msg, "to", strings.ToLower(msg))
		return strings.ToLower(msg)
	})

	if err != nil {
		slog.Error("unable to initialize twitch chat", "error", err)
		return
	}

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal("error at start", err)
		return
	}

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	go func() {
		<-time.After(time.Second * 1)
		ptmx.Write([]byte{'
'})
		<-time.After(time.Second * 1)
		ptmx.Write([]byte{'
'})
		<-time.After(time.Second * 1)
		ptmx.Write([]byte{'
'})
        <-time.After(time.Second * 1)
        ptmx.Write([]byte{'w'})
		<-time.After(time.Second * 1)
		ptmx.Write([]byte{'
'})
		<-time.After(time.Second * 1)
		ptmx.Write([]byte{'
'})
		slog.Warn("game loop#started")

		for {
			slog.Warn("game loop#Waiting for input")
			timer := time.NewTimer(time.Millisecond * 125)
		outer:
			for {
				select {
				case <-timer.C:
					break outer
				case msg := <-ch:
					agg.Add(msg.Msg)
				}
			}

			occ := agg.Reset()
			if occ.Count == 0 {
				continue
			}

			slog.Warn("game loop#Got Input", "input", occ)
			slog.Warn("game stats", "stats", stats)

            stats.executed[occ.Msg]++

            playCount := 10

            // todo didn't think, stupid
		playLoop:
			for i := 0; i < playCount; i++ {
				switch occ.Msg {
				case "w", "s":
					ptmx.Write([]byte(occ.Msg))
				case "a", "d":
					ptmx.Write([]byte(occ.Msg))
				case "1", "2", "3":
					ptmx.Write([]byte(occ.Msg))
					break playLoop
				case "fire", "f":
					ptmx.Write([]byte{'f'})
					break playLoop
				case "use", "u":
					ptmx.Write([]byte{'u'})
					break playLoop
				}
				<-time.After(time.Millisecond * 16)
			}

			// <-time.After(time.Second * 5)
		}

	}()

	//writer := DoomWriter{}
	_, _ = io.Copy(os.Stdout, ptmx)
}
*/


