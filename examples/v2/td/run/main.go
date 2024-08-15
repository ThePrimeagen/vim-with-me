package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/examples/v2/td/players"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
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

	playerOneStr := ""
	flag.StringVar(&playerOneStr, "one", "", "player one string")

	playerTwoStr := ""
	flag.StringVar(&playerTwoStr, "two", "", "player two string")

	roundTimeInt := int64(0)
	flag.Int64Var(&roundTimeInt, "roundTime", 0, "seconds per round time")

	viz := false
	flag.BoolVar(&viz, "viz", false, "displays the game")

	seed := 1337
	flag.IntVar(&seed, "seed", 69420, "the seed value for the game")
	flag.Parse()

	args := flag.Args()
    assert.Assert(len(args) >= 2, "you must provide path to exec and json file")
	name := args[0]
	json := args[1]

    if roundTimeInt == 0 {
        roundTimeInt = int64(time.Second * 20)
    } else {
        roundTimeInt = int64(time.Second * time.Duration(roundTimeInt));
    }
    roundTime := time.Duration(roundTimeInt)

    debug, err := getDebug(debugFile)
    if err != nil {
        log.Fatalf("could not open up debug file: %v\n", err)
    }
    defer debug.Close()

	ctx := context.Background()

    cmdParser := td.NewCmdErrParser(debug)

    prog := cmd.NewCmder(name, ctx).
        AddVArg(json).
        AddKVArg("--seed", fmt.Sprintf("%d", seed)).
        WithErrFn(cmdParser.Parse).
        WithOutFn(func(b []byte) (int, error) {
            if viz {
                fmt.Printf("%s\n", string(b))
            }
            return len(b), nil
        })

    cmdr := td.TDCommander {
        Cmdr: prog,
        Debug: debug,
    }

    go prog.Run()

    one := players.NewTeamPlayerFromString(playerOneStr, debug, ctx, '1', cmdr)
    two := players.NewTeamPlayerFromString(playerTwoStr, debug, ctx, '2', cmdr)

    go one.Player.Run(ctx);
    go two.Player.Run(ctx);

    round := 0
    fmt.Printf("won,round,prompt file,seed,ai total towers,ai guesses,ai bad parses\n")

    defer func() {
        fmt.Println("\x1b[?25h")
    }()

    outer:
    for {
        debug.WriteStrLine(fmt.Sprintf("------------- waiting on game round: %d -----------------", round))
        select {
        case <-ctx.Done():
            break outer;
        case gs := <- cmdParser.Gs:
            debug.WriteStrLine(fmt.Sprintf("ai-placement response: \"%s\"", gs.String()))
            round = int(gs.Round)

            if gs.Finished {
                name := players.GetPromptName(playerOneStr) + players.GetPromptName(playerTwoStr)
                oneStats := one.Player.Stats()
                stats := oneStats.Add(two.Player.Stats())

                if gs.Winner == '1' {
                    fmt.Printf("1,%d,%s,%d,%s\n", round, name, seed, stats.String())
                } else {
                    fmt.Printf("2,%d,%s,%d,%s\n", round, name, seed, stats.String())
                }
                break outer
            }

            if gs.Playing {
                continue
            }

            one.Player.StartRound()
            two.Player.StartRound()

            innerCtx, cancel := context.WithCancel(ctx)

            t := time.NewTimer(roundTime)
            cmdr.Countdown(roundTime)

            go one.StreamMoves(innerCtx, &gs)
            go two.StreamMoves(innerCtx, &gs)

            <-t.C
            cancel()
            cmdr.PlayRound()
        }
    }
}

