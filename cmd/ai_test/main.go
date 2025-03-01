package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type Challenger struct {
	One string
	Two string
}

const GPT3Dot5Turbo = "gpt-3.5-turbo"
const GPT4 = "gpt-4"
const GPT4oMini = "gpt-4o-mini"
const GPTo1 = "o1"
const GPTo3Mini = "o3-mini"

const ClaudeSonnet3_5 = "claude-3-5-sonnet-latest"
const ClaudeSonnet3_7 = "claude-3-7-sonnet-latest"

var claudes = []string{ClaudeSonnet3_5, ClaudeSonnet3_7}
var gpts = []string{GPT3Dot5Turbo, GPT4, GPT4oMini, GPTo1, GPTo3Mini}

func createTmp(name string) string {
	return path.Join("/tmp", fmt.Sprintf("%s-%d", name, time.Now().UnixMilli()))
}

func run(ctx context.Context, out io.WriteCloser, args []string, round int) error {
	vizFileTmp, err := os.CreateTemp("/tmp", "viz-")
	assert.NoError(err, "unable to create the viz file.  Temps should never fail")
	slog.Warn("VizFile", "name", vizFileTmp.Name(), "round", round)
	cmd := exec.Command("go", append([]string{
		"run",
		"./examples/v2/td/run/main.go",
		"--silent",
		"--vizFile", vizFileTmp.Name(),
	}, args...)...)
	cmd.Env = append(os.Environ(), "LEVEL=silent")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		slog.Error("error running cmd", "error", err)
	}

	/*
	buf := bytes.NewBuffer(make([]byte, 0, 1000))
	cmd.Stdout = buf

	ch := make(chan struct{})
	go func() {
		fmt.Printf("running: %v\n", cmd.String())
		err := cmd.Run()
		if err != nil {
			slog.Error("error running cmd", "error", err)
		}
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		cmd.Process.Kill()
		fmt.Printf("killing process", cmd.String())
		return fmt.Errorf("context cancelled before game ran its course: %w", ctx.Err())
	case <-ch:
		output, err := io.ReadAll(buf)
		fmt.Printf("process finished nicely", cmd.String())
		if err != nil {
			slog.Error("somehow got error while reading from buf output of cmd", "error", err)
		}
		out.Write(output)
		out.Close()
	}
	*/

	return nil
}

type RunParams struct {
	One       string `json:"one"`       // ai:anthropic:claude-3-7-sonnet-latest:prompt/THEPRIMEAGEN
	Two       string `json:"two"`       // ai:openai:model:prompt/THEPRIMEAGEN
	RoundTime int    `json:"roundTime"` // 10
	Seed      int    `json:"seed"`      // 42069
	DebugPath string `json:"debugPath"` // /tmp/td
	VizPath string `json:"vizPath"` // /tmp/td
}

func (r *RunParams) String() string {
	return fmt.Sprintf("RunParams{One: %s, Two: %s, RoundTime: %d, Seed: %d, DebugPath: %s, VizPath: %s}", r.One, r.Two, r.RoundTime, r.Seed, r.DebugPath, r.VizPath)
}

func (r *RunParams) toArgs() []string {
	return []string{
		"--one", r.One,
		"--two", r.Two,
		"--roundTime", fmt.Sprintf("%d", r.RoundTime),
		"--seed", fmt.Sprintf("%d", r.Seed),
		"--debug", r.DebugPath,
		"--vizFile", r.VizPath,
		"--",
		"./zig-out/bin/to",
		"main.json",
	}
}

func createChallenger(one, two string) Challenger {
	return Challenger{
		One: one,
		Two: two,
	}
}

func runWithCancel(ctx context.Context, params RunParams, round int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	tmp, err := os.CreateTemp("/tmp", "td-")
	if err != nil {
		slog.Error("unable to run program, no temp file", "error", err)
		return err
	}
	slog.Warn("Running with temp file", "name", tmp.Name())
	err = run(ctx, tmp, params.toArgs(), round)
	if err != nil {
		slog.Error("error running", "error", err)
		return err
	}
	return nil
}

func main() {
	semaphore := make(chan struct{}, 2)
	semaphore <- struct{}{}
	semaphore <- struct{}{}

	round := 0
	for top := range 2 {
		for seed := range 5 {
			for anth := range len(claudes) {
				for gpt := range len(gpts) {
					round++
					one := fmt.Sprintf("ai:anthropic:%s:prompt/THEPRIMEAGEN", claudes[anth])
					two := fmt.Sprintf("ai:openai:%s:prompt/THEPRIMEAGEN", gpts[gpt])

					if top == 1 {
						tmp := one
						one = two
						two = tmp
					}

					params := RunParams{
						One: one,
						Two: two,
						RoundTime: 30,
						Seed: round,
						DebugPath: fmt.Sprintf("/tmp/td-%d-%d-%d-%d", top, seed, anth, gpt),
					}

					fmt.Printf("waiting for semaphore...\n")
					<-semaphore
					fmt.Printf("running: %v\n", params.String())
					go func () {
						err := runWithCancel(context.Background(), params, round)
						if err != nil {
							fmt.Printf("error: %s\n", err)
						}
						semaphore <- struct{}{}
					}()
				}
			}
		}
	}

	fmt.Printf("done\n")
}

func testRunner() {
	ctx := context.Background()
	params := RunParams{
		One: "ai:anthropic:claude-3-7-sonnet-latest:prompt/THEPRIMEAGEN",
		Two: "ai:openai:o3-mini:prompt/THEPRIMEAGEN",
		RoundTime: 30,
		Seed: 42069,
		DebugPath: createTmp("td-debug"),
	}
	err := run(ctx, os.Stdout, params.toArgs(), 1)
	if err != nil {
		slog.Error("error running the program", "error", err)
	}
}
