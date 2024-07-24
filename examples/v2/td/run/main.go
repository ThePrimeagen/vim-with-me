package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
	"github.com/theprimeagen/vim-with-me/pkg/v2/cmd"
)

func getDebug(name string) (*testies.DebugFile, error) {
    if name == "" {
        return testies.EmptyDebugFile(), nil
    }
    return testies.NewDebugFile(name)
}

func main() {

    testies.SetupLogger()

	godotenv.Load()

	debugFile := ""
	flag.StringVar(&debugFile, "debug", "", "runs the file like the program instead of running doom")
	flag.Parse()
    debug, err := getDebug(debugFile)
    if err != nil {
        log.Fatalf("could not open up debug file: %v\n", err)
    }
    defer debug.Close()
    debug.WriteLine([]byte("hello world"))

	ctx := context.Background()

	twitchChat, err := chat.NewTwitchChat(ctx)
	assert.NoError(err, "twitch cannot initialize")
	chtAgg := chat.
		NewChatAggregator().
		WithFilter(td.TDFilter(24, 80));

    ch := make(chan struct{}, 1)

    cmdr := cmd.NewCmder("./zig-out/bin/to").
        WithErrFn(func(b []byte) (int, error) {
            str := string(b)
            parts := strings.Split(str, "-")

            debug.WriteLine([]byte(fmt.Sprintf("state: %s\n", str)))

            if parts[0] == "waiting" {
                assert.Assert(len(parts) == 2, "somehow i didn't communicate back with 2 pieces", "parts", parts)
                parts[1] = strings.TrimSpace(parts[1])

                count, err := strconv.Atoi(parts[1])
                assert.NoError(err, "zig program gave me a non number")

                for i := range count {
                    debug.WriteLine([]byte(fmt.Sprintf("getting input for item %d out of %d", i, count)))
                    ch<-struct{}{}
                    debug.WriteLine([]byte(fmt.Sprintf("got input for %d", i)))
                }
            }

            return len(b), nil
        }).
        WithOutFn(func(b []byte) (int, error) {
            fmt.Printf("Game: %s\n", string(b))
            return len(b), nil
        })

    go cmdr.Run()
	go chtAgg.Pipe(twitchChat)

    for range ch {
        <-time.NewTimer(time.Second * 10).C
        occurrences := chtAgg.ResetWithAll()
        slices.SortFunc(occurrences, func(a, b *chat.Occurrence) int {
            return b.Count - a.Count
        });

        one := false
        two := false
        for _, v := range occurrences {
            coord, err := td.FromString(v.Msg)
            if err != nil {
                continue
            }

            if !one && coord.Team == 1 {
                debug.WriteLine([]byte(fmt.Sprintf("COORD %s\n", coord.String())))
                cmdr.WriteLine([]byte(coord.String()))
                one = true
            }
            if !two && coord.Team == 2 {
                debug.WriteLine([]byte(fmt.Sprintf("COORD %s\n", coord.String())))
                cmdr.WriteLine([]byte(coord.String()))
                two = true
            }

            if one && two {
                break
            }
        }

        if !one {
            debug.WriteLine([]byte(fmt.Sprintf("No Coord for team one so default coord selected\n")))
            cmdr.WriteLine([]byte("13,5"))
        }
        if !two {
            debug.WriteLine([]byte(fmt.Sprintf("No Coord for team two so default coord selected\n")))
            cmdr.WriteLine([]byte("222,5"))
        }
    }
}

