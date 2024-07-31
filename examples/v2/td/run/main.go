package main

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
	"github.com/theprimeagen/vim-with-me/pkg/v2/cmd"
)

func main() {

    testies.SetupLogger()

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

            fmt.Printf("state: %s\n", str)

            if parts[0] == "waiting" {
                assert.Assert(len(parts) == 2, "somehow i didn't communicate back with 2 pieces", "parts", parts)
                parts[1] = strings.TrimSpace(parts[1])

                count, err := strconv.Atoi(parts[1])
                assert.NoError(err, "zig program gave me a non number")

                for range count {
                    ch<-struct{}{}
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

    interval := time.NewTicker(time.Second * 10)

    for range ch {
        <-interval.C
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
                fmt.Printf("COORD %s\n", coord.String())
                cmdr.WriteLine([]byte(coord.String()))
                one = true
            }
            if !two && coord.Team == 2 {
                fmt.Printf("COORD %s\n", coord.String())
                cmdr.WriteLine([]byte(coord.String()))
                two = true
            }

            if one && two {
                break
            }
        }

        if !one {
            cmdr.WriteLine([]byte("13,5"))
        }
        if !two {
            cmdr.WriteLine([]byte("13,5"))
        }
    }
}

