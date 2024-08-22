package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/examples/v2/td/objects"
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

func runMoves(one players.TeamPlayer, two players.TeamPlayer, ctx context.Context, gs *objects.GameState) chan struct{} {
	out := make(chan struct{}, 1)

	go func() {
		waitGroup := sync.WaitGroup{}
		waitGroup.Add(2)

		go func() {
			one.StreamMoves(ctx, gs)
			waitGroup.Done()
		}()

		go func() {
			two.StreamMoves(ctx, gs)
			waitGroup.Done()
		}()
		waitGroup.Wait()
		out <- struct{}{}
		close(out)
	}()

	return out
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

	hasVizFile := false
	flag.BoolVar(&hasVizFile, "vizFile", false, "displays the game in a file")

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
		roundTimeInt = int64(time.Second * time.Duration(roundTimeInt))
	}
	roundTime := time.Duration(roundTimeInt)

	debug, err := getDebug(debugFile)
	if err != nil {
		log.Fatalf("could not open up debug file: %v\n", err)
	}
	defer debug.Close()

	ctx := context.Background()

	cmdParser := td.NewCmdErrParser(debug)
    vizFile, err := os.OpenFile("/tmp/td-viz", os.O_RDWR|os.O_CREATE, 0644)

	prog := cmd.NewCmder(name, ctx).
		AddVArg(json).
		AddKVArg("--seed", fmt.Sprintf("%d", seed)).
		WithErrFn(cmdParser.Parse).
		WithOutFn(func(b []byte) (int, error) {
			if viz {
				fmt.Printf("%s\n", string(b))
			}
            if vizFile != nil {
				fmt.Fprintf(vizFile, "%s\n", string(b))
            }
			return len(b), nil
		})

	cmdr := td.TDCommander{
		Cmdr:  prog,
		Debug: debug,
	}

	go prog.Run()

	one := players.NewTeamPlayerFromString(playerOneStr, debug, ctx, '1', cmdr)
	two := players.NewTeamPlayerFromString(playerTwoStr, debug, ctx, '2', cmdr)

	go one.Player.Run(ctx)
	go two.Player.Run(ctx)

	round := 0
	fmt.Printf("won,round,one-team,two-team,seed,oneTotalTowersBuild,oneTotalProjectiles,oneTotalTowerUpgrades,oneTotalCreepDamage,oneTotalTowerDamage,oneTotalDamageFromCreeps,twoTotalTowersBuild,twoTotalProjectiles,twoTotalTowerUpgrades,twoTotalCreepDamage,twoTotalTowerDamage,twoTotalDamageFromCreeps\n")

	defer func() {
		fmt.Println("\x1b[?25h")
	}()

outer:
	for {
		debug.WriteStrLine(fmt.Sprintf("------------- waiting on game round: %d -----------------", round))
		select {
		case <-ctx.Done():
			break outer
		case gs := <-cmdParser.Gs:
			debug.WriteStrLine(fmt.Sprintf("ai-placement response(%d): \"%s\"", gs.OneTotalDamageFromCreeps, gs.String()))
			round = int(gs.Round)

			if gs.Finished {
                //           w  r 1n 2n seed
                fmt.Printf("%d,%d,%s,%s,%d,", gs.Winner, gs.Round, one.Player.Name(), two.Player.Name(), seed);
                // oneTotalTowersBuild, oneTotalProjectiles, oneTotalTowerUpgrades, oneTotalCreepDamage, oneTotalTowerDamage,
                fmt.Printf("%d,%d,%d,%d,%d,%d,", gs.OneTotalTowersBuild, gs.OneTotalProjectiles, gs.OneTotalTowerUpgrades, gs.OneTotalCreepDamage, gs.OneTotalTowerDamage, gs.OneTotalDamageFromCreeps)
                // twoTotalTowersBuild, twoTotalProjectiles, twoTotalTowerUpgrades, twoTotalCreepDamage, twoTotalTowerDamage")
                fmt.Printf("%d,%d,%d,%d,%d,%d\n", gs.TwoTotalTowersBuild, gs.TwoTotalProjectiles, gs.TwoTotalTowerUpgrades, gs.TwoTotalCreepDamage, gs.TwoTotalTowerDamage, gs.TwoTotalDamageFromCreeps)
				break outer
			}

			if gs.Playing {
				continue
			}

			one.Player.StartRound()
			two.Player.StartRound()

			innerCtx, cancel := context.WithCancel(ctx)
			moves := runMoves(one, two, innerCtx, &gs)

			t := time.NewTimer(roundTime)
			cmdr.Countdown(roundTime)

			select {
			case <-t.C:
			case <-moves:
			}

			t.Stop()
			cancel()

            one.Player.EndRound(&gs, cmdr)
            two.Player.EndRound(&gs, cmdr)

			cmdr.PlayRound()
		}
	}
}
